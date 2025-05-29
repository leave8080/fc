package cloudfunction

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

// executeGoFunction 执行Go函数
func (p *Platform) executeGoFunction(fn *Function, req *ExecuteRequest) (interface{}, error) {
	fnDir := filepath.Join(p.workDir, fn.ID)

	// 创建临时的main.go文件
	mainCode := fmt.Sprintf(`
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

%s

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
	result := %s(event, context)
	
	// 输出结果
	resultBytes, _ := json.Marshal(result)
	fmt.Print(string(resultBytes))
}
`, fn.Code, fn.Handler)

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
const userFunction = %s;

// 从环境变量读取输入
const eventStr = process.env.FUNCTION_EVENT || '{}';
const contextStr = process.env.FUNCTION_CONTEXT || '{}';

const event = JSON.parse(eventStr);
const context = JSON.parse(contextStr);

// 执行用户函数
async function execute() {
    try {
        const result = await userFunction(event, context);
        console.log(JSON.stringify(result));
    } catch (error) {
        console.error(JSON.stringify({error: error.message}));
        process.exit(1);
    }
}

execute();
`, fn.Code)

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

%s

def main():
    try:
        # 从环境变量读取输入
        event_str = os.environ.get('FUNCTION_EVENT', '{}')
        context_str = os.environ.get('FUNCTION_CONTEXT', '{}')
        
        event = json.loads(event_str)
        context = json.loads(context_str)
        
        # 执行用户函数
        result = %s(event, context)
        
        print(json.dumps(result))
    except Exception as e:
        print(json.dumps({"error": str(e)}), file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()
`, fn.Code, fn.Handler)

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
