package cloudfunction

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

// LogLevel 日志级别
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

var logLevelNames = map[LogLevel]string{
	LogLevelDebug: "DEBUG",
	LogLevelInfo:  "INFO",
	LogLevelWarn:  "WARN",
	LogLevelError: "ERROR",
	LogLevelFatal: "FATAL",
}

// Logger 日志记录器
type Logger struct {
	level  LogLevel
	logger *log.Logger
}

// NewLogger 创建新的日志记录器
func NewLogger(level string) *Logger {
	var logLevel LogLevel
	switch level {
	case "debug":
		logLevel = LogLevelDebug
	case "info":
		logLevel = LogLevelInfo
	case "warn":
		logLevel = LogLevelWarn
	case "error":
		logLevel = LogLevelError
	case "fatal":
		logLevel = LogLevelFatal
	default:
		logLevel = LogLevelInfo
	}

	return &Logger{
		level:  logLevel,
		logger: log.New(os.Stdout, "", 0),
	}
}

func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	// 获取调用信息
	_, file, line, _ := runtime.Caller(2)
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	message := fmt.Sprintf(format, args...)
	logEntry := fmt.Sprintf("[%s] %s %s:%d - %s",
		logLevelNames[level], timestamp, file, line, message)

	l.logger.Println(logEntry)
}

// Debug 调试日志
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LogLevelDebug, format, args...)
}

// Info 信息日志
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LogLevelInfo, format, args...)
}

// Warn 警告日志
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LogLevelWarn, format, args...)
}

// Error 错误日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LogLevelError, format, args...)
}

// Fatal 致命错误日志
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(LogLevelFatal, format, args...)
	os.Exit(1)
}

// 全局日志实例
var GlobalLogger = NewLogger("info")

// 便捷的全局日志函数
func Debug(format string, args ...interface{}) {
	GlobalLogger.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	GlobalLogger.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	GlobalLogger.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	GlobalLogger.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	GlobalLogger.Fatal(format, args...)
}
