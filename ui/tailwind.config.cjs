module.exports = {
    content: ['./src/routes/**/*.{svelte,js,ts}'],
    plugins: [require('daisyui')],
    daisyui: {
        themes: [
            {
                dark: {
                    ...require('daisyui/src/colors/themes')['[data-theme=dark]'],
                    primary: '#e92b2b',
                    "secondary": "#ef4444",

                    "accent": "#f59e0b",

                    "neutral": "#172027",

                    "base-100": "#374151",

                    "info": "#4b5563",

                    "success": "#0F7041",

                    "warning": "#F6A431",

                    "error": "#E43F42",
                },
            },
            {
                light: {
                    ...require('daisyui/src/colors/themes')['[data-theme=light]'],
                    primary: '#e92b2b',
                },
            },
        ],
    },
};
