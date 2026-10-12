import { fileURLToPath, URL } from 'node:url'

import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

export default defineConfig(({ command, mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const apiTarget = env.VITE_API_TARGET || 'http://localhost:8888'
  // 调试插件只在本地开发启用，生产构建不注入任何调试标记
  const isDev = command === 'serve'

  return {
    plugins: [vue(), isDev ? vueDevTools() : false],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
    server: {
      port: 5173,
      proxy: {
        '/pub': {
          target: apiTarget,
          changeOrigin: true,
        },
        '/pri': {
          target: apiTarget,
          changeOrigin: true,
        },
      },
    },
  }
})
