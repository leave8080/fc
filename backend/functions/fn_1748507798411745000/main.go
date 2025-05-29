
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

import "fmt"

func calculate(event interface{}, context map[string]string) interface{} {
    eventMap, ok := event.(map[string]interface{})
    if !ok {
        return map[string]interface{}{
            "error": "Invalid event format",
        }
    }

    operation, hasOp := eventMap["operation"].(string)
    numbers, hasNumbers := eventMap["numbers"].([]interface{})
    
    if !hasOp || !hasNumbers {
        return map[string]interface{}{
            "error": "Missing operation or numbers",
        }
    }

    var nums []float64
    for _, v := range numbers {
        if num, ok := v.(float64); ok {
            nums = append(nums, num)
        }
    }

    var result float64
    switch operation {
    case "sum":
        for _, num := range nums {
            result += num
        }
    case "multiply":
        result = 1
        for _, num := range nums {
            result *= num
        }
    case "average":
        for _, num := range nums {
            result += num
        }
        result = result / float64(len(nums))
    default:
        return map[string]interface{}{
            "error": "Unsupported operation",
        }
    }

    return map[string]interface{}{
        "operation": operation,
        "numbers": nums,
        "result": result,
        "message": fmt.Sprintf("Processed by Go function for user: %s", context["user"]),
    }
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
	result := calculate(event, context)
	
	// 输出结果
	resultBytes, _ := json.Marshal(result)
	fmt.Print(string(resultBytes))
}
