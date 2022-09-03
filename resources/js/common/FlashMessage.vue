<template>
    <!-- eslint-disable max-len -->
    <div class="cursor-pointer w-full mt-5 rounded-sm shadow-lg bg-white ring-1 ring-black ring-opacity-5 md:w-96"
         @click="remove(message)">
        <div class="flex items-center p-5 font-medium text-sm">
            <check-circle-icon class="h-5 w-5 text-green-400 mr-2" />
            <span>
                {{ message.value }}
            </span>
            <button class="ml-auto text-gray-400 hover:text-gray-500">
                <x-mark-icon class="h-5 w-5" />
            </button>
        </div>
        <div class="h-2.5 bg-slate-50 transition-all ease-linear duration-100"
             :style="{width: `${progress}%`}"></div>
    </div>
    <!-- eslint-enable max-len -->
</template>

<script>
import {
    CheckCircleIcon,
    XMarkIcon
} from '@heroicons/vue/24/outline';

import { ref, onMounted, onUnmounted } from 'vue';

export default {
    components: {
        CheckCircleIcon,
        XMarkIcon
    },

    props: {
        message: {
            type: Object,
            required: true
        },

        remove: {
            type: Function,
            required: true
        }
    },

    setup(props) {
        const messageTickInterval = ref();
        const messageTickTimer = ref(100);
        const messageTimer = ref(2500);
        const progress = ref(100);

        onMounted(() => {
            messageTickInterval.value = setInterval(() => {
                if (Date.now() - messageTimer.value > props.message.createdAt) {
                    progress.value = 0;
                    props.remove(props.message);
                } else {
                    progress.value = (1 - (
                        (Date.now() - props.message.createdAt) / messageTimer.value
                    )) * 100;
                }
            }, messageTickTimer.value);
        });

        onUnmounted(() => {
            clearInterval(messageTickInterval.value);
        });

        return {
            messageTickInterval,
            messageTickTimer,
            messageTimer,
            progress
        };
    }
};
</script>
