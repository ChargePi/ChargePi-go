<script lang="ts">
	import { page } from '$app/stores';

	export let href: string;
</script>

<a {href} class:selected={$page.url.pathname === href}>
	<slot />
</a>

<style>
	.selected {
		@apply text-white;
		@apply bg-primary;
	}
</style>
