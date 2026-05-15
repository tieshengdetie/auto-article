import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const backendBaseUrl = env.VITE_BACKEND_BASE_URL || env.BACKEND_BASE_URL || 'http://localhost:9001'

  return {
    plugins: [vue()],
    server: {
      port: 5173,
      proxy: {
        '/api': {
          target: backendBaseUrl,
          changeOrigin: true
        },
        '/static': {
          target: backendBaseUrl,
          changeOrigin: true
        }
      }
    }
  }
})
