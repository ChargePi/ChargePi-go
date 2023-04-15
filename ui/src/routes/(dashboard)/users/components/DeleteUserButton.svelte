<script lang="ts">
	import DeleteModal from '$lib/modals/DeleteModal.svelte';
	import { faRemove } from '@fortawesome/free-solid-svg-icons';
	import Fa from 'svelte-fa/src/fa.svelte';

	export let onDelete: () => Promise<void> = async () => {};

	let opened = false;

	async function handleDelete() {
		await onDelete();
		opened = false;
	}
</script>

<DeleteModal bind:opened onDelete={handleDelete}>
	<span slot="title">Delete user</span>
	<span slot="content">Are you sure you want to delete this user?</span>
</DeleteModal>

<button class="btn btn-ghost btn-xs" on:click={() => (opened = true)}>
	<Fa icon={faRemove} />
</button>
