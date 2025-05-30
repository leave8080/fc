<template>
  <div class="cloud-function-manager">
    <!-- 创建函数区域 -->
    <div class="mb-16">
      <div class="glass-card p-8 rounded-3xl backdrop-blur-xl bg-white/70 dark:bg-gray-900/70 shadow-2xl border border-white/20 dark:border-gray-700/20">
        <div class="flex items-center mb-8">
          <div class="p-3 bg-gradient-to-r from-blue-500 to-purple-600 rounded-2xl mr-4">
            <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
            </svg>
          </div>
          <h2 class="text-3xl font-bold text-gray-900 dark:text-white">创建新函数</h2>
        </div>
        
        <form @submit.prevent="createFunction" class="space-y-8">
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
            <!-- 函数名称 -->
            <div class="form-group">
              <label class="form-label">函数名称</label>
              <input
                v-model="newFunction.name"
                type="text"
                required
                class="form-input"
                placeholder="输入函数名称"
              />
            </div>

            <!-- 运行时 -->
            <div class="form-group">
              <label class="form-label">运行时</label>
              <select
                v-model="newFunction.runtime"
                required
                class="form-input"
                @change="loadExampleCode"
              >
                <option value="">选择运行时</option>
                <option value="go">Go</option>
                <option value="nodejs">Node.js</option>
                <option value="python">Python</option>
              </select>
            </div>

            <!-- 入口函数 -->
            <div class="form-group">
              <label class="form-label">入口函数</label>
              <input
                v-model="newFunction.handler"
                type="text"
                required
                class="form-input"
                placeholder="handler"
              />
            </div>

            <!-- 超时时间 -->
            <div class="form-group">
              <label class="form-label">超时时间(秒)</label>
              <input
                v-model="newFunction.timeout"
                type="number"
                min="1"
                max="300"
                required
                class="form-input"
                placeholder="30"
              />
            </div>
          </div>

          <!-- 代码编辑器 -->
          <div class="form-group">
            <label class="form-label">函数代码</label>
            <textarea
              v-model="newFunction.code"
              required
              rows="12"
              class="form-textarea"
              placeholder="输入函数代码..."
            ></textarea>
          </div>

          <!-- 提交按钮 -->
          <div class="flex justify-end space-x-4">
            <button
              type="button"
              @click="resetForm"
              class="btn-secondary"
            >
              重置
            </button>
            <button
              type="submit"
              :disabled="creating"
              class="btn-primary"
            >
              <span v-if="creating" class="flex items-center">
                <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                创建中...
              </span>
              <span v-else>创建函数</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 函数列表区域 -->
    <div>
      <div class="flex items-center justify-between mb-8">
        <h2 class="text-3xl font-bold text-gray-900 dark:text-white">函数列表</h2>
        <button @click="loadFunctions" class="btn-secondary">
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
          </svg>
          刷新
        </button>
      </div>

      <!-- 函数卡片网格 -->
      <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-8">
        <div
          v-for="func in functions"
          :key="func.id"
          class="function-card glass-card p-6 rounded-2xl backdrop-blur-xl bg-white/60 dark:bg-gray-900/60 shadow-xl border border-white/20 dark:border-gray-700/20 hover:shadow-2xl transition-all duration-300 transform hover:-translate-y-1"
        >
          <!-- 函数头部 -->
          <div class="flex items-start justify-between mb-4">
            <div class="flex items-center">
              <div class="p-2 bg-gradient-to-r from-emerald-500 to-blue-600 rounded-xl mr-3">
                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"></path>
                </svg>
              </div>
              <div>
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ func.name }}</h3>
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300">
                  {{ func.runtime }}
                </span>
              </div>
            </div>
            <div class="flex space-x-2">
              <button
                @click="openTestDialog(func)"
                class="p-2 text-blue-600 hover:bg-blue-100 dark:hover:bg-blue-900/50 rounded-lg transition-colors"
                title="测试函数"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1m4 0h1m-6 4h4m6-4v6a2 2 0 01-2 2H7a2 2 0 01-2-2V10a2 2 0 012-2h10a2 2 0 012 2z"></path>
                </svg>
              </button>
              <button
                @click="deleteFunction(func.id)"
                class="p-2 text-red-600 hover:bg-red-100 dark:hover:bg-red-900/50 rounded-lg transition-colors"
                title="删除函数"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                </svg>
              </button>
            </div>
          </div>

          <!-- 函数信息 -->
          <div class="space-y-3">
            <div class="flex justify-between text-sm">
              <span class="text-gray-600 dark:text-gray-400">入口函数:</span>
              <span class="font-mono text-gray-900 dark:text-white">{{ func.handler }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600 dark:text-gray-400">超时时间:</span>
              <span class="text-gray-900 dark:text-white">{{ func.timeout }}s</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600 dark:text-gray-400">创建时间:</span>
              <span class="text-gray-900 dark:text-white">{{ formatDate(func.created_at) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="functions.length === 0" class="text-center py-16">
        <div class="mx-auto w-24 h-24 bg-gradient-to-r from-gray-400 to-gray-600 rounded-full flex items-center justify-center mb-4">
          <svg class="w-12 h-12 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"></path>
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">暂无函数</h3>
        <p class="text-gray-600 dark:text-gray-400">创建您的第一个云函数开始使用吧！</p>
      </div>
    </div>

    <!-- 测试对话框 -->
    <div v-if="testDialog.show" class="fixed inset-0 z-50 overflow-y-auto" @click.self="closeTestDialog">
      <div class="flex items-center justify-center min-h-screen px-4">
        <div class="fixed inset-0 bg-black/50 backdrop-blur-sm transition-opacity"></div>
        <div class="relative glass-card bg-white/90 dark:bg-gray-900/90 backdrop-blur-xl rounded-3xl max-w-4xl w-full max-h-[90vh] overflow-hidden shadow-2xl border border-white/20 dark:border-gray-700/20">
          <!-- 对话框头部 -->
          <div class="flex items-center justify-between p-6 border-b border-gray-200/50 dark:border-gray-700/50">
            <h3 class="text-xl font-semibold text-gray-900 dark:text-white">
              测试函数: {{ testDialog.function?.name }}
            </h3>
            <button @click="closeTestDialog" class="p-2 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-xl transition-colors">
              <svg class="w-6 h-6 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </button>
          </div>

          <!-- 对话框内容 -->
          <div class="p-6 space-y-6 max-h-[70vh] overflow-y-auto">
            <!-- 输入事件 -->
            <div>
              <label class="form-label">输入事件 (JSON)</label>
              <textarea
                v-model="testDialog.event"
                rows="8"
                class="form-textarea font-mono"
                placeholder='{"key": "value"}'
              ></textarea>
            </div>

            <!-- 执行结果 -->
            <div v-if="testDialog.result">
              <label class="form-label">执行结果</label>
              <div class="bg-gray-50 dark:bg-gray-800 rounded-xl p-4 font-mono text-sm overflow-x-auto">
                <pre class="text-gray-900 dark:text-white">{{ JSON.stringify(testDialog.result, null, 2) }}</pre>
              </div>
            </div>
          </div>

          <!-- 对话框底部 -->
          <div class="flex justify-end space-x-4 p-6 border-t border-gray-200/50 dark:border-gray-700/50">
            <button @click="closeTestDialog" class="btn-secondary">
              取消
            </button>
            <button
              @click="executeFunction"
              :disabled="testDialog.executing"
              class="btn-primary"
            >
              <span v-if="testDialog.executing" class="flex items-center">
                <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                执行中...
              </span>
              <span v-else>执行函数</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'CloudFunction',
  data() {
    return {
      functions: [],
      creating: false,
      newFunction: {
        name: '',
        runtime: '',
        handler: '',
        code: '',
        timeout: 30
      },
      testDialog: {
        show: false,
        function: null,
        event: '{}',
        result: null,
        executing: false
      },
      examples: {
        go: {
          handler: 'Handler',
          code: `package main

import (
    "context"
    "encoding/json"
    "fmt"
)

type Event struct {
    Name string \`json:"name"\`
}

type Response struct {
    Message string \`json:"message"\`
}

func Handler(ctx context.Context, event Event) (Response, error) {
    return Response{
        Message: fmt.Sprintf("Hello, %s!", event.Name),
    }, nil
}`
        },
        nodejs: {
          handler: 'handler',
          code: `exports.handler = async (event, context) => {
    console.log('Event:', JSON.stringify(event, null, 2));
    
    return {
        statusCode: 200,
        body: {
            message: \`Hello, \${event.name || 'World'}!\`,
            timestamp: new Date().toISOString()
        }
    };
};`
        },
        python: {
          handler: 'lambda_handler',
          code: `import json
import datetime

def lambda_handler(event, context):
    print(f"Event: {json.dumps(event)}")
    
    return {
        'statusCode': 200,
        'body': {
            'message': f"Hello, {event.get('name', 'World')}!",
            'timestamp': datetime.datetime.now().isoformat()
        }
    }`
        }
      }
    }
  },
  
  mounted() {
    this.loadFunctions()
  },
  
  methods: {
    async loadFunctions() {
      try {
        const response = await axios.get('/api/v1/functions')
        this.functions = response.data.functions || []
      } catch (error) {
        console.error('加载函数列表失败:', error)
      }
    },
    
    async createFunction() {
      this.creating = true
      try {
        const response = await axios.post('/api/v1/functions', this.newFunction)
        console.log('函数创建成功:', response.data)
        this.resetForm()
        await this.loadFunctions()
      } catch (error) {
        console.error('创建函数失败:', error)
        alert('创建函数失败: ' + (error.response?.data?.error || error.message))
      } finally {
        this.creating = false
      }
    },
    
    async deleteFunction(id) {
      if (!confirm('确认删除此函数吗？')) {
        return
      }
      
      try {
        await axios.delete(`/api/v1/functions/${id}`)
        await this.loadFunctions()
      } catch (error) {
        console.error('删除函数失败:', error)
        alert('删除函数失败: ' + (error.response?.data?.error || error.message))
      }
    },
    
    openTestDialog(func) {
      this.testDialog = {
        show: true,
        function: func,
        event: '{\n  "name": "World"\n}',
        result: null,
        executing: false
      }
    },
    
    closeTestDialog() {
      this.testDialog.show = false
      this.testDialog.function = null
      this.testDialog.event = '{}'
      this.testDialog.result = null
      this.testDialog.executing = false
    },
    
    async executeFunction() {
      this.testDialog.executing = true
      this.testDialog.result = null
      
      try {
        let event
        try {
          event = JSON.parse(this.testDialog.event)
        } catch (e) {
          throw new Error('输入的JSON格式不正确')
        }
        
        const response = await axios.post(`/api/v1/functions/${this.testDialog.function.id}/invoke`, {
          event: event,
          context: {
            test: 'true'
          }
        })
        
        this.testDialog.result = response.data
      } catch (error) {
        console.error('执行函数失败:', error)
        this.testDialog.result = {
          error: error.response?.data?.error || error.message
        }
      } finally {
        this.testDialog.executing = false
      }
    },
    
    loadExampleCode() {
      if (this.newFunction.runtime && this.examples[this.newFunction.runtime]) {
        const example = this.examples[this.newFunction.runtime]
        this.newFunction.handler = example.handler
        this.newFunction.code = example.code
      }
    },
    
    resetForm() {
      this.newFunction = {
        name: '',
        runtime: '',
        handler: '',
        code: '',
        timeout: 30
      }
    },
    
    formatDate(dateString) {
      return new Date(dateString).toLocaleString()
    }
  }
}
</script>

<style scoped>
.cloud-function-manager {
  @apply space-y-8;
}

.glass-card {
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
}

.form-group {
  @apply space-y-2;
}

.form-label {
  @apply block text-sm font-semibold text-gray-900 dark:text-white;
}

.form-input {
  @apply w-full px-4 py-3 rounded-xl border-0 bg-white/50 dark:bg-gray-800/50 backdrop-blur-md text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 shadow-lg focus:ring-2 focus:ring-blue-500 focus:bg-white/70 dark:focus:bg-gray-800/70 transition-all duration-200;
}

.form-textarea {
  @apply w-full px-4 py-3 rounded-xl border-0 bg-white/50 dark:bg-gray-800/50 backdrop-blur-md text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 shadow-lg focus:ring-2 focus:ring-blue-500 focus:bg-white/70 dark:focus:bg-gray-800/70 transition-all duration-200 resize-none;
}

.btn-primary {
  @apply px-6 py-3 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white font-semibold rounded-xl shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none;
}

.btn-secondary {
  @apply px-6 py-3 bg-white/50 dark:bg-gray-800/50 backdrop-blur-md text-gray-900 dark:text-white font-semibold rounded-xl shadow-lg hover:shadow-xl hover:bg-white/70 dark:hover:bg-gray-800/70 transform hover:-translate-y-0.5 transition-all duration-200 border border-gray-200/50 dark:border-gray-700/50;
}

.function-card {
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
}

/* 动画 */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.function-card {
  animation: fadeInUp 0.5s ease-out;
}

/* 滚动条样式 */
::-webkit-scrollbar {
  @apply w-2;
}

::-webkit-scrollbar-track {
  @apply bg-transparent;
}

::-webkit-scrollbar-thumb {
  @apply bg-gray-300 dark:bg-gray-600 rounded-full;
}

::-webkit-scrollbar-thumb:hover {
  @apply bg-gray-400 dark:bg-gray-500;
}
</style> 