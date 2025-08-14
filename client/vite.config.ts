import { svelte } from '@sveltejs/vite-plugin-svelte';
import tailwindcss from '@tailwindcss/vite';
import path from 'path';
import { defineConfig } from 'vite';

// https://vite.dev/config/
export default defineConfig({
    plugins: [tailwindcss(), svelte()],
    resolve: {
        alias: {
            $lib: path.resolve('./src/lib'),
            $views: path.resolve('./src/views'),
        },
    },
    build: {
        outDir: '../build/client',
    },
    base: '/app',
});
