<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <div class="grid grid-cols-1 divide-y divide-gray-100 dark:divide-gray-600">
            <div v-if="!calls"
                 class="bg-white dark:bg-slate-700 px-8 py-6">
                No calls.
            </div>
            <!-- eslint-disable max-len -->
            <div v-for="call in calls"
                 :key="call.id"
                 class="bg-white dark:bg-slate-700 px-8 pt-6 cursor-pointer hover:bg-gray-50 dark:hover:bg-slate-800"
                 @click="toggle(call)">
                <card-title>
                    <check-circle-icon v-if="call.status === 'completed'"
                                       class="h-6 w-6 text-green-400 sm:mr-2" />
                    <exclamation-circle-icon v-else-if="call.status === 'failed'"
                                             class="h-6 w-6 text-red-500 sm:mr-2" />
                    <div v-else class="flex h-3 w-3 relative sm:mr-2">
                        <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-cyan-400 opacity-75"></span>
                        <span class="relative inline-flex rounded-full h-3 w-3 bg-cyan-500"></span>
                    </div>
                    <span class="sm:mr-auto">
                        #{{ call.id }}
                    </span>
                    <clock-icon class="h-6 w-6 sm:mr-2" />
                    <span>
                        {{ date(call.created_at) }}
                    </span>
                    <inertia-link class="link sm:ml-2"
                                  :href="`/call/delete?id=${call.id}`"
                                  method="delete"
                                  as="button"
                                  @click.stop>
                        <trash-icon class="h-6 w-6" />
                    </inertia-link>
                </card-title>
                <div v-if="isShow(call)" class="mb-6">
                    <pre class="bg-gray-900 text-slate-300 text-xs p-6 whitespace-pre-wrap">{{ call.output || 'No output.' }}</pre>
                </div>
            </div>
            <!-- eslint-enable max-len -->
        </div>
    </div>
</template>

<script setup>
import {
    ClockIcon,
    CheckCircleIcon,
    ExclamationCircleIcon,
    TrashIcon
} from '@heroicons/vue/24/outline';

import {
    ref,
    onMounted,
    onUnmounted,
    defineProps,
    defineOptions
} from 'vue';

import { router } from '@inertiajs/vue3';
import Breadcrumb from '../../common/Breadcrumb.vue';
import CardTitle from '../../common/CardTitle.vue';
import Layout from '../../common/Layout.vue';
import useToggle from '../../common/useToggle';

const { command } = defineProps({
    command: {
        type: Object,
        required: true
    },

    calls: {
        type: Array,
        default: () => []
    }
});

defineOptions({
    layout: Layout
});

const subtitle = ref('Call History');
const { isShow, toggle } = useToggle();
const reloadTimer = 2500;

let reloadInterval;

const links = ref([
    { name: 'Commands', href: '/' },
    { name: command.name, href: `/command/edit?id=${command.id}` },
    { name: subtitle }
]);

onMounted(() => {
    reloadInterval = setInterval(() => router.reload(), reloadTimer);
});

onUnmounted(() => {
    clearInterval(reloadInterval);
});
</script>
