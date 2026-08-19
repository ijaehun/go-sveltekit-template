import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	server: {
		// dev 전용: /api 요청을 Go 백엔드로 프록시 (프로덕션은 Caddy가 담당)
		// docker dev에서는 BACKEND_ORIGIN=http://backend:8080 로 override
		proxy: {
			'/api': {
				target: process.env.BACKEND_ORIGIN || 'http://localhost:8080',
				changeOrigin: true
			}
		}
	}
});
