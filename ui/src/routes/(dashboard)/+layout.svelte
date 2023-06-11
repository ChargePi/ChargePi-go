<script lang="ts">
    import Fa from 'svelte-fa/src/fa.svelte';
    import {
        faBars,
        faBolt,
        faCloud,
        faCreditCard,
        faFilePen,
        faHome,
        faMicrochip,
        faRightFromBracket,
        faUser,
    } from '@fortawesome/free-solid-svg-icons';
    import NavItem from '$lib/navigation/NavItem.svelte';
    import ThemeToggle from '$lib/theme/ThemeToggle.svelte';

    let expand = true;
</script>

<div class="flex flex-col h-full space-y-4 p-4">
    <div class="navbar rounded-box bg-base-300 shadow-sm">
        <div class="flex-1 px-5">
            <img alt="logo" class="" src="../../../static/logo.svg"/>
        </div>
        <div class="flex-none">
            <div class="dropdown dropdown-end">
                <label class="btn btn-ghost btn-circle avatar" tabindex="0">
                    <div class="w-10 rounded-full">
                        <img src=""/>
                    </div>
                </label>
                <ul class="menu menu-compact dropdown-content mt-3 p-2 shadow bg-base-100 rounded-box w-52"
                    tabindex="0">
                    <li>
                        <div class="flex-none">
                            <ThemeToggle/>
                        </div>
                    </li>
                    <li>
                        <button class="btn btn-square btn-ghost" on:click={() => (expand = !expand)}>
                            <Fa icon={faBars}/>
                        </button>
                    </li>
                    <li><a>Settings</a></li>


                    <li><a href="/login">
                        <Fa icon={faRightFromBracket}/>
                        Logout</a></li>
                </ul>
            </div>

        </div>
    </div>

    <div class="flex h-full space-x-4">
        <div class="card bg-base-300 shadow-sm grow p-4 min-w-0">
            <slot/>
        </div>
        {#if expand}
            <div class="card bg-base-300 shadow-sm">
                <ul class="menu w-64 p-4">
                    <li>
                        <NavItem href="/">
                            <Fa icon={faHome}/>
                            Dashboard
                        </NavItem>
                    </li>
                    <li>
                        <NavItem href="/evse">
                            <Fa icon={faBolt}/>
                            EVSE
                        </NavItem>
                    </li>
                    <li>
                        <NavItem href="/hardware">
                            <Fa icon={faMicrochip}/>
                            Hardware
                        </NavItem>
                    </li>
                    <li>
                        <NavItem href="/ocpp">
                            <Fa icon={faCloud}/>
                            OCPP
                        </NavItem>
                    </li>
                    <li>
                        <NavItem href="/tags">
                            <Fa icon={faCreditCard}/>
                            Tags
                        </NavItem>
                    </li>
                    <li>
                        <NavItem href="/users">
                            <Fa icon={faUser}/>
                            Users
                        </NavItem>
                    </li>
                    <li>
                        <NavItem href="/logs">
                            <Fa icon={faFilePen}/>
                            Logs
                        </NavItem>
                    </li>

                </ul>
            </div>
        {/if}
    </div>
</div>
