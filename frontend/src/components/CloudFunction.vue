<template>
  <div class="cloud-function-manager animate-fade-in-up">
    <!-- 创建函数区域 -->
    <div class="mb-16">
      <div class="apple-card p-8">
        <div class="flex items-center mb-8">
          <div class="w-14 h-14 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center mr-4 shadow-lg">
            <svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
            </svg>
          </div>
          <div>
            <h2 class="text-3xl font-semibold text-gray-900 dark:text-white tracking-tight">创建新函数</h2>
            <p class="text-gray-500 dark:text-gray-400 mt-2">快速部署您的云函数，体验极速响应</p>
          </div>
        </div>
        
        <form @submit.prevent="createFunction" class="space-y-8">
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <!-- 函数名称 -->
            <div class="form-group">
              <label class="form-label">函数名称</label>
              <input
                v-model="newFunction.name"
                type="text"
                required
                class="apple-input"
                placeholder="输入函数名称"
              />
            </div>

            <!-- 运行时 -->
            <div class="form-group">
              <label class="form-label">运行时</label>
              <select
                v-model="newFunction.runtime"
                required
                class="apple-input"
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
                class="apple-input"
                placeholder="handler"
              />
            </div>

            <!-- 超时时间 -->
            <div class="form-group">
              <label class="form-label">超时时间 (秒)</label>
              <input
                v-model="newFunction.timeout"
                type="number"
                min="1"
                max="300"
                required
                class="apple-input"
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
              class="apple-textarea font-mono text-sm"
              placeholder="输入函数代码..."
            ></textarea>
          </div>

          <!-- 提交按钮 -->
          <div class="flex justify-end space-x-4 pt-6">
            <button
              type="button"
              @click="resetForm"
              class="apple-btn-secondary"
            >
              重置
            </button>
            <button
              type="submit"
              :disabled="creating"
              class="apple-btn-primary"
            >
              <span v-if="creating" class="flex items-center">
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 714 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
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
      <div class="flex items-center justify-between mb-10">
        <div>
          <h2 class="text-3xl font-semibold text-gray-900 dark:text-white tracking-tight">函数列表</h2>
          <p class="text-gray-500 dark:text-gray-400 mt-2">管理和监控您的云函数</p>
        </div>
        <button @click="loadFunctions" class="apple-btn-secondary flex items-center">
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
          </svg>
          刷新
        </button>
      </div>

      <!-- 函数卡片网格 -->
      <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
        <div
          v-for="func in functions"
          :key="func.id"
          class="function-card group"
        >
          <!-- 函数头部 -->
          <div class="flex items-start justify-between mb-6">
            <div class="flex items-center flex-1 min-w-0">
              <div class="flex-shrink-0 w-12 h-12 bg-gradient-to-br from-emerald-400 to-blue-500 rounded-2xl flex items-center justify-center mr-4 shadow-lg">
                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"></path>
                </svg>
              </div>
              <div class="min-w-0 flex-1">
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white truncate mb-1">{{ func.name }}</h3>
                <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300">
                  {{ func.runtime }}
                </span>
              </div>
            </div>
            <div class="flex space-x-2 ml-3">
              <button
                @click="openTestDialog(func)"
                class="p-2 rounded-lg bg-gray-100 hover:bg-blue-100 dark:bg-gray-800 dark:hover:bg-blue-900/30 text-blue-600 dark:text-blue-400 transition-all duration-200"
                title="测试函数"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1m4 0h1m-6 4h4m6-4v6a2 2 0 01-2 2H7a2 2 0 01-2-2V10a2 2 0 012-2h10a2 2 0 012 2z"></path>
                </svg>
              </button>
              <button
                @click="deleteFunction(func.id)"
                class="p-2 rounded-lg bg-gray-100 hover:bg-red-100 dark:bg-gray-800 dark:hover:bg-red-900/30 text-red-600 dark:text-red-400 transition-all duration-200"
                title="删除函数"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                </svg>
              </button>
            </div>
          </div>

          <!-- 函数信息 -->
          <div class="space-y-4">
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-500 dark:text-gray-400">入口函数</span>
              <span class="font-mono text-sm text-gray-900 dark:text-white bg-gray-100 dark:bg-gray-700/50 px-3 py-1 rounded-lg">{{ func.handler }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-500 dark:text-gray-400">超时时间</span>
              <span class="text-sm font-medium text-gray-900 dark:text-white">{{ func.timeout }}s</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-500 dark:text-gray-400">创建时间</span>
              <span class="text-sm text-gray-900 dark:text-white">{{ formatDate(func.created_at) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="functions.length === 0" class="text-center py-24">
        <div class="mx-auto w-24 h-24 bg-gray-100 dark:bg-gray-800 rounded-full flex items-center justify-center mb-8 shadow-lg">
          <svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"></path>
          </svg>
        </div>
        <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-3">暂无函数</h3>
        <p class="text-gray-500 dark:text-gray-400 text-lg">创建您的第一个云函数开始使用吧</p>
      </div>
    </div>

    <!-- 测试对话框 -->
    <div v-if="testDialog.show" class="fixed inset-0 z-50 overflow-y-auto" @click.self="closeTestDialog">
      <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:p-0">
        <div class="fixed inset-0 bg-black/40 backdrop-blur-sm transition-opacity"></div>
        <div class="relative inline-block w-full max-w-2xl p-6 my-8 text-left align-middle transition-all transform bg-white/95 dark:bg-gray-900/95 backdrop-blur-xl shadow-2xl rounded-2xl border border-gray-200/50 dark:border-gray-700/50">
          <!-- 对话框头部 -->
          <div class="flex items-center justify-between mb-6">
            <div>
              <h3 class="text-xl font-semibold text-gray-900 dark:text-white tracking-tight">
                测试函数
              </h3>
              <p class="text-gray-500 dark:text-gray-400 mt-1">{{ testDialog.function?.name }}</p>
            </div>
            <button @click="closeTestDialog" class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-all duration-200">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </button>
          </div>

          <!-- 对话框内容 -->
          <div class="space-y-6 max-h-[60vh] overflow-y-auto">
            <!-- 输入事件 -->
            <div class="form-group">
              <label class="form-label">输入事件 (JSON)</label>
              <textarea
                v-model="testDialog.event"
                rows="8"
                class="apple-textarea font-mono text-sm"
                placeholder='{"name": "World"}'
              ></textarea>
            </div>

            <!-- 执行结果 -->
            <div v-if="testDialog.result" class="form-group">
              <label class="form-label">执行结果</label>
              <div class="bg-gray-50/50 dark:bg-gray-800/50 rounded-xl p-4 border border-gray-200/50 dark:border-gray-700/50">
                <pre class="font-mono text-sm text-gray-900 dark:text-white whitespace-pre-wrap overflow-x-auto">{{ JSON.stringify(testDialog.result, null, 2) }}</pre>
              </div>
            </div>
          </div>

          <!-- 对话框底部 -->
          <div class="flex justify-end space-x-4 mt-8 pt-6 border-t border-gray-200/50 dark:border-gray-700/50">
            <button @click="closeTestDialog" class="apple-btn-secondary">
              取消
            </button>
            <button
              @click="executeFunction"
              :disabled="testDialog.executing"
              class="apple-btn-primary"
            >
              <span v-if="testDialog.executing" class="flex items-center">
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 714 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                执行中...
              </span>
              <span v-else>执行函数</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 删除对话框 -->
    <div v-if="deleteDialog.show" class="fixed inset-0 z-50 overflow-y-auto" @click.self="closeDeleteDialog" @keydown.esc="closeDeleteDialog" @keydown.enter="confirmDelete">
      <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:p-0">
        <div class="fixed inset-0 bg-black/40 backdrop-blur-sm transition-opacity"></div>
        <div class="relative inline-block w-full max-w-md p-8 my-8 text-center align-middle transition-all transform bg-white/95 dark:bg-gray-900/95 backdrop-blur-xl shadow-2xl rounded-3xl border border-red-200/50 dark:border-red-700/30 animate-scale-in">
          <!-- 警告图标 -->
          <div class="mx-auto mb-6 w-16 h-16 bg-gradient-to-br from-red-400 to-red-600 rounded-full flex items-center justify-center shadow-lg">
            <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
            </svg>
          </div>
          
          <!-- 对话框内容 -->
          <div class="mb-8">
            <h3 class="text-2xl font-semibold text-gray-900 dark:text-white mb-3 tracking-tight">
              删除函数
            </h3>
            <p class="text-gray-600 dark:text-gray-300 text-lg mb-4">
              确定要删除函数 <span class="font-semibold text-gray-900 dark:text-white">{{ deleteDialog.function?.name }}</span> 吗？
            </p>
            <div class="bg-red-50/80 dark:bg-red-900/20 backdrop-blur-sm rounded-2xl p-4 mb-6 border border-red-200/50 dark:border-red-700/30">
              <div class="flex items-center">
                <svg class="w-5 h-5 text-red-500 mr-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
                </svg>
                <p class="text-red-700 dark:text-red-300 text-sm">
                  此操作无法撤销，函数代码和相关数据将被永久删除。
                </p>
              </div>
            </div>
          </div>

          <!-- 对话框底部 -->
          <div class="flex space-x-4">
            <button @click="closeDeleteDialog" class="flex-1 apple-btn-secondary">
              取消
            </button>
            <button
              @click="confirmDelete"
              :disabled="deleteDialog.deleting"
              class="flex-1 px-6 py-3 bg-red-500 hover:bg-red-600 text-white font-medium rounded-xl shadow-lg hover:shadow-xl transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed transform hover:-translate-y-0.5 active:translate-y-0"
            >
              <span v-if="deleteDialog.deleting" class="flex items-center justify-center">
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 714 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                删除中...
              </span>
              <span v-else class="flex items-center justify-center">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                </svg>
                确认删除
              </span>
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
      deleteDialog: {
        show: false,
        function: null,
        deleting: false
      },
      examples: {
        go: {
          handler: 'Handler',
          code: `func Handler(ctx context.Context, event interface{}) interface{} {
    // 处理事件数据
    if event == nil {
        return map[string]interface{}{
            "message": "Hello, World!",
        }
    }
    
    // 尝试获取name字段
    if eventMap, ok := event.(map[string]interface{}); ok {
        if name, exists := eventMap["name"]; exists {
            return map[string]interface{}{
                "message": fmt.Sprintf("Hello, %v!", name),
            }
        }
    }
    
    return map[string]interface{}{
        "message": "Hello, World!",
        "received": event,
    }
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
    document.addEventListener('keydown', this.handleKeyDown)
  },
  
  beforeUnmount() {
    document.removeEventListener('keydown', this.handleKeyDown)
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
      const func = this.functions.find(f => f.id === id)
      if (!func) return
      
      this.deleteDialog = {
        show: true,
        function: func,
        deleting: false
      }
    },
    
    async confirmDelete() {
      if (!this.deleteDialog.function) return
      
      this.deleteDialog.deleting = true
      
      try {
        await axios.delete(`/api/v1/functions/${this.deleteDialog.function.id}`)
        await this.loadFunctions()
        this.closeDeleteDialog()
      } catch (error) {
        console.error('删除函数失败:', error)
        alert('删除函数失败: ' + (error.response?.data?.error || error.message))
      } finally {
        this.deleteDialog.deleting = false
      }
    },
    
    closeDeleteDialog() {
      this.deleteDialog = {
        show: false,
        function: null,
        deleting: false
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
    },
    
    handleKeyDown(event) {
      if (this.deleteDialog.show) {
        if (event.key === 'Escape') {
          event.preventDefault()
          this.closeDeleteDialog()
        } else if (event.key === 'Enter') {
          event.preventDefault()
          if (!this.deleteDialog.deleting) {
            this.confirmDelete()
          }
        }
      }
    }
  }
}
</script>

<style scoped>
.cloud-function-manager {
  @apply space-y-8;
}

/* 苹果风格卡片 */
.apple-card {
  @apply rounded-3xl shadow-xl;
  backdrop-filter: blur(40px);
  -webkit-backdrop-filter: blur(40px);
  border: 1px solid rgba(255, 255, 255, 0.15);
}

/* 苹果风格函数卡片 */
.apple-function-card {
  @apply p-6 rounded-2xl bg-white/80 dark:bg-gray-900/80 backdrop-blur-2xl border border-gray-200/20 dark:border-gray-700/20 shadow-lg hover:shadow-xl transition-all duration-300 ease-in-out hover:-translate-y-1;
  backdrop-filter: blur(40px);
  -webkit-backdrop-filter: blur(40px);
}

/* 苹果风格输入框 */
.apple-input {
  @apply w-full px-4 py-3 rounded-xl border-0 bg-white/60 dark:bg-gray-800/60 backdrop-blur-md text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 shadow-sm focus:ring-2 focus:ring-blue-500/50 focus:bg-white/80 dark:focus:bg-gray-800/80 transition-all duration-200;
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
}

.apple-textarea {
  @apply w-full px-4 py-3 rounded-xl border-0 bg-white/60 dark:bg-gray-800/60 backdrop-blur-md text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 shadow-sm focus:ring-2 focus:ring-blue-500/50 focus:bg-white/80 dark:focus:bg-gray-800/80 transition-all duration-200 resize-none;
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
}

/* 苹果风格按钮 */
.apple-btn-primary {
  @apply px-6 py-3 bg-blue-500 hover:bg-blue-600 text-white font-medium rounded-xl shadow-lg hover:shadow-xl transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed transform hover:-translate-y-0.5 active:translate-y-0;
}

.apple-btn-secondary {
  @apply px-6 py-3 bg-white/70 dark:bg-gray-800/70 backdrop-blur-md text-gray-700 dark:text-gray-300 font-medium rounded-xl shadow-md hover:shadow-lg hover:bg-white/80 dark:hover:bg-gray-800/80 transition-all duration-200 border border-gray-200/50 dark:border-gray-700/50 transform hover:-translate-y-0.5 active:translate-y-0;
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
}

/* 苹果风格图标按钮 */
.apple-icon-btn {
  @apply p-2 rounded-lg backdrop-blur-md transition-all duration-200 hover:bg-white/80 dark:hover:bg-gray-800/80;
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
}

/* 苹果风格模态框 */
.apple-modal {
  @apply bg-white/90 dark:bg-gray-900/90 backdrop-blur-2xl rounded-3xl max-w-4xl w-full max-h-[90vh] overflow-hidden shadow-2xl border border-white/20 dark:border-gray-700/20;
  backdrop-filter: blur(40px);
  -webkit-backdrop-filter: blur(40px);
}

.form-group {
  @apply space-y-2;
}

.form-label {
  @apply block text-sm font-medium text-gray-900 dark:text-white;
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

@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(0.9) translateY(-10px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.animate-fade-in-up {
  animation: fadeInUp 0.6s ease-out;
}

.animate-scale-in {
  animation: scaleIn 0.3s ease-out;
}

.apple-function-card {
  animation: fadeInUp 0.6s ease-out;
}

/* 滚动条样式 */
::-webkit-scrollbar {
  @apply w-2;
}

::-webkit-scrollbar-track {
  @apply bg-transparent;
}

::-webkit-scrollbar-thumb {
  @apply bg-gray-300/60 dark:bg-gray-600/60 rounded-full;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
}

::-webkit-scrollbar-thumb:hover {
  @apply bg-gray-400/60 dark:bg-gray-500/60;
}
</style>