
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

import "strings"

func processText(event interface{}, context map[string]string) interface{} {
    eventMap, ok := event.(map[string]interface{})
    if !ok {
        return map[string]interface{}{"error": "Invalid event format"}
    }

    text, hasText := eventMap["text"].(string)
    operation, hasOp := eventMap["operation"].(string)
    
    if !hasText {
        return map[string]interface{}{"error": "Missing text field"}
    }
    
    if !hasOp {
        operation = "info"
    }

    var result interface{}
    switch operation {
    case "upper":
        result = strings.ToUpper(text)
    case "lower":
        result = strings.ToLower(text)
    case "title":
        result = strings.Title(text)
    case "reverse":
        runes := []rune(text)
        for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
            runes[i], runes[j] = runes[j], runes[i]
        }
        result = string(runes)
    case "words":
        result = strings.Fields(text)
    case "info":
        words := strings.Fields(text)
        result = map[string]interface{}{
            "length": len(text),
            "word_count": len(words),
            "has_numbers": strings.ContainsAny(text, "0123456789"),
            "first_word": func() string {
                if len(words) > 0 {
                    return words[0]
                }
                return ""
            }(),
        }
    default:
        return map[string]interface{}{"error": "Unsupported operation"}
    }

    return map[string]interface{}{
        "original": text,
        "operation": operation,
        "result": result,
        "processed_by": context["user"],
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
	result := processText(event, context)
	
	// 输出结果
	resultBytes, _ := json.Marshal(result)
	fmt.Print(string(resultBytes))
}
