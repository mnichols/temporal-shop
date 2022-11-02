// import adapter from '@sveltejs/adapter-auto';
//
// /** @type {import('@sveltejs/kit').Config} */
// const config = {
// 	// Consult https://github.com/sveltejs/svelte-preprocess
// 	// for more information about preprocessors
// 	preprocess: preprocess(),
//
// 	kit: {
// 		adapter: adapter()
// 	}
// };
//
// export default config;

import adapter from '@sveltejs/adapter-static';
import preprocess from 'svelte-preprocess';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: preprocess(),
	kit: {
		adapter: adapter({
			pages: '../bff/generated',
			assets: '../bff/generated',
			//fallback: 'index.html',
		})
	}
};

export default config;