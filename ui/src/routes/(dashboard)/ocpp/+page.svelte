<script lang="ts">

    let keys = [
        {key: "Key"},
        {key: "Key"}
    ]

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

                <label class="label">
                    <span class="label-text">CA</span>
                    <input class="file-input file-input-bordered w-full max-w-xs" type="file"/>
                </label>
                <label class="label">
                    <span class="label-text">Private key</span>
                    <input class="file-input file-input-bordered w-full max-w-xs" type="file"/>
                </label>
                <label class="label">
                    <span class="label-text">Certificate</span>
                    <input class="file-input file-input-bordered w-full max-w-xs" type="file"/>
                </label>

            </div>
        </div>
        <div class="divider lg:divider-horizontal"></div>
        <div class="card w-96 bg-base-100 text-primary-content">
            <div class="card-body">
                <h2 class="card-title">OCPP information</h2>
                <div class="form-control">
                    <label class="label">
                        <span class="label-text">Connection URL</span>
                    </label>
                    <label class="input-group">
                        <input class="input input-bordered" placeholder="ws://" type="text"/>
                    </label>
                </div>

                <div class="form-control">
                    <label class="label">
                        <span class="label-text">Connection URL</span>
                    </label>
                    <label class="input-group">
                        <input class="input input-bordered" placeholder="ws://" type="text"/>
                    </label>
                </div>

                <div class="form-control">
                    <label class="label">
                        <span class="label-text">Connection URL</span>
                    </label>
                    <label class="input-group">
                        <input class="input input-bordered" placeholder="ws://" type="text"/>
                    </label>
                </div>

                <div class="form-control">
                    <label class="label">
                        <span class="label-text">Connection URL</span>
                    </label>
                    <label class="input-group">
                        <input class="input input-bordered" placeholder="ws://" type="text"/>
                    </label>
                </div>

                <div class="form-control">
                    <label class="label">
                        <span class="label-text">Connection URL</span>
                    </label>
                    <label class="input-group">
                        <input class="input input-bordered" placeholder="ws://" type="text"/>
                    </label>
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
                    <td>{key.value}</td>
            {/each}
            </tbody>
        </table>
        <button class="btn btn-success">Apply</button>
    </div>
</div>
