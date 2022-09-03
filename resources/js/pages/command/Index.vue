<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <div class="grid grid-cols-1 gap-5 xl:grid-cols-4">
            <div v-if="!commands"
                 class="bg-white p-8">
                No commands.
            </div>
            <div v-for="command in commands"
                 :key="command.id"
                 class="bg-white p-8">
                <card-title>
                    <component :is="iconName(command)" class="h-6 w-6 sm:mr-2" />
                    <span class="flex-1 sm:mr-auto">
                        {{ command.name }}
                    </span>
                    <inertia-link class="link sm:ml-2"
                                  title="Call History"
                                  :href="`/call/history?commandID=${command.id}`">
                        <command-line-icon class="h-6 w-6" />
                    </inertia-link>
                    <inertia-link class="link sm:ml-2"
                                  title="Edit Command"
                                  :href="`/command/edit?id=${command.id}`">
                        <pencil-square-icon class="h-6 w-6" />
                    </inertia-link>
                </card-title>
                <div class="flex items-center">
                    <clock-icon class="h-5 w-5 mr-2" />
                    <span>
                        {{ date(command.created_at) }}
                    </span>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import {
    BoltIcon,
    ClockIcon,
    CpuChipIcon,
    CommandLineIcon,
    EllipsisHorizontalCircleIcon,
    KeyIcon,
    LightBulbIcon,
    PencilSquareIcon
} from '@heroicons/vue/24/outline';

import { ref } from 'vue';
import Breadcrumb from '../../common/Breadcrumb.vue';
import CardTitle from '../../common/CardTitle.vue';
import Layout from '../../common/Layout.vue';

export default {
    components: {
        BoltIcon,
        ClockIcon,
        CpuChipIcon,
        CommandLineIcon,
        EllipsisHorizontalCircleIcon,
        KeyIcon,
        LightBulbIcon,
        PencilSquareIcon,
        Breadcrumb,
        CardTitle
    },

    layout: Layout,

    props: {
        commands: {
            type: Array,
            default: () => []
        }
    },

    setup() {
        const subtitle = ref('Commands');

        const links = ref([
            { name: subtitle }
        ]);

        const iconName = command => {
            if (command.image === 'door') {
                return 'key-icon';
            }

            if (command.image === 'light') {
                return 'light-bulb-icon';
            }

            if (command.image === 'outlet') {
                return 'ellipsis-horizontal-circle-icon';
            }

            if (command.image === 'plug') {
                return 'bolt-icon';
            }

            if (command.image === 'sensor') {
                return 'cpu-chip-icon';
            }

            return 'command';
        };

        return {
            subtitle,
            links,
            iconName
        };
    }
};
</script>
