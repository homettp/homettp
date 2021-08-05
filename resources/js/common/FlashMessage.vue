<template>
    <div ref="root" class="toast fade hide">
        <div class="toast-body d-flex align-items-center">
            <svg class="bi"
                 width="1em"
                 height="1em"
                 fill="currentColor">
                <use :xlink:href="icon('check-circle')" />
            </svg>
            <span>
                {{ message }}
            </span>
        </div>
    </div>
</template>

<script>
import { ref, onUnmounted } from 'vue';
import { Inertia } from '@inertiajs/inertia';

export default {
    setup() {
        const root = ref();
        const message = ref('');

        const show = flash => {
            if (flash) {
                message.value = flash;
                $(root.value).toast('show');
            }
        };

        const onFlash = e => {
            show(e.detail.flash);
        };

        document.addEventListener('flash', onFlash);

        onUnmounted(() => {
            document.removeEventListener('flash', onFlash);
        });

        onUnmounted(
            Inertia.on('success', e => {
                show(e.detail.page.props.flash);
            })
        );

        return {
            root,
            message
        };
    }
};
</script>
