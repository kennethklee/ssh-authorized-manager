import { resolve } from 'path'
import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import { optimizeCss, optimizeImports } from 'carbon-preprocess-svelte'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte({ preprocess: [optimizeImports()] }), optimizeCss()],
  envPrefix: 'APP_',

  resolve:{
    alias: {
      '$app': resolve('src'),
      '$lib': resolve('src/lib'),
    },
  },
})
