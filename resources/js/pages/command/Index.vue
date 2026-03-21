<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <div class="grid grid-cols-1 gap-5 xl:grid-cols-4">
            <div v-if="!commands"
                 class="bg-white dark:bg-slate-700 p-8">
                No commands.
            </div>
            <div v-for="command in commands"
                 :key="command.id"
                 class="bg-white dark:bg-slate-700 p-8">
                <card-title>
                    <component :is="iconName(command)" class="h-6 w-6 sm:mr-2" />
                    <span class="flex-1 sm:mr-auto">
                        {{ command.name }}
                    </span>
                    <inertia-link class="link sm:ml-2"
                                  title="Call History"
                                  :href="callHistoryPath(command.id)">
                        <command-line-icon class="h-6 w-6" />
                    </inertia-link>
                    <inertia-link class="link sm:ml-2"
                                  title="Edit Command"
                                  :href="commandEditPath(command.id)">
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

<script setup lang="ts">
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
import useDate from '../../use/useDate';
import usePaths from '../../use/usePaths';
import Breadcrumb from '../../base/Breadcrumb.vue';
import CardTitle from '../../base/CardTitle.vue';
import Layout from '../../base/Layout.vue';

const {
    commands = []
} = defineProps<{
    commands?: object[]
}>();

defineOptions({
    layout: Layout
});

const { date } = useDate();
const { callHistoryPath, commandEditPath } = usePaths();

const subtitle = ref('Commands');

const links = ref([
    { name: subtitle }
]);

const iconName = command => {
    if (command.image === 'door') {
        return KeyIcon;
    }

    if (command.image === 'light') {
        return LightBulbIcon;
    }

    if (command.image === 'outlet') {
        return EllipsisHorizontalCircleIcon;
    }

    if (command.image === 'plug') {
        return BoltIcon;
    }

    if (command.image === 'sensor') {
        return CpuChipIcon;
    }

    return 'command';
};
</script>
