import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      '/api/homes': 'http://localhost:8081',
      '/api/rooms': 'http://localhost:8082',
      '/api/objects': 'http://localhost:8080',
      '/api/users': 'http://localhost:8083'
    }
  }
})