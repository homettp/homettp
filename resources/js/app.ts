import { createApp, h } from 'vue';
import {
    createInertiaApp, Head, Link
} from '@inertiajs/vue3';
import type { DefineComponent } from 'vue';
import AppTitle from './base/AppTitle.vue';

import '../css/app.css';

window.isDark = () => localStorage.theme === 'dark'
    || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches);

window.currentTheme = () => (!('theme' in localStorage)
    ? 'system'
    : localStorage.theme);

window.setTheme = (theme: string) => {
    document.documentElement.classList.toggle('dark', window.isDark());

    document.dispatchEvent(new CustomEvent('set-theme', {
        detail: { theme }
    }));
};

createInertiaApp({
    resolve: (name: string) => {
        const pages = import.meta.glob<{ default: DefineComponent }>('./pages/**/*.vue', {
            eager: true
        });

        return pages[`./pages/${name}.vue`];
    },
    setup({
        el, App, props, plugin
    }) {
        createApp({ render: () => h(App, props) })
            .use(plugin)
            .component('InertiaHead', Head)
            .component('InertiaLink', Link)
            .component('AppTitle', AppTitle)
            .mount(el);
    }
});
