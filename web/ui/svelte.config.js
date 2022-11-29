import adapter from '@sveltejs/adapter-static'
import preprocess from 'svelte-preprocess'

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: preprocess(),
  env: {
    dir: process.cwd(),
    publicPrefix: 'PUBLIC_'
  },
  kit: {
    adapter: adapter({
      pages: '../bff/generated',
      assets: '../bff/generated',
    }),
  },
}

export default config
