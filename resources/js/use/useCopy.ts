import { ref } from 'vue';

export default (timeout = 2000) => {
    const copied = ref(false);
    let copiedTimeout: ReturnType<typeof setTimeout>;

    const copy = async (text: string) => {
        await navigator.clipboard.writeText(text);
        copied.value = true;
        clearTimeout(copiedTimeout);
        copiedTimeout = setTimeout(() => {
            copied.value = false;
        }, timeout);

        document.dispatchEvent(new CustomEvent('flash', {
            detail: { flash: 'Copied successfully.' }
        }));
    };

    return {
        copied,
        copy
    };
};
