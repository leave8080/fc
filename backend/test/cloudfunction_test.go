package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080/api/v1"

func main() {
	fmt.Println("=== 云函数平台测试 ===")

	// 等待服务器启动
	fmt.Println("等待服务器启动...")
	time.Sleep(2 * time.Second)

	// 测试健康检查
	testHealth()

	// 测试Go函数
	testGoFunction()

	// 测试Node.js函数
	testNodeJSFunction()

	// 测试Python函数
	testPythonFunction()
}

func testHealth() {
	fmt.Println("\n1. 测试健康检查...")
	resp, err := http.Get("http://localhost:8080/api/health")
	if err != nil {
		fmt.Printf("❌ 健康检查失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("✅ 健康检查通过")
	} else {
		fmt.Printf("❌ 健康检查失败: %d\n", resp.StatusCode)
	}
}

func testGoFunction() {
	fmt.Println("\n2. 测试Go函数...")

	// 创建Go函数
	goFunction := map[string]interface{}{
		"name":    "test-go",
		"runtime": "go",
		"handler": "handler",
		"code": `import "time"

func handler(event interface{}, context map[string]string) interface{} {
    return map[string]interface{}{
        "message": "Hello from Go!",
        "timestamp": time.Now().Format(time.RFC3339),
        "event": event,
    }
}`,
		"timeout": 30,
		"memory":  128,
	}

	functionID := createFunction(goFunction)
	if functionID == "" {
		return
	}

	// 测试函数调用
	testData := map[string]interface{}{
		"event": map[string]interface{}{
			"test": "go function",
			"data": "hello world",
		},
		"context": map[string]string{
			"test": "true",
		},
	}

	invokeFunction(functionID, testData)

	// 清理
	deleteFunction(functionID)
}

func testNodeJSFunction() {
	fmt.Println("\n3. 测试Node.js函数...")

	// 创建Node.js函数
	nodeFunction := map[string]interface{}{
		"name":    "test-nodejs",
		"runtime": "nodejs",
		"handler": "handler",
		"code": `function handler(event, context) {
    return {
        message: "Hello from Node.js!",
        timestamp: new Date().toISOString(),
        event: event
    };
}`,
		"timeout": 30,
		"memory":  128,
	}

	functionID := createFunction(nodeFunction)
	if functionID == "" {
		return
	}

	// 测试函数调用
	testData := map[string]interface{}{
		"event": map[string]interface{}{
			"test": "nodejs function",
			"data": "hello world",
		},
		"context": map[string]string{
			"test": "true",
		},
	}

	invokeFunction(functionID, testData)

	// 清理
	deleteFunction(functionID)
}

func testPythonFunction() {
	fmt.Println("\n4. 测试Python函数...")

	// 创建Python函数
	pythonFunction := map[string]interface{}{
		"name":    "test-python",
		"runtime": "python",
		"handler": "handler",
		"code": `def handler(event, context):
    import datetime
    return {
        'message': 'Hello from Python!',
        'timestamp': datetime.datetime.now().isoformat(),
        'event': event
    }`,
		"timeout": 30,
		"memory":  128,
	}

	functionID := createFunction(pythonFunction)
	if functionID == "" {
		return
	}

	// 测试函数调用
	testData := map[string]interface{}{
		"event": map[string]interface{}{
			"test": "python function",
			"data": "hello world",
		},
		"context": map[string]string{
			"test": "true",
		},
	}

	invokeFunction(functionID, testData)

	// 清理
	deleteFunction(functionID)
}

func createFunction(funcData map[string]interface{}) string {
	jsonData, _ := json.Marshal(funcData)

	resp, err := http.Post(baseURL+"/functions", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ 创建函数失败: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 201 {
		var result map[string]interface{}
		json.Unmarshal(body, &result)

		if function, ok := result["function"].(map[string]interface{}); ok {
			if id, ok := function["id"].(string); ok {
				fmt.Printf("✅ 函数创建成功: %s (ID: %s)\n", funcData["name"], id)
				return id
			}
		}
	} else {
		fmt.Printf("❌ 创建函数失败: %d - %s\n", resp.StatusCode, string(body))
	}

	return ""
}

func invokeFunction(functionID string, testData map[string]interface{}) {
	jsonData, _ := json.Marshal(testData)

	resp, err := http.Post(baseURL+"/functions/"+functionID+"/invoke", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ 调用函数失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 200 {
		var result map[string]interface{}
		json.Unmarshal(body, &result)

		if success, ok := result["success"].(bool); ok && success {
			fmt.Printf("✅ 函数执行成功\n")
			if resultData, ok := result["result"]; ok {
				fmt.Printf("   结果: %v\n", resultData)
			}
			if duration, ok := result["duration"].(float64); ok {
				fmt.Printf("   耗时: %.0fms\n", duration)
			}
		} else {
			fmt.Printf("❌ 函数执行失败: %v\n", result["error"])
		}
	} else {
		fmt.Printf("❌ 调用函数失败: %d - %s\n", resp.StatusCode, string(body))
	}
}

func deleteFunction(functionID string) {
	req, _ := http.NewRequest("DELETE", baseURL+"/functions/"+functionID, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ 删除函数失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Printf("✅ 函数删除成功: %s\n", functionID)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("❌ 删除函数失败: %d - %s\n", resp.StatusCode, string(body))
	}
}
