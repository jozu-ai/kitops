import { URL, fileURLToPath } from 'node:url'

import vue from '@vitejs/plugin-vue'
import icons from 'unplugin-icons/vite'
import { defineConfig } from 'vite'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    icons({ autoInstall: true })
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      'node-fetch': 'just-use-native-fetch',
    }
  }
})
