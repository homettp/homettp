<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <div class="grid grid-cols-1 gap-5 xl:grid-cols-4">
            <div v-for="user in users"
                 :key="user.id"
                 class="bg-white p-8">
                <card-title>
                    <img class="rounded-full h-12 w-12 mr-2"
                         :src="gravatars[user.username]">
                    <span class="flex-1 mr-auto">
                        {{ user.username }}
                    </span>
                    <inertia-link class="link ml-2"
                                  title="Edit User"
                                  :href="`/user/edit?id=${user.id}`">
                        <pencil-alt-icon class="h-6 w-6" />
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

<script>
import {
    ClockIcon,
    PencilAltIcon
} from '@heroicons/vue/outline';

import { ref } from 'vue';
import Breadcrumb from '../../common/Breadcrumb.vue';
import CardTitle from '../../common/CardTitle.vue';
import Layout from '../../common/Layout.vue';

export default {
    components: {
        ClockIcon,
        PencilAltIcon,
        Breadcrumb,
        CardTitle
    },

    layout: Layout,

    props: {
        users: {
            type: Array,
            required: true
        },

        gravatars: {
            type: Object,
            required: true
        }
    },

    setup() {
        const subtitle = ref('Users');

        const links = ref([
            { name: subtitle }
        ]);

        return {
            subtitle,
            links
        };
    }
};
</script>
