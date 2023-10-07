import { resolve } from 'path'
import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  envPrefix: 'APP_',

  resolve:{
    alias: {
      '$app': resolve('src'),
      '$lib': resolve('src/lib'),
    },
  },

  build: {
    sourcemap: true
  },
  server: {
    watch: {
      usePolling: process.env.USE_POLLING ?? false
    }
  }
})
