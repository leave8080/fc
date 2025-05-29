# 自实现云函数平台

## 概述

这是一个自己实现的云函数平台，支持多种编程语言运行时，提供完整的函数生命周期管理，无需依赖第三方云服务商。

## 核心特性

### 🚀 多运行时支持
- **Go**: 原生支持，编译执行
- **Node.js**: JavaScript/TypeScript支持
- **Python**: Python3支持

### 🔧 完整的函数管理
- ✅ 创建函数
- ✅ 更新函数
- ✅ 删除函数
- ✅ 列出函数
- ✅ 函数执行
- ✅ 实时监控

### 🛡️ 安全与隔离
- 函数沙盒环境
- 超时控制
- 内存限制
- 进程隔离

### 📊 性能监控
- 执行时间统计
- 内存使用监控
- 错误日志记录

## 架构设计

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend UI   │    │  HTTP API       │    │  Function       │
│                 │◄──►│                 │◄──►│  Platform       │
│ Vue.js + Axios  │    │ Gin Router      │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                      │
                       ┌─────────────────┬─────────────────┬─────────────────┐
                       │                 │                 │                 │
                ┌─────────────┐   ┌─────────────┐   ┌─────────────┐
                │ Go Executor │   │ Node.js     │   │ Python      │
                │             │   │ Executor    │   │ Executor    │
                └─────────────┘   └─────────────┘   └─────────────┘
```

## 快速开始

### 1. 环境要求

**必需**:
- Go 1.16+ (必须)
- Git

**可选**:
- Node.js 14+ (运行Node.js函数)
- Python 3.7+ (运行Python函数)

### 2. 启动服务

```bash
# 进入后端目录
cd backend

# 使用启动脚本（推荐）
./start_cloudfunction.sh

# 或者手动启动
go mod tidy
go run main.go
```

### 3. 验证服务

访问健康检查接口：
```bash
curl http://localhost:8080/api/health
```

预期响应：
```json
{
  "status": "ok",
  "services": {
    "crawler": "running",
    "cloudfunction": "running"
  }
}
```

## API文档

### 基础信息
- **Base URL**: `http://localhost:8080/api/v1`
- **Content-Type**: `application/json`

### 函数管理API

#### 1. 创建函数
```http
POST /functions
Content-Type: application/json

{
  "name": "hello-world",
  "runtime": "go",
  "handler": "handler",
  "code": "func handler(event interface{}, context map[string]string) interface{} {\n    return map[string]interface{}{\n        \"message\": \"Hello World!\"\n    }\n}",
  "timeout": 30,
  "memory": 128,
  "environment": {
    "KEY": "value"
  }
}
```

#### 2. 列出函数
```http
GET /functions
```

#### 3. 获取函数详情
```http
GET /functions/{id}
```

#### 4. 更新函数
```http
PUT /functions/{id}
Content-Type: application/json

{
  "name": "updated-hello",
  "timeout": 60
}
```

#### 5. 删除函数
```http
DELETE /functions/{id}
```

#### 6. 执行函数
```http
POST /functions/{id}/invoke
Content-Type: application/json

{
  "event": {
    "message": "Hello World",
    "data": [1, 2, 3]
  },
  "context": {
    "user_id": "123",
    "request_id": "req_456"
  }
}
```

## 函数编写指南

### Go函数

```go
// 基本格式
func handler(event interface{}, context map[string]string) interface{} {
    // 你的业务逻辑
    return result
}

// 示例：数据处理
func processData(event interface{}, context map[string]string) interface{} {
    eventMap := event.(map[string]interface{})
    data := eventMap["data"].([]interface{})
    
    sum := 0.0
    for _, v := range data {
        sum += v.(float64)
    }
    
    return map[string]interface{}{
        "sum": sum,
        "count": len(data),
        "average": sum / float64(len(data)),
    }
}
```

### Node.js函数

