import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    // 开发服务器配置
    port: 5173,
    // 可选：配置代理（如果不使用nginx反向代理）
    proxy: {
      // 如果想在开发时也使用代理，可以取消注释以下配置
      // '/api': {
      //   target: 'http://localhost:8080',
      //   changeOrigin: true,
      //   secure: false,
      // },
      // '/ws': {
      //   target: 'ws://localhost:8080',
      //   ws: true,
      //   changeOrigin: true,
      // }
    }
  }
})
