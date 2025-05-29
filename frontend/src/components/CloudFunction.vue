<template>
  <div class="cloud-function-manager">
    <div class="container mx-auto p-6">
      <h1 class="text-3xl font-bold mb-8 text-gray-800">云函数管理平台</h1>
      
      <!-- 创建函数表单 -->
      <div class="bg-white rounded-lg shadow-lg p-6 mb-8">
        <h2 class="text-xl font-semibold mb-4">创建新函数</h2>
        <form @submit.prevent="createFunction" class="space-y-4">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">函数名称</label>
              <input
                v-model="newFunction.name"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="输入函数名称"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">运行时</label>
              <select
                v-model="newFunction.runtime"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="">选择运行时</option>
                <option value="go">Go</option>
                <option value="nodejs">Node.js</option>
                <option value="python">Python</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">入口函数</label>
              <input
                v-model="newFunction.handler"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="handler"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">超时时间(秒)</label>
              <input
                v-model="newFunction.timeout"
                type="number"
                min="1"
                max="300"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="30"
              />
            </div>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">函数代码</label>
            <textarea
              v-model="newFunction.code"
              required
              rows="8"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
              placeholder="输入函数代码..."
            ></textarea>
          </div>
          
          <div class="flex space-x-4">
            <button
              type="submit"
              :disabled="creating"
              class="px-6 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
            >
              {{ creating ? '创建中...' : '创建函数' }}
            </button>
            <button
              type="button"
              @click="loadExampleCode"
              class="px-6 py-2 bg-gray-500 text-white rounded-md hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-gray-500"
            >
              加载示例
            </button>
          </div>
        </form>
      </div>

      <!-- 函数列表 -->
      <div class="bg-white rounded-lg shadow-lg p-6 mb-8">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-semibold">函数列表</h2>
          <button
            @click="loadFunctions"
            class="px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-500"
          >
            刷新
          </button>
        </div>
        
        <div v-if="loading" class="text-center py-8">
          <div class="text-gray-500">加载中...</div>
        </div>
        
        <div v-else-if="functions.length === 0" class="text-center py-8">
          <div class="text-gray-500">暂无函数</div>
        </div>
        
        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div
            v-for="func in functions"
            :key="func.id"
            class="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow"
          >
            <div class="flex justify-between items-start mb-2">
              <h3 class="font-semibold text-lg">{{ func.name }}</h3>
              <span class="px-2 py-1 bg-blue-100 text-blue-800 rounded text-xs">
                {{ func.runtime }}
              </span>
            </div>
            
            <div class="text-sm text-gray-600 mb-3">
              <div>入口: {{ func.handler }}</div>
              <div>超时: {{ func.timeout }}s</div>
              <div>创建时间: {{ formatDate(func.created_at) }}</div>
            </div>
            
            <div class="flex space-x-2">
              <button
                @click="testFunction(func)"
                class="flex-1 px-3 py-1 bg-green-500 text-white rounded text-sm hover:bg-green-600"
              >
                测试
              </button>
              <button
                @click="deleteFunction(func.id)"
                class="flex-1 px-3 py-1 bg-red-500 text-white rounded text-sm hover:bg-red-600"
              >
                删除
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 测试函数对话框 -->
      <div v-if="testDialog.show" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg p-6 w-full max-w-2xl max-h-96 overflow-y-auto">
          <h3 class="text-lg font-semibold mb-4">测试函数: {{ testDialog.function.name }}</h3>
          
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-1">输入数据 (JSON)</label>
            <textarea
              v-model="testDialog.input"
              rows="6"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
              placeholder='{"key": "value"}'
            ></textarea>
          </div>
          
          <div v-if="testDialog.result" class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-1">执行结果</label>
            <pre class="bg-gray-100 p-3 rounded-md text-sm overflow-x-auto">{{ JSON.stringify(testDialog.result, null, 2) }}</pre>
          </div>
          
          <div class="flex space-x-4">
            <button
              @click="executeTest"
              :disabled="testDialog.executing"
              class="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 disabled:opacity-50"
            >
              {{ testDialog.executing ? '执行中...' : '执行' }}
            </button>
            <button
              @click="testDialog.show = false"
              class="px-4 py-2 bg-gray-500 text-white rounded-md hover:bg-gray-600"
            >
              关闭
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
      loading: false,
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
        function: {},
        input: '{"message": "Hello World"}',
        result: null,
        executing: false
      },
      examples: {
        go: {
          handler: 'handler',
          code: `func handler(event interface{}, context map[string]string) interface{} {
    return map[string]interface{}{
        "message": "Hello from Go Cloud Function!",
        "event": event,
        "timestamp": time.Now().Format(time.RFC3339),
    }
}`
        },
        nodejs: {
          handler: 'handler',
          code: `function handler(event, context) {
    return {
        message: "Hello from Node.js Cloud Function!",
        event: event,
        timestamp: new Date().toISOString()
    };
}`
        },
        python: {
          handler: 'handler',
          code: `def handler(event, context):
    import datetime
    return {
        'message': 'Hello from Python Cloud Function!',
        'event': event,
        'timestamp': datetime.datetime.now().isoformat()
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
      this.loading = true
      try {
        const response = await axios.get('/api/v1/functions')
        this.functions = response.data.functions || []
      } catch (error) {
        console.error('加载函数列表失败:', error)
        alert('加载函数列表失败')
      } finally {
        this.loading = false
      }
    },
    
    async createFunction() {
      this.creating = true
      try {
        await axios.post('/api/v1/functions', this.newFunction)
        alert('函数创建成功!')
        this.resetForm()
        this.loadFunctions()
      } catch (error) {
        console.error('创建函数失败:', error)
        alert('创建函数失败: ' + (error.response?.data?.error || error.message))
      } finally {
        this.creating = false
      }
    },
    
    async deleteFunction(id) {
      if (!confirm('确定要删除这个函数吗？')) return
      
      try {
        await axios.delete(`/api/v1/functions/${id}`)
        alert('函数删除成功!')
        this.loadFunctions()
      } catch (error) {
        console.error('删除函数失败:', error)
        alert('删除函数失败: ' + (error.response?.data?.error || error.message))
      }
    },
    
    testFunction(func) {
      this.testDialog.function = func
      this.testDialog.show = true
      this.testDialog.result = null
    },
    
    async executeTest() {
      this.testDialog.executing = true
      try {
        let event = {}
        try {
          event = JSON.parse(this.testDialog.input)
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
  min-height: 100vh;
  background-color: #f3f4f6;
}
</style> 