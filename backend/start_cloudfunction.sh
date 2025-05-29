#!/bin/bash

echo "=== 启动云函数平台 ==="

# 检查依赖
echo "检查运行时依赖..."

# 检查Go
if ! command -v go &> /dev/null; then
    echo "❌ Go编译器未安装，请先安装Go"
    exit 1
fi
echo "✅ Go编译器已安装: $(go version)"

# 检查Node.js
if ! command -v node &> /dev/null; then
    echo "⚠️  Node.js未安装，无法运行Node.js函数"
else
    echo "✅ Node.js已安装: $(node --version)"
fi

# 检查Python
if ! command -v python3 &> /dev/null; then
    echo "⚠️  Python3未安装，无法运行Python函数"
else
    echo "✅ Python3已安装: $(python3 --version)"
fi

# 创建工作目录
echo "创建工作目录..."
mkdir -p functions

# 安装Go依赖
echo "安装Go依赖..."
go mod tidy

# 启动服务器
echo "启动云函数平台服务器..."
echo "服务地址: http://localhost:8080"
echo "健康检查: http://localhost:8080/api/health"
echo "云函数API: http://localhost:8080/api/v1/functions"
echo ""
echo "按Ctrl+C停止服务"

go run main.go 