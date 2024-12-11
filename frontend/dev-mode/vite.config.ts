import { URL, fileURLToPath } from 'node:url'

import vue from '@vitejs/plugin-vue'
import icons from 'unplugin-icons/vite'
import { defineConfig } from 'vite'
import dynamicImport from 'vite-plugin-dynamic-import'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    icons({ autoInstall: true }),
    dynamicImport({
      filter(id: string) {
        if (id.includes('@speed-highlight/core'))
          return true
      }
    })
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      'node-fetch': 'just-use-native-fetch',
    }
  }
})
