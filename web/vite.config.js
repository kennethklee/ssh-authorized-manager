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
})
