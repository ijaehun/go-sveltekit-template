// API 클라이언트. 항상 상대경로(/api/...)를 쓴다:
//  - dev: vite proxy가 Go 백엔드(8080)로 전달 (vite.config.ts)
//  - prod: Caddy가 /api/* 를 backend 컨테이너로 라우팅
// feature가 늘어나면 services/<feature>.ts 로 나누고 이 함수를 재사용.

export class ApiError extends Error {
	constructor(
		public status: number,
		message: string
	) {
		super(message);
	}
}

export async function api<T>(path: string, init?: RequestInit): Promise<T> {
	const res = await fetch(path, {
		headers: { 'Content-Type': 'application/json' },
		...init
	});
	if (!res.ok) {
		const body = await res.text().catch(() => '');
		throw new ApiError(res.status, body || res.statusText);
	}
	return res.json() as Promise<T>;
}
