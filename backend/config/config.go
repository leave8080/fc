package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig
	Storage  StorageConfig
	Runtime  RuntimeConfig
	Security SecurityConfig
	Monitor  MonitorConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string
	Port         int
	Mode         string // debug, release
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

// StorageConfig 存储配置
type StorageConfig struct {
	Type        string
	DataDir     string
	MaxConns    int
	MaxIdleTime int
}

// RuntimeConfig 运行时配置
type RuntimeConfig struct {
	WorkDir         string
	MaxConcurrent   int
	DefaultTimeout  int
	DefaultMemory   int
	EnabledRuntimes []string
	MaxCodeSize     int64 // KB
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	EnableAuth     bool
	JWTExpiry      int
	AllowedOrigins []string
	RateLimit      int
	EnableHTTPS    bool
}

// MonitorConfig 监控配置
type MonitorConfig struct {
	EnableMetrics   bool
	MetricsPath     string
	EnableProfiling bool
	ProfilingPath   string
	LogLevel        string
	LogFormat       string // json, text
	LogOutput       string // stdout, file
	EnableTracing   bool
}

// Load 加载配置
func Load() *Config {
	config := &Config{
		Server: ServerConfig{
			Host:         GetEnv("SERVER_HOST", "0.0.0.0"),
			Port:         GetEnvInt("PORT", 8080),
			Mode:         GetEnv("GIN_MODE", "debug"),
			ReadTimeout:  30,
			WriteTimeout: 30,
			IdleTimeout:  120,
		},
		Storage: StorageConfig{
			Type:        GetEnv("STORAGE_TYPE", "file"),
			DataDir:     GetEnv("DATA_DIR", "./data"),
			MaxConns:    10,
			MaxIdleTime: 300,
		},
		Runtime: RuntimeConfig{
			WorkDir:         GetEnv("FUNCTIONS_DIR", "./functions"),
			MaxConcurrent:   GetEnvInt("MAX_CONCURRENT", 10),
			DefaultTimeout:  GetEnvInt("DEFAULT_TIMEOUT", 30),
			DefaultMemory:   GetEnvInt("DEFAULT_MEMORY", 128),
			EnabledRuntimes: []string{"go", "nodejs", "python"},
			MaxCodeSize:     1024, // 1MB
		},
		Security: SecurityConfig{
			EnableAuth:     GetEnvBool("ENABLE_AUTH", false),
			JWTExpiry:      24 * 3600, // 24小时
			AllowedOrigins: []string{"*"},
			RateLimit:      100,
			EnableHTTPS:    GetEnvBool("ENABLE_HTTPS", false),
		},
		Monitor: MonitorConfig{
			EnableMetrics:   GetEnvBool("ENABLE_METRICS", true),
			MetricsPath:     "/metrics",
			EnableProfiling: false,
			ProfilingPath:   "/debug/pprof",
			LogLevel:        GetEnv("LOG_LEVEL", "info"),
			LogFormat:       "json",
			LogOutput:       "stdout",
			EnableTracing:   false,
		},
	}

	// 验证配置
	if err := validate(config); err != nil {
		fmt.Printf("配置验证失败: %v\n", err)
		os.Exit(1)
	}

	return config
}

func validate(config *Config) error {
	// 验证必要的配置
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("无效的端口号: %d", config.Server.Port)
	}

	if config.Runtime.MaxConcurrent <= 0 {
		return fmt.Errorf("最大并发数必须大于0")
	}

	if config.Runtime.DefaultTimeout <= 0 {
		return fmt.Errorf("默认超时时间必须大于0")
	}

	// 验证日志级别
	validLogLevels := []string{"debug", "info", "warn", "error"}
	found := false
	for _, level := range validLogLevels {
		if config.Monitor.LogLevel == level {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("无效的日志级别: %s", config.Monitor.LogLevel)
	}

	// 创建必要的目录
	dirs := []string{
		config.Storage.DataDir,
		config.Runtime.WorkDir,
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %v", dir, err)
		}
	}

	return nil
}

// GetEnv 获取环境变量，如果不存在则返回默认值
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt 获取整数类型的环境变量
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvBool 获取布尔类型的环境变量
func GetEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
