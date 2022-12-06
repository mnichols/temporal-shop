import { sveltekit } from '@sveltejs/kit/vite'
import type { UserConfig } from 'vite'
import path from "path";

const config: UserConfig = {
  base: '/app',
  plugins: [sveltekit()],
  resolve: {
    alias: {
      $lib: path.resolve(__dirname, './src/lib'),
      $types: path.resolve(__dirname, './src/types'),
      $components: path.resolve(__dirname, './src/lib/components/'),
      //$app: path.resolve(__dirname, './src/lib/svelte-mocks/app/'),
      $fixtures: path.resolve(__dirname, './src/fixtures/'),
    },
  },
  build: {
    sourcemap: true
  }
}

export default config
