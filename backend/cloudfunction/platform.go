package cloudfunction

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Function 表示一个云函数
type Function struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Runtime     string            `json:"runtime"`     // go, nodejs, python
	Code        string            `json:"code"`        // 函数代码
	Handler     string            `json:"handler"`     // 入口函数
	Environment map[string]string `json:"environment"` // 环境变量
	Timeout     int               `json:"timeout"`     // 超时时间(秒)
	Memory      int               `json:"memory"`      // 内存限制(MB)
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// ExecuteRequest 函数执行请求
type ExecuteRequest struct {
	Event   interface{}       `json:"event"`   // 事件数据
	Context map[string]string `json:"context"` // 执行上下文
}

// ExecuteResponse 函数执行响应
type ExecuteResponse struct {
	Success    bool        `json:"success"`
	Result     interface{} `json:"result"`
	Error      string      `json:"error,omitempty"`
	Duration   int64       `json:"duration"` // 执行时间(毫秒)
	MemoryUsed int         `json:"memory_used,omitempty"`
}

// Platform 云函数平台
type Platform struct {
	functions map[string]*Function
	workDir   string
	dataFile  string // 数据持久化文件路径
	mutex     sync.RWMutex
}

// NewPlatform 创建新的云函数平台
func NewPlatform(workDir string) *Platform {
	platform := &Platform{
		functions: make(map[string]*Function),
		workDir:   workDir,
		dataFile:  filepath.Join(workDir, "functions.json"),
	}

	// 尝试从文件加载现有函数
	platform.loadFromFile()

	return platform
}

// loadFromFile 从JSON文件加载函数列表
func (p *Platform) loadFromFile() error {
	// 检查文件是否存在
	if _, err := os.Stat(p.dataFile); os.IsNotExist(err) {
		// 文件不存在，创建空文件
		return p.saveToFile()
	}

	data, err := ioutil.ReadFile(p.dataFile)
	if err != nil {
		return fmt.Errorf("读取数据文件失败: %v", err)
	}

	// 如果文件为空，直接返回
	if len(data) == 0 {
		return nil
	}

	var functions []*Function
	if err := json.Unmarshal(data, &functions); err != nil {
		return fmt.Errorf("解析数据文件失败: %v", err)
	}

	// 将函数加载到内存
	for _, fn := range functions {
		p.functions[fn.ID] = fn
	}

	fmt.Printf("从文件加载了 %d 个函数\n", len(functions))
	return nil
}

// saveToFile 将函数列表保存到JSON文件
func (p *Platform) saveToFile() error {
	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(p.dataFile), 0755); err != nil {
		return fmt.Errorf("创建数据目录失败: %v", err)
	}

	functions := make([]*Function, 0, len(p.functions))
	for _, fn := range p.functions {
		functions = append(functions, fn)
	}

	data, err := json.MarshalIndent(functions, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化数据失败: %v", err)
	}

	if err := ioutil.WriteFile(p.dataFile, data, 0644); err != nil {
		return fmt.Errorf("写入数据文件失败: %v", err)
	}

	return nil
}

// CreateFunction 创建新函数
func (p *Platform) CreateFunction(fn *Function) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	fn.ID = generateID()
	fn.CreatedAt = time.Now()
	fn.UpdatedAt = time.Now()

	// 创建函数工作目录
	fnDir := filepath.Join(p.workDir, fn.ID)
	if err := os.MkdirAll(fnDir, 0755); err != nil {
		return fmt.Errorf("创建函数目录失败: %v", err)
	}

	// 保存函数代码
	if err := p.saveFunction(fn); err != nil {
		return fmt.Errorf("保存函数失败: %v", err)
	}

	p.functions[fn.ID] = fn

	// 持久化到文件
	if err := p.saveToFile(); err != nil {
		// 如果持久化失败，回滚内存操作
		delete(p.functions, fn.ID)
		os.RemoveAll(fnDir)
		return fmt.Errorf("持久化函数失败: %v", err)
	}

	return nil
}

// GetFunction 获取函数
func (p *Platform) GetFunction(id string) (*Function, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	fn, exists := p.functions[id]
	if !exists {
		return nil, fmt.Errorf("函数不存在: %s", id)
	}
	return fn, nil
}

// ListFunctions 列出所有函数
func (p *Platform) ListFunctions() []*Function {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	functions := make([]*Function, 0, len(p.functions))
	for _, fn := range p.functions {
		functions = append(functions, fn)
	}
	return functions
}

// UpdateFunction 更新函数
func (p *Platform) UpdateFunction(id string, fn *Function) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	existing, exists := p.functions[id]
	if !exists {
		return fmt.Errorf("函数不存在: %s", id)
	}

	// 备份原函数用于回滚
	backup := *existing

	fn.ID = id
	fn.CreatedAt = existing.CreatedAt
	fn.UpdatedAt = time.Now()

	if err := p.saveFunction(fn); err != nil {
		return fmt.Errorf("保存函数失败: %v", err)
	}

	p.functions[id] = fn

	// 持久化到文件
	if err := p.saveToFile(); err != nil {
		// 如果持久化失败，回滚
		p.functions[id] = &backup
		return fmt.Errorf("持久化函数失败: %v", err)
	}

	return nil
}

// DeleteFunction 删除函数
func (p *Platform) DeleteFunction(id string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	fn, exists := p.functions[id]
	if !exists {
		return fmt.Errorf("函数不存在: %s", id)
	}

	// 备份函数用于回滚
	backup := *fn

	// 删除函数目录
	fnDir := filepath.Join(p.workDir, id)
	if err := os.RemoveAll(fnDir); err != nil {
		return fmt.Errorf("删除函数目录失败: %v", err)
	}

	delete(p.functions, id)

	// 持久化到文件
	if err := p.saveToFile(); err != nil {
		// 如果持久化失败，回滚
		p.functions[id] = &backup
		// 尝试恢复目录（尽力而为）
		os.MkdirAll(fnDir, 0755)
		p.saveFunction(&backup)
		return fmt.Errorf("持久化删除操作失败: %v", err)
	}

	return nil
}

// ExecuteFunction 执行函数
func (p *Platform) ExecuteFunction(id string, req *ExecuteRequest) (*ExecuteResponse, error) {
	fn, err := p.GetFunction(id)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()

	// 根据运行时执行函数
	var result interface{}
	var execErr error

	switch fn.Runtime {
	case "go":
		result, execErr = p.executeGoFunction(fn, req)
	case "nodejs":
		result, execErr = p.executeNodeJSFunction(fn, req)
	case "python":
		result, execErr = p.executePythonFunction(fn, req)
	default:
		execErr = fmt.Errorf("不支持的运行时: %s", fn.Runtime)
	}

	duration := time.Since(startTime).Milliseconds()

	response := &ExecuteResponse{
		Success:  execErr == nil,
		Result:   result,
		Duration: duration,
	}

	if execErr != nil {
		response.Error = execErr.Error()
	}

	return response, nil
}

// saveFunction 保存函数代码到文件
func (p *Platform) saveFunction(fn *Function) error {
	fnDir := filepath.Join(p.workDir, fn.ID)

	var filename string
	switch fn.Runtime {
	case "go":
		filename = "main.go"
	case "nodejs":
		filename = "index.js"
	case "python":
		filename = "main.py"
	default:
		return fmt.Errorf("不支持的运行时: %s", fn.Runtime)
	}

	filepath := filepath.Join(fnDir, filename)
	return ioutil.WriteFile(filepath, []byte(fn.Code), 0644)
}

// generateID 生成唯一ID
func generateID() string {
	return fmt.Sprintf("fn_%d", time.Now().UnixNano())
}
