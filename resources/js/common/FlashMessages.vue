<template>
    <div class="fixed inset-x-5 top-20 z-10 md:left-auto">
        <flash-message v-for="message in messages"
                       :key="message.id"
                       :message="message"
                       :remove="remove" />
    </div>
</template>

<script>
import { ref, onUnmounted } from 'vue';
import { Inertia } from '@inertiajs/inertia';
import FlashMessage from './FlashMessage.vue';

export default {
    components: {
        FlashMessage
    },

    setup() {
        const messages = ref([]);

        const remove = message => {
            const index = messages.value.findIndex(current => current.id === message.id);

            if (index !== -1) {
                messages.value.splice(index, 1);
            }
        };

        const add = flash => {
            if (flash) {
                const createdAt = Date.now();
                const id = flash + Math.random() * createdAt;

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
            Inertia.on('success', e => {
                add(e.detail.page.props.flash);
            })
        );

        return {
            messages,
            remove
        };
    }
};
</script>
