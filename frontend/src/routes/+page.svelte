<script lang="ts">
	import { api } from '$lib/api';

	// ── Svelte 5 runes 데모 ──────────────────────────────
	// $state: 반응형 상태 (React의 useState와 달리 그냥 변수처럼 읽고 쓴다)
	let hello = $state<string>('불러오는 중...');
	let name = $state('');
	let greeting = $state<string>('');
	let clock = $state<string>('');
	let streaming = $state(false);

	// $derived: 다른 상태에서 파생되는 값 (React의 useMemo 역할, 의존성 배열 없음)
	let canGreet = $derived(name.trim().length > 0);

	// 페이지 로드 시 GET /api/hello 호출
	$effect(() => {
		api<{ message: string; time: string }>('/api/hello')
			.then((res) => (hello = res.message))
			.catch((e) => (hello = `백엔드 연결 실패: ${e.message}`));
	});

	// 경로 변수 데모: GET /api/hello/{name}
	async function greet() {
		const res = await api<{ message: string }>(`/api/hello/${encodeURIComponent(name)}`);
		greeting = res.message;
	}

	// ── SSE 데모: Go 백엔드가 1초마다 push하는 서버 시간 구독 ──
	let source: EventSource | null = null;

	function toggleStream() {
		if (streaming) {
			source?.close();
			source = null;
			streaming = false;
			clock = '';
		} else {
			source = new EventSource('/api/events');
			source.onmessage = (e) => (clock = e.data);
			source.onerror = () => toggleStream();
			streaming = true;
		}
	}

	// 컴포넌트 unmount 시 연결 정리
	$effect(() => {
		return () => source?.close();
	});
</script>

<svelte:head>
	<title>Go + SvelteKit Template</title>
</svelte:head>

<div class="space-y-6">
	<section class="rounded-lg border border-gray-200 bg-white p-6">
		<h1 class="text-2xl font-bold text-gray-900">🚀 Go + SvelteKit Template</h1>
		<p class="mt-2 text-gray-600">
			Go(표준 라이브러리) 백엔드 + SvelteKit(Svelte 5) 프론트 + Caddy 프록시.
			아래 데모들이 동작하면 풀스택이 연결된 것.
		</p>
	</section>

	<!-- 데모 1: 기본 API 호출 -->
	<section class="rounded-lg border border-gray-200 bg-white p-6">
		<h2 class="font-semibold text-gray-900">1. GET /api/hello</h2>
		<p class="mt-2 rounded bg-gray-100 px-3 py-2 font-mono text-sm text-gray-800">{hello}</p>
	</section>

	<!-- 데모 2: 경로 변수 -->
	<section class="rounded-lg border border-gray-200 bg-white p-6">
		<h2 class="font-semibold text-gray-900">2. GET /api/hello/&#123;name&#125; — 경로 변수</h2>
		<div class="mt-3 flex gap-2">
			<input
				bind:value={name}
				placeholder="이름 입력"
				class="flex-1 rounded border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none"
			/>
			<button
				onclick={greet}
				disabled={!canGreet}
				class="rounded bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-40"
			>
				인사받기
			</button>
		</div>
		{#if greeting}
			<p class="mt-3 rounded bg-gray-100 px-3 py-2 font-mono text-sm text-gray-800">{greeting}</p>
		{/if}
	</section>

	<!-- 데모 3: SSE 실시간 스트리밍 -->
	<section class="rounded-lg border border-gray-200 bg-white p-6">
		<h2 class="font-semibold text-gray-900">3. GET /api/events — SSE 실시간 스트리밍</h2>
		<p class="mt-1 text-sm text-gray-500">Go 서버가 1초마다 push하는 서버 시간</p>
		<div class="mt-3 flex items-center gap-4">
			<button
				onclick={toggleStream}
				class="rounded px-4 py-2 text-sm font-medium text-white {streaming
					? 'bg-red-600 hover:bg-red-700'
					: 'bg-green-600 hover:bg-green-700'}"
			>
				{streaming ? '중지' : '스트리밍 시작'}
			</button>
			{#if streaming}
				<span class="font-mono text-2xl tabular-nums text-gray-900">{clock || '…'}</span>
			{/if}
		</div>
	</section>
</div>
