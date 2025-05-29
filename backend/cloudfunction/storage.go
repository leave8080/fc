package cloudfunction

import (
	"context"
	"time"
)

// Storage 定义存储接口
type Storage interface {
	// 函数CRUD操作
	CreateFunction(ctx context.Context, fn *Function) error
	GetFunction(ctx context.Context, id string) (*Function, error)
	UpdateFunction(ctx context.Context, fn *Function) error
	DeleteFunction(ctx context.Context, id string) error
	ListFunctions(ctx context.Context, filters map[string]interface{}) ([]*Function, error)

	// 代码存储
	SaveFunctionCode(ctx context.Context, functionID string, runtime string, code []byte) error
	GetFunctionCode(ctx context.Context, functionID string) ([]byte, error)

	// 执行日志
	SaveExecutionLog(ctx context.Context, log *ExecutionLog) error
	GetExecutionHistory(ctx context.Context, functionID string, limit int) ([]*ExecutionLog, error)

	// 健康检查
	HealthCheck(ctx context.Context) error
	Close() error
}

// ExecutionLog 执行日志
type ExecutionLog struct {
	ID         string                 `json:"id"`
	FunctionID string                 `json:"function_id"`
	RequestID  string                 `json:"request_id"`
	Event      map[string]interface{} `json:"event"`
	Result     map[string]interface{} `json:"result"`
	Duration   time.Duration          `json:"duration"`
	Success    bool                   `json:"success"`
	Error      string                 `json:"error,omitempty"`
	ExecutedAt time.Time              `json:"executed_at"`
	UserID     string                 `json:"user_id,omitempty"`
}

// StorageConfig 存储配置
type StorageConfig struct {
	Type        string `yaml:"type"` // postgresql, mongodb, redis, s3
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Database    string `yaml:"database"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	SSLMode     string `yaml:"ssl_mode"`
	MaxConns    int    `yaml:"max_conns"`
	MaxIdleTime int    `yaml:"max_idle_time"`
	Bucket      string `yaml:"bucket,omitempty"`
	Region      string `yaml:"region,omitempty"`
	AccessKey   string `yaml:"access_key,omitempty"`
	SecretKey   string `yaml:"secret_key,omitempty"`
}

// FileStorage 文件存储实现（当前使用）
type FileStorage struct {
	dataFile string
	workDir  string
}

// NewFileStorage 创建文件存储实例
func NewFileStorage(workDir, dataFile string) Storage {
	return &FileStorage{
		dataFile: dataFile,
		workDir:  workDir,
	}
}

// 实现Storage接口的占位方法（实际实现在platform.go中）
func (f *FileStorage) CreateFunction(ctx context.Context, fn *Function) error {
	// 实际实现在Platform中
	return nil
}

func (f *FileStorage) GetFunction(ctx context.Context, id string) (*Function, error) {
	return nil, nil
}

func (f *FileStorage) UpdateFunction(ctx context.Context, fn *Function) error {
	return nil
}

func (f *FileStorage) DeleteFunction(ctx context.Context, id string) error {
	return nil
}

func (f *FileStorage) ListFunctions(ctx context.Context, filters map[string]interface{}) ([]*Function, error) {
	return nil, nil
}

func (f *FileStorage) SaveFunctionCode(ctx context.Context, functionID string, runtime string, code []byte) error {
	return nil
}

func (f *FileStorage) GetFunctionCode(ctx context.Context, functionID string) ([]byte, error) {
	return nil, nil
}

func (f *FileStorage) SaveExecutionLog(ctx context.Context, log *ExecutionLog) error {
	return nil
}

func (f *FileStorage) GetExecutionHistory(ctx context.Context, functionID string, limit int) ([]*ExecutionLog, error) {
	return nil, nil
}

func (f *FileStorage) HealthCheck(ctx context.Context) error {
	return nil
}

func (f *FileStorage) Close() error {
	return nil
}
