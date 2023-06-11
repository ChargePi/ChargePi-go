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

<div class="hero min-h-screen bg-base-200">
    <div class="hero-content flex-col lg:flex-row-reverse">
        <div class="text-center lg:text-left px-10">
            <h1 class="text-5xl font-bold">ChargePi </h1>
            <p class="py-6 space-x-4">An open source Charge Point.</p>
        </div>
        <div class="card flex-shrink-0 w-full max-w-sm shadow-2xl bg-base-100">
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
    </div>
</div>
