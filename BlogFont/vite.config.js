import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    proxy: {
      // 同源 /api → BlogBack，便于 HttpOnly Cookie
      '/api': {
        target: 'http://127.0.0.1:3000',
        changeOrigin: true,
      },
      // 图片开发期直接用 https://cspona.top/Pictures/（见 src/utils/image.js）
      // 不再代理 /Pictures，避免和 public/Pictures 本地文件打架
    },
  },
})
