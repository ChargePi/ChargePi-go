<script lang="ts">
    type CreateUser = {
        username: string;
        password: string;
        role: string;
    };

    const roles = ['Monke', 'Admin'];

    export let onCreate: (user: CreateUser) => Promise<void> = async () => {
        // todo call api
    };

    let user: CreateUser = {username: '', password: '', role: 'Role'};
    let loading = false;

    async function handleCreate() {
        loading = true;
        await onCreate(user);
        loading = false;
    }
</script>

<form class="flex space-x-4 w-full justify-center" on:submit|preventDefault={handleCreate}>
    <input
            bind:value={user.username}
            class="input w-full max-w-xs"
            placeholder="Username"
            required
            type="text"
    />
    <input
            bind:value={user.password}
            class="input w-full max-w-xs"
            placeholder="Password"
            required
            type="password"
    />
    <select bind:value={user.role} class="select w-full max-w-xs" required>
        <option disabled selected>Role</option>
        {#each roles as role}
            <option>{role}</option>
        {/each}
    </select>
    <button class="btn" class:loading>Create</button>
</form>
