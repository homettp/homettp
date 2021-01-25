<template>
    <div class="user__index">
        <ol class="breadcrumb">
            <li class="breadcrumb-item">
                <inertia-link href="/">
                    Home
                </inertia-link>
            </li>
            <li class="breadcrumb-item active">
                {{ $metaInfo.title }}
            </li>
        </ol>
        <div class="user__row row">
            <div v-for="user in users"
                 :key="user.id"
                 class="user__col col-xl-3">
                <div class="card">
                    <div class="card-header">
                        <img class="gravatar img-fluid rounded-circle"
                             :src="gravatars[user.username]">
                        <span class="mr-auto">
                            {{ user.username }}
                        </span>
                        <inertia-link class="ml-3" :href="`/user/edit?id=${user.id}`">
                            <svg class="bi"
                                 width="1em"
                                 height="1em"
                                 fill="currentColor">
                                <use :xlink:href="icon('pencil')" />
                            </svg>
                        </inertia-link>
                    </div>
                    <div class="card-body d-flex align-items-center">
                        <svg class="bi"
                             width="1em"
                             height="1em"
                             fill="currentColor">
                            <use :xlink:href="icon('clock-history')" />
                        </svg>
                        <span>
                            {{ user.created_at | date }}
                        </span>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import Layout from '../../common/Layout.vue';

export default {
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

    layout: Layout,

    metaInfo() {
        return {
            title: 'Users'
        };
    }
};
</script>
