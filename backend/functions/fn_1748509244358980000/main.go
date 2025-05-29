
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func handler(event interface{}, context map[string]string) interface{} {
	return map[string]interface{}{
		"message": "Hello from Go Cloud Function!",
		"event": event,
		"user": context["user"]   }
}


func main() {
	// 从环境变量读取输入
	eventStr := os.Getenv("FUNCTION_EVENT")
	contextStr := os.Getenv("FUNCTION_CONTEXT")
	
	var event interface{}
	var context map[string]string
	
	if eventStr != "" {
		json.Unmarshal([]byte(eventStr), &event)
	}
	if contextStr != "" {
		json.Unmarshal([]byte(contextStr), &context)
	}
	
	// 调用用户函数
	result := sum(event, context)
	
	// 输出结果
	resultBytes, _ := json.Marshal(result)
	fmt.Print(string(resultBytes))
}
