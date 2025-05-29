# 云函数平台使用示例

## 1. Go语言函数示例

### 简单的Hello World函数

```go
func handler(event interface{}, context map[string]string) interface{} {
    return map[string]interface{}{
        "message": "Hello from Go Cloud Function!",
        "timestamp": time.Now().Format(time.RFC3339),
        "event": event,
    }
}
```

**创建请求**:
```json
{
    "name": "go-hello",
    "runtime": "go",
    "handler": "handler",
    "code": "import \"time\"\n\nfunc handler(event interface{}, context map[string]string) interface{} {\n    return map[string]interface{}{\n        \"message\": \"Hello from Go Cloud Function!\",\n        \"timestamp\": time.Now().Format(time.RFC3339),\n        \"event\": event,\n    }\n}",
    "timeout": 30,
    "memory": 128
}
```

### 数据处理函数

```go
import (
    "encoding/json"
    "strings"
)

func processData(event interface{}, context map[string]string) interface{} {
    eventMap, ok := event.(map[string]interface{})
    if !ok {
        return map[string]interface{}{"error": "Invalid event data"}
    }
    
    data, exists := eventMap["data"]
    if !exists {
        return map[string]interface{}{"error": "Missing data field"}
    }
    
    dataStr, ok := data.(string)
    if !ok {
        return map[string]interface{}{"error": "Data must be string"}
    }
    
    // 处理数据：转大写并统计单词
    processed := strings.ToUpper(dataStr)
    wordCount := len(strings.Fields(dataStr))
    
    return map[string]interface{}{
        "original": dataStr,
        "processed": processed,
        "wordCount": wordCount,
    }
}
```

## 2. Node.js函数示例

### 异步HTTP请求函数

```javascript
async function fetchData(event, context) {
    const axios = require('axios');
    
    try {
        const url = event.url || 'https://jsonplaceholder.typicode.com/users/1';
        const response = await axios.get(url);
        
        return {
            success: true,
            data: response.data,
            status: response.status
        };
    } catch (error) {
        return {
            success: false,
            error: error.message
        };
    }
}
```

### JSON处理函数

```javascript
function processJson(event, context) {
    const data = event.data;
    
    if (!data) {
        return { error: 'Missing data field' };
    }
    
    return {
        processed: true,
        keys: Object.keys(data),
        values: Object.values(data),
        count: Object.keys(data).length,
        timestamp: new Date().toISOString()
    };
}
```

## 3. Python函数示例

### 文本分析函数

```python
import re
from collections import Counter

def analyze_text(event, context):
    text = event.get('text', '')
    
    if not text:
        return {'error': 'Missing text field'}
    
    # 基本统计
    char_count = len(text)
    word_count = len(text.split())
    line_count = len(text.split('\n'))
    
    # 单词频率
    words = re.findall(r'\w+', text.lower())
    word_freq = dict(Counter(words).most_common(10))
    
    return {
        'stats': {
            'characters': char_count,
            'words': word_count,
            'lines': line_count
        },
        'top_words': word_freq
    }
```

### 数学计算函数

```python
import math

def calculate(event, context):
    operation = event.get('operation')
    numbers = event.get('numbers', [])
    
    if not operation or not numbers:
        return {'error': 'Missing operation or numbers'}
    
    try:
        if operation == 'sum':
            result = sum(numbers)
        elif operation == 'average':
            result = sum(numbers) / len(numbers)
        elif operation == 'max':
            result = max(numbers)
        elif operation == 'min':
            result = min(numbers)
        elif operation == 'sqrt':
            result = [math.sqrt(x) for x in numbers]
        else:
            return {'error': f'Unsupported operation: {operation}'}
        
        return {
            'operation': operation,
            'input': numbers,
            'result': result
        }
    except Exception as e:
        return {'error': str(e)}
```

## 4. API使用示例

### 创建函数

```bash
curl -X POST http://localhost:8080/api/v1/functions \
  -H "Content-Type: application/json" \
  -d '{
    "name": "hello-world",
    "runtime": "go",
    "handler": "handler",
    "code": "func handler(event interface{}, context map[string]string) interface{} {\n    return map[string]interface{}{\n        \"message\": \"Hello World!\",\n        \"timestamp\": \"2024-01-01T00:00:00Z\"\n    }\n}",
    "timeout": 30,
    "memory": 128
  }'
```

### 调用函数

```bash
curl -X POST http://localhost:8080/api/v1/functions/{function_id}/invoke \
  -H "Content-Type: application/json" \
  -d '{
    "event": {
      "name": "World",
      "type": "greeting"
    },
    "context": {
      "user_id": "123",
      "request_id": "req_456"
    }
  }'
```

### 列出所有函数

```bash
curl http://localhost:8080/api/v1/functions
```

### 获取函数详情

```bash
curl http://localhost:8080/api/v1/functions/{function_id}
```

### 更新函数

```bash
curl -X PUT http://localhost:8080/api/v1/functions/{function_id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "updated-hello",
    "timeout": 60
  }'
```

### 删除函数

```bash
curl -X DELETE http://localhost:8080/api/v1/functions/{function_id}
```

## 5. 响应格式

### 成功响应

```json
{
  "success": true,
  "result": {
    "message": "Hello World!",
    "timestamp": "2024-01-01T00:00:00Z"
  },
  "duration": 123
}
```

### 错误响应

```json
{
  "success": false,
  "error": "函数执行失败: compilation error",
  "duration": 45
}
```

## 6. 注意事项

1. **运行时要求**：
   - Go: 需要安装Go编译器
   - Node.js: 需要安装Node.js运行时
   - Python: 需要安装Python3

2. **安全限制**：
   - 函数在沙盒环境中运行
   - 有超时限制（默认30秒）
   - 有内存限制（默认128MB）

3. **最佳实践**：
   - 保持函数简单和专一
   - 使用环境变量传递配置
   - 处理错误情况
   - 合理设置超时时间

4. **性能优化**：
   - 避免在函数中执行耗时操作
   - 复用连接和资源
   - 使用合适的内存设置 