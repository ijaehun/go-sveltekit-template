<script lang="ts">
	import '../app.css';
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';

	// Svelte 5: children은 $props()로 받아 {@render}로 렌더링
	let { children } = $props();
</script>

<div class="flex min-h-screen flex-col bg-gray-50">
	<Header />
	<main class="mx-auto w-full max-w-3xl flex-1 px-4 py-8">
		{@render children()}
	</main>
	<Footer />
</div>
