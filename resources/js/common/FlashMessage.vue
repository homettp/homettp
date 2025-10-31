<template>
    <!-- eslint-disable max-len -->
    <div class="cursor-pointer w-full mt-5 rounded-xs shadow-lg bg-white dark:bg-slate-700 ring-1 ring-black/5 md:w-96"
         @click="remove(message)">
        <div class="flex items-center p-5 font-medium text-sm">
            <check-circle-icon class="h-5 w-5 text-green-500 dark:text-green-400 mr-2" />
            <span>
                {{ message.value }}
            </span>
            <button class="ml-auto text-gray-400 hover:text-gray-500">
                <x-mark-icon class="h-5 w-5" />
            </button>
        </div>
        <div class="h-2.5 bg-slate-50 dark:bg-slate-600 transition-all ease-linear duration-100"
             :style="{width: `${progress}%`}"></div>
    </div>
    <!-- eslint-enable max-len -->
</template>

<script setup>
import {
    CheckCircleIcon,
    XMarkIcon
} from '@heroicons/vue/24/outline';

import {
    ref,
    onMounted,
    onUnmounted,
    defineProps
} from 'vue';

const {
    message,
    remove
} = defineProps({
    message: {
        type: Object,
        required: true
    },

    remove: {
        type: Function,
        required: true
    }
});

const messageTickInterval = ref();
const messageTickTimer = ref(100);
const messageTimer = ref(2500);
const progress = ref(100);

onMounted(() => {
    messageTickInterval.value = setInterval(() => {
        if (Date.now() - messageTimer.value > message.createdAt) {
            progress.value = 0;
            remove(message);
        } else {
            progress.value = (1 - (
                (Date.now() - message.createdAt) / messageTimer.value
            )) * 100;
        }
    }, messageTickTimer.value);
});

onUnmounted(() => {
    clearInterval(messageTickInterval.value);
});
</script>
