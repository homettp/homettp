<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <!-- eslint-disable max-len -->
        <table-view :columns="columns" :items="filteredUsers">
            <template #title>
                <user-group-icon class="h-6 w-6" />
                <input v-model.trim="keyword"
                       class="form-input flex-1 sm:ml-2"
                       placeholder="Search">
            </template>
            <template #content="contentProps">
                <td class="border-t border-gray-100 dark:border-slate-600 px-6 py-3">
                    <inertia-link class="link flex items-center"
                                  :href="userEditPath(contentProps.item.id)">
                        <img class="rounded-full h-8 w-8 mr-3"
                             :src="gravatars[contentProps.item.username]">
                        {{ contentProps.item.username }}
                    </inertia-link>
                </td>
                <td class="border-t border-gray-100 dark:border-slate-600 px-6 py-3">
                    {{ date(contentProps.item.created_at) }}
                </td>
                <td class="border-t border-gray-100 dark:border-slate-600 px-6 py-3">
                    <div class="flex justify-end gap-3">
                        <inertia-link class="link"
                                      title="Edit User"
                                      :href="userEditPath(contentProps.item.id)">
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
    PencilSquareIcon,
    UserGroupIcon
} from '@heroicons/vue/24/outline';

import { computed, ref } from 'vue';
import useDate from '../../use/useDate';
import usePaths from '../../use/usePaths';
import Breadcrumb from '../../base/Breadcrumb.vue';
import Layout from '../../base/Layout.vue';
import TableView from '../../base/TableView.vue';

type User = {
    id: number
    username: string
    created_at: string
};

const {
    users = [],
    gravatars
} = defineProps<{
    users?: User[]
    gravatars: Record<string, string>
}>();

defineOptions({
    layout: Layout
});

const { date } = useDate();
const { userEditPath } = usePaths();

const keyword = ref('');
const subtitle = ref('Users');

const links = ref([
    { name: subtitle }
]);

const columns = ref<string[]>([
    'Username',
    'Created'
]);

const filteredUsers = computed(() => {
    const value = keyword.value.toLowerCase();

    return (users || []).filter(user => user.username.toLowerCase().includes(value));
});
</script>
