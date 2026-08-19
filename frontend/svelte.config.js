import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),
	kit: {
		// adapter-node: 빌드 결과를 node 서버(build/index.js)로 실행 (SSR).
		// SEO가 필요 없는 SPA로 바꾸려면 @sveltejs/adapter-static으로 교체 (README 참고).
		adapter: adapter()
	}
};

export default config;
