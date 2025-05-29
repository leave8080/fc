package cloudfunction

import (
	"sync"
	"sync/atomic"
	"time"
)

// Metrics 指标收集器
type Metrics struct {
	// 函数执行统计
	FunctionExecutions   int64 `json:"function_executions"`
	SuccessfulExecutions int64 `json:"successful_executions"`
	FailedExecutions     int64 `json:"failed_executions"`

	// 响应时间统计
	TotalExecutionTime int64 `json:"total_execution_time_ms"`
	MinExecutionTime   int64 `json:"min_execution_time_ms"`
	MaxExecutionTime   int64 `json:"max_execution_time_ms"`

	// 函数管理统计
	CreatedFunctions int64 `json:"created_functions"`
	DeletedFunctions int64 `json:"deleted_functions"`

	// 运行时统计
	RuntimeUsage map[string]int64 `json:"runtime_usage"`

	// 错误统计
	ErrorsByType map[string]int64 `json:"errors_by_type"`

	// 系统指标
	StartTime time.Time `json:"start_time"`

	mu sync.RWMutex
}

// NewMetrics 创建新的指标收集器
func NewMetrics() *Metrics {
	return &Metrics{
		RuntimeUsage:     make(map[string]int64),
		ErrorsByType:     make(map[string]int64),
		StartTime:        time.Now(),
		MinExecutionTime: 999999999, // 初始化为一个很大的值
	}
}

// RecordExecution 记录函数执行
func (m *Metrics) RecordExecution(runtime string, duration time.Duration, success bool, errorType string) {
	atomic.AddInt64(&m.FunctionExecutions, 1)

	durationMs := duration.Milliseconds()
	atomic.AddInt64(&m.TotalExecutionTime, durationMs)

	// 更新最小执行时间
	for {
		current := atomic.LoadInt64(&m.MinExecutionTime)
		if durationMs >= current || atomic.CompareAndSwapInt64(&m.MinExecutionTime, current, durationMs) {
			break
		}
	}

	// 更新最大执行时间
	for {
		current := atomic.LoadInt64(&m.MaxExecutionTime)
		if durationMs <= current || atomic.CompareAndSwapInt64(&m.MaxExecutionTime, current, durationMs) {
			break
		}
	}

	if success {
		atomic.AddInt64(&m.SuccessfulExecutions, 1)
	} else {
		atomic.AddInt64(&m.FailedExecutions, 1)

		// 记录错误类型
		m.mu.Lock()
		m.ErrorsByType[errorType]++
		m.mu.Unlock()
	}

	// 记录运行时使用情况
	m.mu.Lock()
	m.RuntimeUsage[runtime]++
	m.mu.Unlock()
}

// RecordFunctionCreated 记录函数创建
func (m *Metrics) RecordFunctionCreated() {
	atomic.AddInt64(&m.CreatedFunctions, 1)
}

// RecordFunctionDeleted 记录函数删除
func (m *Metrics) RecordFunctionDeleted() {
	atomic.AddInt64(&m.DeletedFunctions, 1)
}

