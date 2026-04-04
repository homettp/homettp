import {
    createInertiaApp, Head, Link
} from '@inertiajs/vue3';
import AppTitle from './base/AppTitle.vue';

import '../css/app.css';

if (typeof window !== 'undefined') {
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
}

createInertiaApp({
    pages: './pages',
    withApp(app) {
        app
            .component('InertiaHead', Head)
            .component('InertiaLink', Link)
            .component('AppTitle', AppTitle);
    }
});
