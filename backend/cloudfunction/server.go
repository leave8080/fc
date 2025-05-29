package cloudfunction

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Server 云函数HTTP服务器
type Server struct {
	platform *Platform
	router   *gin.Engine
}

// NewServer 创建新的服务器
func NewServer(platform *Platform) *Server {
	server := &Server{
		platform: platform,
		router:   gin.Default(),
	}
	server.setupRoutes()
	return server
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// 添加CORS中间件
	s.router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	api := s.router.Group("/api/v1")
	{
		// 函数管理
		api.POST("/functions", s.createFunction)
		api.GET("/functions", s.listFunctions)
		api.GET("/functions/:id", s.getFunction)
		api.PUT("/functions/:id", s.updateFunction)
		api.DELETE("/functions/:id", s.deleteFunction)

		// 函数执行
		api.POST("/functions/:id/invoke", s.invokeFunction)

		// 健康检查
		api.GET("/health", s.healthCheck)
	}
}

// createFunction 创建函数
func (s *Server) createFunction(c *gin.Context) {
	var req struct {
		Name        string            `json:"name" binding:"required"`
		Runtime     string            `json:"runtime" binding:"required"`
		Code        string            `json:"code" binding:"required"`
		Handler     string            `json:"handler" binding:"required"`
		Environment map[string]string `json:"environment"`
		Timeout     int               `json:"timeout"`
		Memory      int               `json:"memory"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 设置默认值
	if req.Timeout == 0 {
		req.Timeout = 30 // 30秒默认超时
	}
	if req.Memory == 0 {
		req.Memory = 128 // 128MB默认内存
	}
	if req.Environment == nil {
		req.Environment = make(map[string]string)
	}

	fn := &Function{
		Name:        req.Name,
		Runtime:     req.Runtime,
		Code:        req.Code,
		Handler:     req.Handler,
		Environment: req.Environment,
		Timeout:     req.Timeout,
		Memory:      req.Memory,
	}

	if err := s.platform.CreateFunction(fn); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建函数失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "函数创建成功",
		"function": fn,
	})
}

// listFunctions 列出所有函数
func (s *Server) listFunctions(c *gin.Context) {
	functions := s.platform.ListFunctions()
	c.JSON(http.StatusOK, gin.H{
		"functions": functions,
		"count":     len(functions),
	})
}

// getFunction 获取函数详情
func (s *Server) getFunction(c *gin.Context) {
	id := c.Param("id")
	fn, err := s.platform.GetFunction(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"function": fn})
}

// updateFunction 更新函数
func (s *Server) updateFunction(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name        string            `json:"name"`
		Runtime     string            `json:"runtime"`
		Code        string            `json:"code"`
		Handler     string            `json:"handler"`
		Environment map[string]string `json:"environment"`
		Timeout     int               `json:"timeout"`
		Memory      int               `json:"memory"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 获取现有函数
	existing, err := s.platform.GetFunction(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 更新字段（只更新非空字段）
	fn := *existing
	if req.Name != "" {
		fn.Name = req.Name
	}
	if req.Runtime != "" {
		fn.Runtime = req.Runtime
	}
	if req.Code != "" {
		fn.Code = req.Code
	}
	if req.Handler != "" {
		fn.Handler = req.Handler
	}
	if req.Environment != nil {
		fn.Environment = req.Environment
	}
	if req.Timeout > 0 {
		fn.Timeout = req.Timeout
	}
	if req.Memory > 0 {
		fn.Memory = req.Memory
	}

	if err := s.platform.UpdateFunction(id, &fn); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新函数失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "函数更新成功",
		"function": &fn,
	})
}

// deleteFunction 删除函数
func (s *Server) deleteFunction(c *gin.Context) {
	id := c.Param("id")

	if err := s.platform.DeleteFunction(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "函数删除成功"})
}

// invokeFunction 调用函数
func (s *Server) invokeFunction(c *gin.Context) {
	id := c.Param("id")

	var req ExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	response, err := s.platform.ExecuteFunction(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "函数执行失败: " + err.Error()})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusInternalServerError, response)
	}
}

// healthCheck 健康检查
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "云函数平台",
		"version": "1.0.0",
	})
}

// GetRouter 获取gin路由器实例
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

// Run 启动服务器
func (s *Server) Run(port int) error {
	addr := ":" + strconv.Itoa(port)
	fmt.Printf("云函数平台启动在端口 %d\n", port)
	return s.router.Run(addr)
}
