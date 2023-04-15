<script lang="ts">

    let Status = "Operative";
    let ActiveSessions = 1;
    let EVSEs = 2;
    let Sessions = 3;

    let Evses = [
        {
            EvseID: 1,
            Status: "Operative",
            EVCC: "EVCC 1",
            PowerMeter: "PowerMeter 1",
            MaxPower: 10000
        },
        {
            EvseID: 2,
            Status: "Operative",
            EVCC: "EVCC 2",
            PowerMeter: "PowerMeter 2",
            MaxPower: 10000
        }
    ]
</script>

<div class="flex flex-col w-full">
    <div class="stats stats-vertical lg:stats-horizontal shadow">
        <div class="stat">
            <div class="stat-title">EVSEs</div>
            <div class="stat-value">{EVSEs}</div>
            <div class="stat-desc"></div>
        </div>

        <div class="stat">
            <div class="stat-title">Active sessions</div>
            <div class="stat-value">{ActiveSessions}</div>
        </div>

        <div class="stat">
            <div class="stat-title">Sessions in past 24h</div>
            <div class="stat-value">{Sessions}</div>
        </div>

    </div>
    <div class="divider"></div>
    <div class="overflow-x-auto">
        <table class="table w-full">
            <!-- head -->
            <thead>
            <tr>
                <th>EvseID</th>
                <th>Status</th>
                <th>EVCC</th>
                <th>Power Meter</th>
                <th>Max Power</th>
                <th></th>
            </tr>
            </thead>
            <tbody>
            {#each Evses as evse}
                <tr>
                    <td>{evse.EvseID}</td>
                    <td>{evse.Status}</td>
                    <td>{evse.EVCC}</td>
                    <td>{evse.PowerMeter}</td>
                    <td>{evse.MaxPower}</td>
                    <td>

                    </td>
            {/each}
            </tbody>
        </table>
    </div>
</div>

