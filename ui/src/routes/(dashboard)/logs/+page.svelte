<script lang="ts">
	const logs = [
		{
			timestamp: new Date(),
			message: 'Lorem ipsum, dolor sit amet consectetur adipisicing elit.',
		},
		{
			timestamp: new Date(),
			message: 'Lorem ipsum, dolor sit amet consectetur adipisicing elit.',
		},
		{
			timestamp: new Date(),
			message: 'Lorem ipsum, dolor sit amet consectetur adipisicing elit.',
		},
	];
</script>

<div class="overflow-x-auto">
	<button class="btn">Button</button>
	<table class="table w-full">
		<thead>
		<tr>
			<th>Timestamp</th>
			<th>Severity</th>
			<th>Message</th>
		</tr>
		</thead>
		<tbody>
		{#each logs as log}
			<tr>
					<td>{log.timestamp.toISOString()}</td>
				<td>{""}</td>
					<td>{log.message}</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
