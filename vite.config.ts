import { defineConfig } from 'vite';
import laravel from 'laravel-vite-plugin';
import tailwindcss from '@tailwindcss/vite';
import vue from '@vitejs/plugin-vue';

export default defineConfig(async () => {
    const { default: inertia } = await import('@inertiajs/vite');

    return {
        plugins: [
            laravel({
                input: 'resources/js/app.ts',
                publicDirectory: 'static'
            }),
            tailwindcss(),
            vue({
                template: {
                    transformAssetUrls: {
                        base: null,
                        includeAbsolute: false
                    }
                }
            }),
            inertia({
                ssr: false
            })
        ]
    };
});
