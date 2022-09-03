import { fileURLToPath } from 'url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import unocss from 'unocss/vite'
import { presetAttributify, presetIcons, presetUno } from 'unocss'

// https://vitejs.dev/config/
export default defineConfig({
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      '@wails': fileURLToPath(new URL('./wailsjs', import.meta.url)),
    },
  },
  plugins: [
    vue(),
    unocss({
      presets: [
        presetIcons(), presetAttributify(), presetUno(),
      ],
    }),
  ],
})
