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
      // 友链等新上传图在后端 Pictures 目录；public 没有时回源到 Go 静态
      '/Pictures': {
        target: 'http://127.0.0.1:3000',
        changeOrigin: true,
      },
    },
  },
})
