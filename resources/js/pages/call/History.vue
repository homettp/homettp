<template>
    <div class="call__history">
        <ol class="breadcrumb">
            <li class="breadcrumb-item">
                <inertia-link href="/">
                    Home
                </inertia-link>
            </li>
            <li class="breadcrumb-item">
                <inertia-link href="/">
                    Commands
                </inertia-link>
            </li>
            <li class="breadcrumb-item">
                <inertia-link :href="`/command/edit?id=${command.id}`">
                    {{ command.name }}
                </inertia-link>
            </li>
            <li class="breadcrumb-item active">
                {{ $metaInfo.title }}
            </li>
        </ol>
        <div v-if="!calls" class="call__history--card card card-body">
            No calls.
        </div>
        <div v-for="call in calls"
             :key="call.id"
             class="call__history--card card">
            <div class="card-header">
                <svg v-if="call.status === 'completed'"
                     class="bi text-success"
                     width="1em"
                     height="1em"
                     fill="currentColor">
                    <use :xlink:href="icon('check-circle')" />
                </svg>
                <svg v-else-if="call.status === 'failed'"
                     class="bi text-danger"
                     width="1em"
                     height="1em"
                     fill="currentColor">
                    <use :xlink:href="icon('x-circle')" />
                </svg>
                <div v-else class="spinner-border spinner-border-sm text-primary"></div>
                <span class="mr-auto">
                    #{{ call.id }}
                </span>
                <svg class="bi"
                     width="1em"
                     height="1em"
                     fill="currentColor">
                    <use :xlink:href="icon('clock-history')" />
                </svg>
                <span>
                    {{ call.created_at | date }}
                </span>
                <inertia-link class="btn btn-link ml-3"
                              :href="`/call/delete?id=${call.id}`"
                              method="delete"
                              as="button">
                    <svg class="bi"
                         width="1em"
                         height="1em"
                         fill="currentColor">
                        <use :xlink:href="icon('trash')" />
                    </svg>
                </inertia-link>
            </div>
            <div class="card-body pt-0">
                <pre v-if="call.output"><code>{{ call.output }}</code></pre>
            </div>
        </div>
    </div>
</template>

<script>
import { Inertia } from '@inertiajs/inertia';
import Layout from '../../common/Layout.vue';

export default {
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

    layout: Layout,

    data() {
        return {
            reloadInterval: undefined,
            reloadTimer: 2500
        };
    },

    metaInfo() {
        return {
            title: 'Call History'
        };
    },

    mounted() {
        this.reloadInterval = setInterval(() => Inertia.reload(), this.reloadTimer);
    },

    beforeDestroy() {
        clearInterval(this.reloadInterval);
    }
};
</script>