// GetSnapshot 获取指标快照
func (m *Metrics) GetSnapshot() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	totalExecutions := atomic.LoadInt64(&m.FunctionExecutions)
	totalTime := atomic.LoadInt64(&m.TotalExecutionTime)

	var avgExecutionTime float64
	if totalExecutions > 0 {
		avgExecutionTime = float64(totalTime) / float64(totalExecutions)
	}

	var successRate float64
	if totalExecutions > 0 {
		successRate = float64(atomic.LoadInt64(&m.SuccessfulExecutions)) / float64(totalExecutions) * 100
	}

	// 复制map以避免并发问题
	runtimeUsage := make(map[string]int64)
	for k, v := range m.RuntimeUsage {
		runtimeUsage[k] = v
	}

	errorsByType := make(map[string]int64)
	for k, v := range m.ErrorsByType {
		errorsByType[k] = v
	}

	return map[string]interface{}{
		"uptime_seconds":          time.Since(m.StartTime).Seconds(),
		"function_executions":     totalExecutions,
		"successful_executions":   atomic.LoadInt64(&m.SuccessfulExecutions),
		"failed_executions":       atomic.LoadInt64(&m.FailedExecutions),
		"success_rate_percent":    successRate,
		"avg_execution_time_ms":   avgExecutionTime,
		"min_execution_time_ms":   atomic.LoadInt64(&m.MinExecutionTime),
		"max_execution_time_ms":   atomic.LoadInt64(&m.MaxExecutionTime),
		"total_execution_time_ms": totalTime,
		"created_functions":       atomic.LoadInt64(&m.CreatedFunctions),
		"deleted_functions":       atomic.LoadInt64(&m.DeletedFunctions),
		"runtime_usage":           runtimeUsage,
		"errors_by_type":          errorsByType,
		"start_time":              m.StartTime,
	}
}

// Reset 重置指标
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	atomic.StoreInt64(&m.FunctionExecutions, 0)
	atomic.StoreInt64(&m.SuccessfulExecutions, 0)
	atomic.StoreInt64(&m.FailedExecutions, 0)
	atomic.StoreInt64(&m.TotalExecutionTime, 0)
	atomic.StoreInt64(&m.MinExecutionTime, 999999999)
	atomic.StoreInt64(&m.MaxExecutionTime, 0)
	atomic.StoreInt64(&m.CreatedFunctions, 0)
	atomic.StoreInt64(&m.DeletedFunctions, 0)

	m.RuntimeUsage = make(map[string]int64)
	m.ErrorsByType = make(map[string]int64)
	m.StartTime = time.Now()
}

// 全局指标实例
var GlobalMetrics = NewMetrics()

// PerformanceMonitor 性能监控
type PerformanceMonitor struct {
	metrics *Metrics
	alerts  []AlertRule
	mu      sync.RWMutex
}

// AlertRule 告警规则
type AlertRule struct {
	Name      string
	Condition func(metrics map[string]interface{}) bool
	Message   string
	Triggered bool
}

// NewPerformanceMonitor 创建性能监控器
func NewPerformanceMonitor(metrics *Metrics) *PerformanceMonitor {
	monitor := &PerformanceMonitor{
		metrics: metrics,
		alerts:  []AlertRule{},
	}

	// 添加默认告警规则
	monitor.AddAlert("high_error_rate", func(m map[string]interface{}) bool {
		if rate, ok := m["success_rate_percent"].(float64); ok {
			return rate < 90.0 // 成功率低于90%
		}
		return false
	}, "函数执行成功率低于90%")

	monitor.AddAlert("slow_execution", func(m map[string]interface{}) bool {
		if avg, ok := m["avg_execution_time_ms"].(float64); ok {
			return avg > 5000 // 平均执行时间超过5秒
		}
		return false
	}, "函数平均执行时间超过5秒")

	return monitor
}

// AddAlert 添加告警规则
func (pm *PerformanceMonitor) AddAlert(name string, condition func(map[string]interface{}) bool, message string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.alerts = append(pm.alerts, AlertRule{
		Name:      name,
		Condition: condition,
		Message:   message,
		Triggered: false,
	})
}

// CheckAlerts 检查告警
func (pm *PerformanceMonitor) CheckAlerts() []string {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	snapshot := pm.metrics.GetSnapshot()
	var triggeredAlerts []string

	for i := range pm.alerts {
		alert := &pm.alerts[i]
		shouldTrigger := alert.Condition(snapshot)

		if shouldTrigger && !alert.Triggered {
			triggeredAlerts = append(triggeredAlerts, alert.Message)
			alert.Triggered = true
		} else if !shouldTrigger && alert.Triggered {
			alert.Triggered = false
		}
	}

	return triggeredAlerts
}

// 全局性能监控实例
var GlobalPerformanceMonitor = NewPerformanceMonitor(GlobalMetrics)
