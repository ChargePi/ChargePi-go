<script lang="ts">
    let Status = "Operative";
    let ActiveSessions = 1;
    let Errors = 0;
    let EVSEs = 2;
    let Active = "Connected";
    let user = "Admin";

</script>

<div class="flex flex-col w-full">
    <div class="flex pb-5">
        <h1 class="text-4xl">Welcome, {user}!</h1>
    </div>

    <div class="stats shadow">
        <div class="stat">
            <div class="stat-figure text-secondary">
                <svg
                        class="inline-block w-8 h-8 stroke-current"
                        fill="none"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                >
                    <path
                            d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                    />
                </svg
                >
            </div>
            <div class="stat-title">Number of EVSEs</div>
            <div class="stat-value">{EVSEs}</div>
        </div>

        <div class="stat">
            <div class="stat-figure text-secondary">
                <svg
                        class="inline-block w-8 h-8 stroke-current"
                        fill="none"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                >
                    <path
                            d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                    />
                </svg
                >
            </div>
            <div class="stat-title">Active sessions</div>
            <div class="stat-value">{ActiveSessions}</div>
            <div class="stat-desc">↗︎ 400 (22%)</div>
        </div>

        <div class="stat">
            <div class="stat-figure text-secondary">
                <svg
                        class="inline-block w-8 h-8 stroke-current"
                        fill="none"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                >
                    <path
                            d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                    />
                </svg
                >
            </div>
            <div class="stat-title">Errors in past 24h</div>
            <div class="stat-value">{Errors}</div>
        </div>

        <div class="stat">
            <div class="stat-title">Status</div>
            <div class="stat-value">{Status}</div>
        </div>
        <div class="stat">
            <div class="stat-title">Connection</div>
            <div class="stat-value">{Active}</div>
        </div>
    </div>
    <div class="divider"></div>
</div>