# 云函数管理平台前端

这是云函数平台的前端管理界面，基于 Vue 3 + Vite 构建。

## 功能特性

- 🔧 云函数创建与管理
- 📝 可视化代码编辑器
- 🎯 一键函数测试
- 📊 执行结果展示
- 🔄 实时状态监控

## 技术栈

- Vue 3 (使用 `<script setup>` 语法)
- Vite
- Axios (HTTP客户端)
- Tailwind CSS (样式框架)

## 开发

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

## 组件说明

- `CloudFunction.vue` - 云函数管理主组件
  - 函数创建表单
  - 函数列表展示
  - 函数测试对话框

确保后端云函数平台已启动并运行在 `http://localhost:8080`.
