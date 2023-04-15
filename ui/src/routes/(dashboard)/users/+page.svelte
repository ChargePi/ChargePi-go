<script lang="ts">
	import CreateUserForm from './components/CreateUserForm.svelte';
	import DeleteUserButton from './components/DeleteUserButton.svelte';

	type User = {
		username: string;
		role: string;
	};

	async function handleCreate(user: User) {
		return new Promise<void>((res, rej) => {
			setTimeout(() => {
				users = [...users, user];
				res();
			}, 500);
		});
	}

	async function handleDelete(user: User) {
		return new Promise<void>((res, rej) => {
			setTimeout(() => {
				users = users.filter((u) => u.username !== user.username);
				res();
			}, 500);
		});
	}

	let users: User[] = [
		{
			username: 'User1',
			role: 'Monkey',
		},
		{
			username: 'User2',
			role: 'Monkey',
		},
		{
			username: 'User3',
			role: 'Monkey',
		},
		{
			username: 'User4',
			role: 'Monkey',
		},
		{
			username: 'User5',
			role: 'Monkey',
		},
	];
</script>

<div class="overflow-x-auto">
	<table class="table w-full">
		<thead>
			<tr>
				<th>Username</th>
				<th>Role</th>
				<th />
			</tr>
		</thead>
		<tbody>
			{#each users as user}
				<tr>
					<td>{user.username}</td>
					<td>{user.role}</td>
					<td class="flex justify-end">
						<DeleteUserButton onDelete={() => handleDelete(user)} />
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>

<div class="divider" />

<CreateUserForm onCreate={handleCreate} />
