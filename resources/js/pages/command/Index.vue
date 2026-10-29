<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <!-- eslint-disable max-len -->
        <table-view :columns="columns" :items="filteredCommands">
            <template #title>
                <command-line-icon class="h-6 w-6" />
                <input v-model.trim="keyword"
                       class="form-input flex-1 sm:ml-2"
                       placeholder="Search">
            </template>
            <template #content="contentProps">
                <td class="border-t border-gray-100 dark:border-slate-600 px-6 py-3">
                    <inertia-link class="link flex items-center"
                                  :href="commandEditPath(contentProps.item.id)">
                        <component :is="iconName(contentProps.item)" class="h-6 w-6 mr-3" />
                        {{ contentProps.item.name }}
                    </inertia-link>
                </td>
                <td class="border-t border-gray-100 dark:border-slate-600 px-6 py-3">
                    {{ date(contentProps.item.created_at) }}
                </td>
                <td class="border-t border-gray-100 dark:border-slate-600 px-6 py-3">
                    <div class="flex justify-end gap-3">
                        <inertia-link class="link"
                                      title="Call History"
                                      :href="callHistoryPath(contentProps.item.id)">
                            <command-line-icon class="h-6 w-6" />
                        </inertia-link>
                        <inertia-link class="link"
                                      title="Edit Command"
                                      :href="commandEditPath(contentProps.item.id)">
                            <pencil-square-icon class="h-6 w-6" />
                        </inertia-link>
                    </div>
                </td>
            </template>
        </table-view>
        <!-- eslint-enable max-len -->
    </div>
</template>

<script setup lang="ts">
import {
    BoltIcon,
    CpuChipIcon,
    CommandLineIcon,
    EllipsisHorizontalCircleIcon,
    KeyIcon,
    LightBulbIcon,
    PencilSquareIcon
} from '@heroicons/vue/24/outline';

import { computed, ref } from 'vue';
import useDate from '../../use/useDate';
import usePaths from '../../use/usePaths';
import Breadcrumb from '../../base/Breadcrumb.vue';
import Layout from '../../base/Layout.vue';
import TableView from '../../base/TableView.vue';

type Command = {
    id: number
    name: string
    image: string
    created_at: string
};

const {
    commands = []
} = defineProps<{
    commands?: Command[]
}>();

defineOptions({
    layout: Layout
});

const { date } = useDate();
const { callHistoryPath, commandEditPath } = usePaths();

const keyword = ref('');
const subtitle = ref('Commands');

const links = ref([
    { name: subtitle }
]);

const columns = ref<string[]>([
    'Name',
    'Created'
]);

const filteredCommands = computed(() => {
    const value = keyword.value.toLowerCase();

    return (commands || []).filter(command => command.name.toLowerCase().includes(value));
});

const iconName = (command: Command) => {
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

    return CommandLineIcon;
};
</script>
