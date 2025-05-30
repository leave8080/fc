<script setup>
import { ref, onMounted } from 'vue'
import CloudFunction from './components/CloudFunction.vue'

const isDark = ref(false)

const toggleTheme = () => {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
}

onMounted(() => {
  // 检查系统主题偏好
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  isDark.value = prefersDark
  document.documentElement.classList.toggle('dark', prefersDark)
})
</script>

<template>
  <div class="min-h-screen bg-white dark:bg-black transition-all duration-700">
    <!-- 苹果风格的导航栏 -->
    <nav class="fixed top-0 w-full z-50 backdrop-blur-xl bg-white/80 dark:bg-black/80 border-b border-gray-200/20 dark:border-gray-800/30">
      <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-[52px]">
          <div class="flex items-center">
            <div class="text-lg font-semibold text-gray-900 dark:text-white tracking-tight">
              CloudFunc
            </div>
          </div>
          
          <!-- 导航菜单 -->
          <div class="hidden md:flex items-center space-x-8">
            <a href="#" class="text-sm text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white transition-colors">产品</a>
            <a href="#" class="text-sm text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white transition-colors">文档</a>
            <a href="#" class="text-sm text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white transition-colors">支持</a>
          </div>
          
          <button
            @click="toggleTheme"
            class="p-2 rounded-full hover:bg-gray-100/80 dark:hover:bg-gray-800/80 transition-all duration-200"
          >
            <svg v-if="isDark" class="w-5 h-5 text-yellow-500" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 2a1 1 0 011 1v1a1 1 0 11-2 0V3a1 1 0 011-1zm4 8a4 4 0 11-8 0 4 4 0 018 0zm-.464 4.95l.707.707a1 1 0 001.414-1.414l-.707-.707a1 1 0 00-1.414 1.414zm2.12-10.607a1 1 0 010 1.414l-.706.707a1 1 0 11-1.414-1.414l.707-.707a1 1 0 011.414 0zM17 11a1 1 0 100-2h-1a1 1 0 100 2h1zm-7 4a1 1 0 011 1v1a1 1 0 11-2 0v-1a1 1 0 011-1zM5.05 6.464A1 1 0 106.465 5.05l-.708-.707a1 1 0 00-1.414 1.414l.707.707zm1.414 8.486l-.707.707a1 1 0 01-1.414-1.414l.707-.707a1 1 0 011.414 1.414zM4 11a1 1 0 100-2H3a1 1 0 000 2h1z" clip-rule="evenodd" />
            </svg>
            <svg v-else class="w-5 h-5 text-gray-700" fill="currentColor" viewBox="0 0 20 20">
              <path d="M17.293 13.293A8 8 0 016.707 2.707a8.001 8.001 0 1010.586 10.586z" />
            </svg>
          </button>
        </div>
      </div>
    </nav>

    <!-- 主要内容区域 -->
    <main class="relative pt-[52px]">
      <!-- 英雄区域 - 苹果风格 -->
      <section class="relative px-4 pt-20 pb-16 lg:pt-32 lg:pb-24">
        <div class="max-w-4xl mx-auto text-center">
          <!-- 主标题 -->
          <h1 class="text-4xl sm:text-6xl lg:text-7xl font-semibold tracking-tight text-gray-900 dark:text-white mb-6 leading-tight">
            <span class="block">云函数</span>
            <span class="block bg-gradient-to-r from-blue-600 via-purple-600 to-indigo-600 bg-clip-text text-transparent">
              重新定义
            </span>
          </h1>
          
          <!-- 副标题 -->
          <p class="text-xl sm:text-2xl text-gray-600 dark:text-gray-300 mb-12 max-w-3xl mx-auto font-light leading-relaxed">
            简单至极，强大无比。<br class="hidden sm:block">
            让云计算如此优雅。
          </p>
          
          <!-- CTA按钮 -->
          <div class="flex flex-col sm:flex-row gap-4 justify-center mb-16">
            <button class="apple-cta-primary">
              开始体验
            </button>
            <button class="apple-cta-secondary">
              了解更多
            </button>
          </div>
          
          <!-- 产品预览图/动画区域 -->
          <div class="relative max-w-5xl mx-auto">
            <div class="aspect-[16/10] bg-gradient-to-br from-gray-100 to-gray-200 dark:from-gray-800 dark:to-gray-900 rounded-[2rem] shadow-2xl relative overflow-hidden">
              <!-- 模拟产品界面 -->
              <div class="absolute inset-4 bg-white dark:bg-gray-900 rounded-[1.5rem] shadow-inner">
                <div class="h-full flex items-center justify-center">
                  <div class="text-center">
                    <div class="w-16 h-16 mx-auto mb-4 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center">
                      <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
                      </svg>
                    </div>
                    <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">极速部署</h3>
                    <p class="text-sm text-gray-500 dark:text-gray-400">毫秒级响应，秒级部署</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 内容区域 -->
      <section class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 pb-24">
        <CloudFunction />
      </section>
    </main>
  </div>
</template>

<style scoped>
/* 苹果风格的CTA按钮 */
.apple-cta-primary {
  @apply px-8 py-4 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-full text-lg transition-all duration-300 transform hover:scale-105 shadow-lg hover:shadow-xl;
}

.apple-cta-secondary {
  @apply px-8 py-4 border-2 border-blue-600 text-blue-600 dark:text-blue-400 dark:border-blue-400 hover:bg-blue-600 hover:text-white dark:hover:bg-blue-400 dark:hover:text-gray-900 font-medium rounded-full text-lg transition-all duration-300 transform hover:scale-105;
}

/* 自定义滚动条 */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  @apply bg-gray-100 dark:bg-gray-800;
}

::-webkit-scrollbar-thumb {
  @apply bg-gray-300 dark:bg-gray-600 rounded-full;
}

::-webkit-scrollbar-thumb:hover {
  @apply bg-gray-400 dark:bg-gray-500;
}
</style>
