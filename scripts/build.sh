#!/bin/bash

# 云函数平台构建脚本

set -e

echo "🚀 开始构建云函数平台..."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ Go未安装，请先安装Go"
    exit 1
fi

# 进入后端目录
cd backend

echo "📦 下载依赖..."
go mod tidy
go mod download

echo "🔍 运行代码检查..."
if command -v golangci-lint &> /dev/null; then
    golangci-lint run
else
    echo "⚠️  golangci-lint未安装，跳过代码检查"
fi

echo "🧪 运行测试..."
go test ./... -v

echo "🔨 构建后端..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cloudfunction-server .

echo "✅ 后端构建完成"

# 检查前端
cd ../frontend

if [ -f "package.json" ]; then
    echo "📦 安装前端依赖..."
    npm install
    
    echo "🔨 构建前端..."
    npm run build
    
    echo "✅ 前端构建完成"
else
    echo "⚠️  未找到package.json，跳过前端构建"
fi

cd ..

echo "🎉 构建完成！"
echo ""
echo "启动方式："
echo "  后端: cd backend && ./cloudfunction-server"
echo "  前端: cd frontend && npm run dev" 