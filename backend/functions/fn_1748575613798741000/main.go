package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

func Handler(ctx context.Context, event interface{}) interface{} {
	// 处理事件数据
	if event == nil {
		return map[string]interface{}{
			"message": "Hello, World!",
		}
	}

	// 尝试获取name字段
	if eventMap, ok := event.(map[string]interface{}); ok {
		if name, exists := eventMap["name"]; exists {
			return map[string]interface{}{
				"message": fmt.Sprintf("Hello, %v!", name),
			}
		}
	}

	return map[string]interface{}{
		"message":  "Hello, World!",
		"received": event,
	}
}

func main() {
	// 从环境变量读取输入
	eventStr := os.Getenv("FUNCTION_EVENT")
	contextStr := os.Getenv("FUNCTION_CONTEXT")

	var eventData interface{}
	var contextMap map[string]string
	ctx := context.Background()

	if eventStr != "" {
		json.Unmarshal([]byte(eventStr), &eventData)
	}
	if contextStr != "" {
		json.Unmarshal([]byte(contextStr), &contextMap)
	}

	// 调用用户函数
	defer func() {
		if r := recover(); r != nil {
			errorResult := map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("函数执行出错: %v", r),
			}
			resultBytes, _ := json.Marshal(errorResult)
			fmt.Print(string(resultBytes))
		}
	}()

	// 直接调用用户的Handler函数
	result := Handler(ctx, eventData)

	// 包装结果
	finalResult := map[string]interface{}{
		"success": true,
		"result":  result,
	}

	// 输出结果
	resultBytes, _ := json.Marshal(finalResult)
	fmt.Print(string(resultBytes))
}
