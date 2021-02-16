<template>
    <div class="command__index layout__index">
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
        <div class="layout__row row">
            <div v-if="!commands" class="layout__col col">
                <div class="card card-body">
                    No commands.
                </div>
            </div>
            <div v-for="command in commands"
                 :key="command.id"
                 class="layout__col col-xl-3">
                <div class="card">
                    <div class="card-header">
                        <span class="command__icon">
                            <svg class="bi"
                                 width="1em"
                                 height="1em"
                                 fill="currentColor">
                                <use :xlink:href="icon(iconName(command))" />
                            </svg>
                        </span>
                        <span class="mr-auto">
                            {{ command.name }}
                        </span>
                        <inertia-link v-tooltip
                                      class="ml-3"
                                      data-title="Call History"
                                      :href="`/call/history?id=${command.id}`">
                            <svg class="bi"
                                 width="1em"
                                 height="1em"
                                 fill="currentColor">
                                <use :xlink:href="icon('terminal')" />
                            </svg>
                        </inertia-link>
                        <inertia-link v-tooltip
                                      class="ml-3"
                                      data-title="Edit Command"
                                      :href="`/command/edit?id=${command.id}`">
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
                            {{ command.created_at | date }}
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
        commands: {
            type: Array,
            default: () => []
        }
    },

    layout: Layout,

    metaInfo() {
        return {
            title: 'Commands'
        };
    },

    methods: {
        iconName(command) {
            if (command.image === 'door') {
                return 'door-open';
            }

            if (command.image === 'light') {
                return 'lightbulb';
            }

            if (command.image === 'outlet') {
                return 'outlet';
            }

            if (command.image === 'plug') {
                return 'plug';
            }

            if (command.image === 'sensor') {
                return 'cpu';
            }

            return 'command';
        }
    }
};
</script>
