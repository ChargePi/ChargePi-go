<script lang="ts">
    import {goto} from '$app/navigation';
    import ThemeToggle from '$lib/theme/ThemeToggle.svelte';

    let loading = false;

    function handleLogin() {
        loading = true;
        // todo login logic

        setTimeout(() => {
            goto('/');
        }, 500);
    }
</script>

<div class="flex justify-end">
    <ThemeToggle class="p-4"/>
</div>

<div class="grid place-content-center h-3/4">
    <h1 class="text-5xl font-bold mb-12 text-center">ChargePi</h1>

    <form
            class="card w-96 bg-base-300 shadow-xl flex flex-col space-y-4 items-center p-10"
            on:submit|preventDefault={handleLogin}
    >
        <input
                class="input input-bordered w-full max-w-xs"
                name="username"
                placeholder="Username"
                required
                type="text"
        />

        <input
                class="input input-bordered w-full max-w-xs"
                name="password"
                placeholder="Password"
                required
                type="password"
        />

        <button class="btn btn-primary w-2/3" class:loading type="submit">Login</button>
    </form>
</div>
