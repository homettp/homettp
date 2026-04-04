<template>
    <div class="fixed inset-x-5 top-20 z-10 md:left-auto">
        <flash-message v-for="message in messages"
                       :key="message.id"
                       :message="message"
                       :remove="remove" />
    </div>
</template>

<script setup lang="ts">
import { ref, onUnmounted } from 'vue';
import type { FlashMessageItem } from '../types';
import FlashMessage from './FlashMessage.vue';

const messages = ref<FlashMessageItem[]>([]);

const remove = (message: FlashMessageItem) => {
    const index = messages.value.findIndex(current => current.id === message.id);

    if (index !== -1) {
        messages.value.splice(index, 1);
    }
};

const add = (flash: string) => {
    if (flash) {
        const createdAt = Date.now();
        const id = flash.length + Math.random() * createdAt;

        messages.value.unshift({
            id,
            value: flash,
            createdAt
        });
    }
};

const onInertiaFlash = (e: Event) => {
    const { flash } = (e as CustomEvent).detail;

    if (flash?.success) {
        add(flash.success);
    }
};

const onFlash = (e: Event) => {
    add((e as CustomEvent).detail.flash);
};

document.addEventListener('inertia:flash', onInertiaFlash);
document.addEventListener('flash', onFlash);

onUnmounted(() => {
    document.removeEventListener('inertia:flash', onInertiaFlash);
    document.removeEventListener('flash', onFlash);
});
</script>
