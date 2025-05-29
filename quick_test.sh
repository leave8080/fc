#!/bin/bash

echo "=== 云函数平台快速验证 ==="

# 检查后端目录
if [ ! -d "backend" ]; then
    echo "❌ backend目录不存在"
    exit 1
fi

# 检查前端目录
if [ ! -d "frontend" ]; then
    echo "❌ frontend目录不存在"
    exit 1
fi

# 检查核心文件
echo "检查核心文件..."

required_files=(
    "backend/main.go"
    "backend/cloudfunction/platform.go"
    "backend/cloudfunction/executors.go"
    "backend/cloudfunction/server.go"
    "backend/start_cloudfunction.sh"
    "frontend/src/components/CloudFunction.vue"
    "README.md"
)

for file in "${required_files[@]}"; do
    if [ -f "$file" ]; then
        echo "✅ $file"
    else
        echo "❌ $file 缺失"
        exit 1
    fi
done

echo ""
echo "✅ 所有核心文件检查通过！"
echo ""
echo "🚀 启动说明："
echo "1. 启动后端: cd backend && ./start_cloudfunction.sh"
echo "2. 启动前端: cd frontend && npm install && npm run dev"
echo "3. 访问: http://localhost:3000 (前端) 和 http://localhost:8080 (后端API)"
echo ""
echo "📚 详细文档: backend/README_CloudFunction.md" 