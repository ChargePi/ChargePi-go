<script lang="ts">
	import '../app.css';

	export const ssr = false;

    const ocppConnected = false;
    const version = '0.1.0';
    let year = new Date().getFullYear();
</script>

<div class="flex flex-col h-screen">
	<div class="h-full bg-base-200">
		<slot />
	</div>

	<footer class="footer items-center p-4 bg-neutral text-neutral-content">
		<div class="items-center grid-flow-col">
            <p>Copyright ChargePi © {year} - All right reserved.</p>
		</div>
		<div class="grid-flow-col gap-2 md:place-self-center md:justify-self-end">
			<div
                    class="badge badge-success"
                    class:badge-error={!ocppConnected}
                    class:badge-success={ocppConnected}
			>
				{ocppConnected ? 'Connected' : 'Not Connected'}
			</div>
			<div class="badge badge-info">v{version}</div>
		</div>
	</footer>
</div>
