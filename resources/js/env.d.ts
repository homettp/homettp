/// <reference types="vite/client" />

declare global {
    interface Window {
        currentTheme: () => string
        isDark: () => boolean
        setTheme: (theme: string) => void
    }
}

declare module '*.vue';

export {};
