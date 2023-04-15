<script lang="ts">
    import {onMount} from 'svelte';
    import {themeChange} from 'theme-change';

    let clazz: string = '';
    export {clazz as class};

    let theme: string | null = null;

    onMount(() => {
        themeChange(false);
        theme = document.documentElement.getAttribute('data-theme');
    });
</script>

<label class={`swap swap-rotate label cursor-pointer ${clazz}`}>
    <input
            checked={theme === 'light'}
            class="toggle"
            data-act-class="data-theme"
            data-toggle-theme="dark,light"
            type="checkbox"
    />

</label>