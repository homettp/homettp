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
import { router } from '@inertiajs/vue3';
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

const onFlash = e => {
    add(e.detail.flash);
};

document.addEventListener('flash', onFlash);

onUnmounted(() => {
    document.removeEventListener('flash', onFlash);
});

onUnmounted(
    router.on('success', e => {
        add(e.detail.page.props.flash as string);
    })
);
</script>
