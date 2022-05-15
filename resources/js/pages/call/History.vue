<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <div class="grid grid-cols-1 gap-5">
            <div v-if="!calls"
                 class="bg-white p-8">
                No calls.
            </div>
            <div v-for="call in calls"
                 :key="call.id"
                 class="bg-white p-8">
                <card-title>
                    <check-circle-icon v-if="call.status === 'completed'"
                                       class="h-6 w-6 text-green-400 sm:mr-2" />
                    <exclamation-circle-icon v-else-if="call.status === 'failed'"
                                             class="h-6 w-6 text-red-500 sm:mr-2" />
                    <!-- eslint-disable max-len -->
                    <div v-else class="flex h-3 w-3 relative sm:mr-2">
                        <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-cyan-400 opacity-75"></span>
                        <span class="relative inline-flex rounded-full h-3 w-3 bg-cyan-500"></span>
                    </div>
                    <!-- eslint-enable max-len -->
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
                                  as="button">
                        <trash-icon class="h-6 w-6" />
                    </inertia-link>
                </card-title>
                <!-- eslint-disable max-len -->
                <pre v-if="call.output" class="bg-gray-900 text-slate-300 text-xs p-6 whitespace-pre-wrap">{{ call.output }}</pre>
                <!-- eslint-enable max-len -->
            </div>
        </div>
    </div>
</template>

<script>
import {
    ClockIcon,
    CheckCircleIcon,
    ExclamationCircleIcon,
    TrashIcon
} from '@heroicons/vue/outline';

import { ref, onMounted, onUnmounted } from 'vue';
import { Inertia } from '@inertiajs/inertia';
import Breadcrumb from '../../common/Breadcrumb.vue';
import CardTitle from '../../common/CardTitle.vue';
import Layout from '../../common/Layout.vue';

export default {
    components: {
        ClockIcon,
        CheckCircleIcon,
        ExclamationCircleIcon,
        TrashIcon,
        Breadcrumb,
        CardTitle
    },

    layout: Layout,

    props: {
        command: {
            type: Object,
            required: true
        },

        calls: {
            type: Array,
            default: () => []
        }
    },

    setup(props) {
        const subtitle = ref('Call History');
        const reloadInterval = ref(undefined);
        const reloadTimer = 2500;

        const links = ref([
            { name: 'Commands', href: '/' },
            { name: props.command.name, href: `/command/edit?id=${props.command.id}` },
            { name: subtitle }
        ]);

        onMounted(() => {
            reloadInterval.value = setInterval(() => Inertia.reload(), reloadTimer);
        });

        onUnmounted(() => {
            clearInterval(reloadInterval.value);
        });

        return {
            subtitle,
            reloadInterval,
            reloadTimer,
            links
        };
    }
};
</script>
