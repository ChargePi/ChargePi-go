<script lang="ts">

    let keys = [
        {key: "Key"},
        {key: "Key"}
    ]

    let ocppInfo = {
        vendor: "xBlaz3kx",
        model: "ChargePi",
        version: "v0.1.1",
    }

</script>

<div class="w-full">
    <div class="flex flex-col w-full lg:flex-row">
        <div class="card w-96 bg-base-100 text-primary-content">
            <div class="card-body">
                <h2 class="card-title b">Backend connection</h2>
                <div class="form-control">
                    <label class="label">
                        <span class="label-text">Connection URL</span>
                    </label>
                    <label class="input-group">
                        <input class="input input-bordered" placeholder="ws://" type="text"/>
                        <button class="btn btn">Apply</button>
                    </label>
                </div>

                <div class="collapse">
                    <input type="checkbox"/>
                    <div class="collapse-title font-medium">
                        TLS
                    </div>
                    <div class="collapse-content rounded-box bg-base-200 ">
                        <label class="label">
                            <span class="label-text">CA</span>
                            <input class="file-input file-input-xs w-full max-w-xs" type="file"/>
                        </label>
                        <label class="label">
                            <span class="label-text">Private key</span>
                            <input class="file-input file-input-xs w-full max-w-xs" type="file"/>
                        </label>
                        <label class="label">
                            <span class="label-text">Certificate</span>
                            <input class="file-input file-input-xs w-full max-w-xs" type="file"/>
                        </label>
                    </div>
                </div>

            </div>
        </div>
        <div class="divider lg:divider-horizontal"></div>
        <div class="card w-96 bg-base-100 text-primary-content">
            <div class="card-body">
                <h2 class="card-title">OCPP information</h2>
                {#each Object.keys(ocppInfo) as key}
                    <div class="form-control">
                        <label class="label">
                            <span class="label-text">{key}</span>
                        </label>
                        <label class="input-group">
                            <input class="input input-bordered" placeholder="{ocppInfo[key]}" type="text"/>
                        </label>
                    </div>
                {/each}
                <div class="flex flex-col py-5">
                    <button class="btn btn-primary">Apply</button>
                </div>
            </div>
        </div>
        <div class="divider lg:divider-horizontal"></div>
        <div class="card w-96 bg-base-100 text-primary-content">
            <div class="card-body">
                <h2 class="card-title">Operation details</h2>
                <p></p>
            </div>
        </div>
    </div>

    <div class="divider"></div>

    <div class="flex w-full overflow-x-auto">

        <table class="table w-full">
            <thead>
            <tr>
                <th>Key</th>
                <th>ReadOnly</th>
                <th>Value</th>
            </tr>
            </thead>
            <tbody>
            {#each keys as key}
                <tr>
                    <td>{key.key}</td>
                    <td>{key.readOnly}</td>
                    <td><input class="input input-bordered" placeholder="{key.value}" type="text"/></td>
                    <td>
                        <button class="btn btn-primary">Apply</button>
                    </td>
                </tr>
            {/each}
            </tbody>
        </table>
    </div>
</div>
