package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"testChat/backend/cloudfunction"
)

func main() {
	// 1. 初始化配置和日志
	initializeSystem()

	// 2. 创建云函数平台
	platform := createPlatform()

	// 3. 启动服务器
	startServer(platform)
}

// initializeSystem 初始化系统组件
func initializeSystem() {
	// 初始化日志系统
	logLevel := getEnv("LOG_LEVEL", "info")
	logger := cloudfunction.NewLogger(logLevel)
	cloudfunction.GlobalLogger = logger

	logger.Info("云函数平台启动中...")

	// 初始化指标收集
	metrics := cloudfunction.NewMetrics()
	cloudfunction.GlobalMetrics = metrics
}

// createPlatform 创建云函数平台
func createPlatform() *cloudfunction.Platform {
	// 设置云函数工作目录
	functionsDir := getEnv("FUNCTIONS_DIR", "./functions")
	if err := os.MkdirAll(functionsDir, 0755); err != nil {
		cloudfunction.GlobalLogger.Fatal("创建云函数目录失败: %v", err)
	}

	// 创建云函数平台
	platform := cloudfunction.NewPlatform(functionsDir)
	cloudfunction.GlobalLogger.Info("云函数平台初始化完成")

	return platform
}

// startServer 启动服务器
func startServer(platform *cloudfunction.Platform) {
	// 创建服务器
	server := cloudfunction.NewServer(platform)

	// 获取端口配置
	port := getEnvInt("PORT", 8080)

	// 启动服务器（在goroutine中）
	go func() {
		cloudfunction.GlobalLogger.Info("云函数平台启动成功，端口: %d", port)
		cloudfunction.GlobalLogger.Info("- 云函数API: /api/v1/functions/*")
		cloudfunction.GlobalLogger.Info("- 健康检查: /api/v1/health")

		if err := server.Run(port); err != nil {
			cloudfunction.GlobalLogger.Fatal("启动服务器失败: %v", err)
		}
	}()

	// 等待退出信号
	waitForShutdown()
}

// waitForShutdown 等待关闭信号
func waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cloudfunction.GlobalLogger.Info("收到关闭信号，正在关闭服务器...")

	// 这里可以添加优雅关闭逻辑
	time.Sleep(time.Second)
	cloudfunction.GlobalLogger.Info("服务器已关闭")
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 获取整数类型的环境变量
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
