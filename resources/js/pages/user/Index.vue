<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <div class="grid grid-cols-1 gap-5 xl:grid-cols-4">
            <div v-for="user in users"
                 :key="user.id"
                 class="bg-white dark:bg-slate-700 p-8">
                <card-title>
                    <img class="rounded-full h-12 w-12 sm:mr-2"
                         :src="gravatars[user.username]">
                    <span class="flex-1 sm:mr-auto">
                        {{ user.username }}
                    </span>
                    <inertia-link class="link sm:ml-2"
                                  title="Edit User"
                                  :href="`/user/edit?id=${user.id}`">
                        <pencil-square-icon class="h-6 w-6" />
                    </inertia-link>
                </card-title>
                <div class="flex items-center">
                    <clock-icon class="h-5 w-5 mr-2" />
                    <span>
                        {{ date(user.created_at) }}
                    </span>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import {
    ClockIcon,
    PencilSquareIcon
} from '@heroicons/vue/24/outline';

import {
    ref,
    defineProps,
    defineOptions
} from 'vue';

import Breadcrumb from '../../common/Breadcrumb.vue';
import CardTitle from '../../common/CardTitle.vue';
import Layout from '../../common/Layout.vue';

defineProps({
    users: {
        type: Array,
        required: true
    },

    gravatars: {
        type: Object,
        required: true
    }
});

defineOptions({
    layout: Layout
});

const subtitle = ref('Users');

const links = ref([
    { name: subtitle }
]);
</script>
