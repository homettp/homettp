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
            </div>
            <div class="card-body">
                <pre class="mb-0"><code>{{ call.output }}</code></pre>
            </div>
        </div>
    </div>
</template>

<script>
import Layout from '../../common/Layout.vue';

export default {
    props: {
        command: {
            type: Object,
            required: true
        },

        calls: {
            type: Array,
            required: true
        }
    },

    layout: Layout,

    metaInfo() {
        return {
            title: 'Call History'
        };
    }
};
</script>
