package cloudfunction

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// executeGoFunction 执行Go函数
func (p *Platform) executeGoFunction(fn *Function, req *ExecuteRequest) (interface{}, error) {
	fnDir := filepath.Join(p.workDir, fn.ID)

	// 分析用户代码，检查是否已包含package和import
	userCode := fn.Code
	hasPackage := false
	hasImports := false

	// 检查是否已有package声明
	if strings.Contains(userCode, "package ") {
		hasPackage = true
	}

	// 检查是否已有import声明
	if strings.Contains(userCode, "import ") {
		hasImports = true
	}

	// 调试输出
	fmt.Printf("用户代码包含package: %v, 包含import: %v\n", hasPackage, hasImports)
	fmt.Printf("用户代码内容: %s\n", userCode)

	// 生成完整的Go程序
	mainCode := fmt.Sprintf(`
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

%s

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
				"error": fmt.Sprintf("函数执行出错: %%v", r),
			}
			resultBytes, _ := json.Marshal(errorResult)
			fmt.Print(string(resultBytes))
		}
	}()
	
	// 直接调用用户的Handler函数
	result := %s(ctx, eventData)
	
	// 包装结果
	finalResult := map[string]interface{}{
		"success": true,
		"result": result,
	}
	
	// 输出结果
	resultBytes, _ := json.Marshal(finalResult)
	fmt.Print(string(resultBytes))
}
`, userCode, fn.Handler)

	mainPath := filepath.Join(fnDir, "main.go")
	if err := os.WriteFile(mainPath, []byte(mainCode), 0644); err != nil {
		return nil, fmt.Errorf("写入main.go失败: %v", err)
	}

	// 准备环境变量
	env := os.Environ()
	if req.Event != nil {
		eventBytes, _ := json.Marshal(req.Event)
		env = append(env, "FUNCTION_EVENT="+string(eventBytes))
	}
	if req.Context != nil {
		contextBytes, _ := json.Marshal(req.Context)
		env = append(env, "FUNCTION_CONTEXT="+string(contextBytes))
	}

	// 设置函数的环境变量
	for k, v := range fn.Environment {
		env = append(env, k+"="+v)
	}

	// 编译并执行
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(fn.Timeout)*time.Second)
	defer cancel()

	// 编译
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", "function", "main.go")
	buildCmd.Dir = fnDir
	buildCmd.Env = env
	if err := buildCmd.Run(); err != nil {
		return nil, fmt.Errorf("编译失败: %v", err)
	}

	// 执行
	runCmd := exec.CommandContext(ctx, "./function")
	runCmd.Dir = fnDir
	runCmd.Env = env
	output, err := runCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("执行失败: %v", err)
	}

	// 解析结果
	var result interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return string(output), nil // 如果不是JSON，返回原始字符串
	}

	return result, nil
}

// executeNodeJSFunction 执行Node.js函数
func (p *Platform) executeNodeJSFunction(fn *Function, req *ExecuteRequest) (interface{}, error) {
	fnDir := filepath.Join(p.workDir, fn.ID)

	// 创建Node.js执行文件
	nodeCode := fmt.Sprintf(`
%s

// 从环境变量读取输入
const eventStr = process.env.FUNCTION_EVENT || '{}';
const contextStr = process.env.FUNCTION_CONTEXT || '{}';

let event, context;

try {
    event = JSON.parse(eventStr);
    context = JSON.parse(contextStr);
} catch (e) {
    console.error(JSON.stringify({error: "解析输入数据失败: " + e.message}));
    process.exit(1);
}

// 执行用户函数
async function execute() {
    try {
        const userFunction = %s;
        
        let result;
        
        // 尝试不同的函数调用方式
        if (typeof userFunction === 'function') {
            if (userFunction.constructor.name === 'AsyncFunction') {
                result = await userFunction(event, context);
            } else {
                result = userFunction(event, context);
            }
        } else {
            throw new Error('Handler不是一个有效的函数');
        }
        
        console.log(JSON.stringify({
            success: true,
            result: result
        }));
    } catch (error) {
        console.error(JSON.stringify({
            success: false,
            error: error.message,
            stack: error.stack
        }));
        process.exit(1);
    }
}

execute();
`, fn.Code, fn.Handler)

	indexPath := filepath.Join(fnDir, "index.js")
	if err := os.WriteFile(indexPath, []byte(nodeCode), 0644); err != nil {
		return nil, fmt.Errorf("写入index.js失败: %v", err)
	}

	// 准备环境变量
	env := os.Environ()
	if req.Event != nil {
		eventBytes, _ := json.Marshal(req.Event)
		env = append(env, "FUNCTION_EVENT="+string(eventBytes))
	}
	if req.Context != nil {
		contextBytes, _ := json.Marshal(req.Context)
		env = append(env, "FUNCTION_CONTEXT="+string(contextBytes))
	}

	// 设置函数的环境变量
	for k, v := range fn.Environment {
		env = append(env, k+"="+v)
	}

	// 执行Node.js
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(fn.Timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "node", "index.js")
	cmd.Dir = fnDir
	cmd.Env = env
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("执行失败: %v", err)
	}

	// 解析结果
	var result interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return string(output), nil
	}

	return result, nil
}

// executePythonFunction 执行Python函数
func (p *Platform) executePythonFunction(fn *Function, req *ExecuteRequest) (interface{}, error) {
	fnDir := filepath.Join(p.workDir, fn.ID)

	// 创建Python执行文件
	pythonCode := fmt.Sprintf(`
import json
import os
import sys
import traceback

%s

def main():
    try:
        # 从环境变量读取输入
        event_str = os.environ.get('FUNCTION_EVENT', '{}')
        context_str = os.environ.get('FUNCTION_CONTEXT', '{}')
        
        try:
            event = json.loads(event_str)
            context = json.loads(context_str)
        except json.JSONDecodeError as e:
            print(json.dumps({
                "success": False,
                "error": f"解析输入数据失败: {str(e)}"
            }))
            sys.exit(1)
        
        # 获取用户函数
        handler_func = globals().get('%s')
        if not handler_func:
            raise Exception(f"找不到Handler函数: %s")
        
        if not callable(handler_func):
            raise Exception(f"Handler不是一个可调用的函数: %s")
        
        # 执行用户函数
        result = handler_func(event, context)
        
        print(json.dumps({
            "success": True,
            "result": result
        }))
        
    except Exception as e:
        print(json.dumps({
            "success": False,
            "error": str(e),
            "traceback": traceback.format_exc()
        }))
        sys.exit(1)

if __name__ == "__main__":
    main()
`, fn.Code, fn.Handler, fn.Handler, fn.Handler)

	mainPath := filepath.Join(fnDir, "main.py")
	if err := os.WriteFile(mainPath, []byte(pythonCode), 0644); err != nil {
		return nil, fmt.Errorf("写入main.py失败: %v", err)
	}

	// 准备环境变量
	env := os.Environ()
	if req.Event != nil {
		eventBytes, _ := json.Marshal(req.Event)
		env = append(env, "FUNCTION_EVENT="+string(eventBytes))
	}
	if req.Context != nil {
		contextBytes, _ := json.Marshal(req.Context)
		env = append(env, "FUNCTION_CONTEXT="+string(contextBytes))
	}

	// 设置函数的环境变量
	for k, v := range fn.Environment {
		env = append(env, k+"="+v)
	}

	// 执行Python
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(fn.Timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "python3", "main.py")
	cmd.Dir = fnDir
	cmd.Env = env

	// 设置进程组，便于杀死子进程
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("执行失败: %v", err)
	}

	// 解析结果
	var result interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return string(output), nil
	}

	return result, nil
}