```javascript
// 同步函数
function handler(event, context) {
    return {
        message: "Hello from Node.js",
        event: event
    };
}

// 异步函数
async function fetchData(event, context) {
    const axios = require('axios');
    
    try {
        const response = await axios.get('https://api.example.com/data');
        return {
            success: true,
            data: response.data
        };
    } catch (error) {
        return {
            success: false,
            error: error.message
        };
    }
}
```

### Python函数

```python
# 基本函数
def handler(event, context):
    return {
        'message': 'Hello from Python',
        'event': event
    }

# 数据分析函数
def analyze_data(event, context):
    import statistics
    
    data = event.get('data', [])
    
    if not data:
        return {'error': 'No data provided'}
    
    return {
        'mean': statistics.mean(data),
        'median': statistics.median(data),
        'mode': statistics.mode(data),
        'stdev': statistics.stdev(data)
    }
```

## 前端管理界面

访问 `http://localhost:3000` (需要启动前端服务)

### 功能特性
- 📝 可视化函数编辑器
- 🎯 一键函数测试
- 📋 函数列表管理
- 🔄 实时执行结果
- 📊 执行时间统计

## 测试

### 自动化测试
```bash
# 启动服务器（在另一个终端）
./start_cloudfunction.sh

# 运行测试
cd test
go run cloudfunction_test.go
```

### 手动测试
```bash
# 创建测试函数
curl -X POST http://localhost:8080/api/v1/functions \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-hello",
    "runtime": "go",
    "handler": "handler",
    "code": "func handler(event interface{}, context map[string]string) interface{} {\n    return map[string]interface{}{\n        \"message\": \"Hello World!\",\n        \"event\": event\n    }\n}"
  }'

# 调用函数（替换{id}为实际函数ID）
curl -X POST http://localhost:8080/api/v1/functions/{id}/invoke \
  -H "Content-Type: application/json" \
  -d '{
    "event": {"test": "data"},
    "context": {"user": "test"}
  }'
```

## 配置选项

### 环境变量
- `CLOUDFUNCTION_PORT`: 服务端口 (默认: 8080)
- `CLOUDFUNCTION_WORKDIR`: 函数工作目录 (默认: ./functions)
- `CLOUDFUNCTION_TIMEOUT`: 默认超时时间 (默认: 30秒)

### 安全配置
- 函数执行超时限制
- 内存使用限制
- 进程隔离
- 禁止访问系统敏感目录

## 监控和日志

### 执行日志
所有函数执行都会记录：
- 开始时间
- 执行时长
- 内存使用
- 成功/失败状态
- 错误信息

### 性能监控
- 平均执行时间
- 成功率统计
- 错误类型分析

## 部署

### Docker部署
```dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o cloudfunction main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates nodejs python3
WORKDIR /root/
COPY --from=builder /app/cloudfunction .
EXPOSE 8080
CMD ["./cloudfunction"]
```

### Kubernetes部署
参考 `k8s-manifests.yaml` 文件

## 故障排除

### 常见问题

1. **函数编译失败**
   - 检查代码语法
   - 确认运行时已安装
   - 查看错误日志

2. **函数执行超时**
   - 调整超时时间
   - 优化函数代码
   - 检查网络依赖

3. **内存不足**
   - 增加内存限制
   - 优化内存使用
   - 避免内存泄漏

### 日志查看
```bash
# 查看服务日志
journalctl -u cloudfunction -f

# 查看函数执行日志
tail -f functions/*/logs/*.log
```

## 扩展开发

### 添加新运行时
1. 在 `executors.go` 中添加新的执行器
2. 更新 `platform.go` 中的运行时判断
3. 添加对应的测试用例

### 自定义中间件
```go
func customMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 自定义逻辑
        c.Next()
    }
}
```

## 贡献

欢迎提交 Issue 和 Pull Request！

### 开发流程
1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 创建 Pull Request

## 许可证

MIT License

## 联系方式

如有问题或建议，请提交 Issue 或联系维护者。 