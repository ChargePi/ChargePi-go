<script lang="ts">
	export let opened = false;
	export let onDelete: Function = () => {};

	let loading = false;

	$: if (!opened) {
		setTimeout(() => (loading = false), 100);
	}

	function handleDelete() {
		loading = true;
		onDelete();
	}
</script>

<input class="modal-toggle" type="checkbox"/>
<div class="modal" class:modal-open={opened}>
    <div class="modal-box">
        <h3 class="font-bold text-lg">
            <slot name="title"/>
        </h3>
        <p class="py-4">
            <slot name="content"/>
        </p>
        <div class="modal-action">
            <button class="btn" on:click={() => (opened = false)}>Cancel</button>
            <button class="btn btn-primary" class:loading on:click={handleDelete}>Delete</button>
        </div>
    </div>
</div>
