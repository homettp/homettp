<template>
    <inertia-head :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <div class="bg-white p-8">
            <card-title>
                <component :is="iconName" class="h-6 w-6 mr-2" />
                <span>
                    {{ subtitle }}
                </span>
            </card-title>
            <form @submit.prevent="form.post(url)">
                <div class="grid grid-cols-1 gap-6">
                    <label class="block">
                        <div class="lg:flex lg:items-center">
                            <div class="lg:w-64">
                                Name
                                <span class="text-cyan-500">
                                    *
                                </span>
                            </div>
                            <input v-model="form.name"
                                   class="form-input mt-3 block w-full lg:mt-0 lg:flex-1"
                                   :class="{'form-invalid': form.errors.name}"
                                   type="text"
                                   name="name"
                                   placeholder="Name"
                                   required>
                        </div>
                        <div v-if="form.errors.name" class="mt-2 text-sm text-red-500 lg:ml-64">
                            {{ form.errors.name[0] }}
                        </div>
                    </label>
                    <label class="block">
                        <div class="lg:flex lg:items-center">
                            <div class="lg:w-64">
                                Image
                                <span class="text-cyan-500">
                                    *
                                </span>
                            </div>
                            <select v-model="form.image"
                                    class="form-select mt-3 block w-full lg:mt-0 lg:flex-1"
                                    :class="{'form-invalid': form.errors.image}"
                                    name="image"
                                    required>
                                <option v-for="commandImage in commandImages"
                                        :key="commandImage.value"
                                        :value="commandImage.value">
                                    {{ commandImage.name }}
                                </option>
                            </select>
                        </div>
                        <div v-if="form.errors.image" class="mt-2 text-sm text-red-500 lg:ml-64">
                            {{ form.errors.image[0] }}
                        </div>
                    </label>
                    <label class="block">
                        <div class="lg:flex lg:items-center">
                            <div class="lg:w-64">
                                Value
                                <span class="text-cyan-500">
                                    *
                                </span>
                            </div>
                            <input v-model="form.value"
                                   class="form-input mt-3 block w-full lg:mt-0 lg:flex-1"
                                   :class="{'form-invalid': form.errors.value}"
                                   type="text"
                                   name="name"
                                   placeholder="Value"
                                   required>
                        </div>
                        <div v-if="form.errors.value" class="mt-2 text-sm text-red-500 lg:ml-64">
                            {{ form.errors.value[0] }}
                        </div>
                        <div class="text-sm mt-3 lg:ml-64">
                            Payload Variable:
                            <code class="text-cyan-500">{{ commandPayload }}</code>
                        </div>
                    </label>
                    <div class="lg:ml-64">
                        <button class="btn" type="submit">
                            {{ subtitle }}
                        </button>
                    </div>
                </div>
            </form>
        </div>
        <div v-if="!isNew" class="bg-white p-8 mt-5">
            <card-title>
                <terminal-icon class="h-6 w-6 mr-2" />
                <span class="mr-auto">
                    Call Command
                </span>
                <inertia-link class="link flex items-center ml-2"
                              :href="`/call/history?id=${command.id}`">
                    <span class="mr-2">
                        Call History
                    </span>
                    <chevron-right-icon class="h-6 w-6" />
                </inertia-link>
            </card-title>
            <div class="mb-6">
                Making a
                <code class="text-cyan-500">GET</code>
                or
                <code class="text-cyan-500">POST</code>
                request to this URL will call your command:
                <div class="flex mt-3">
                    <input ref="commandPathEl"
                           class="form-input rounded-r-none flex-1"
                           type="text"
                           :value="commandPath"
                           readonly>
                    <a class="btn rounded-l-none flex items-center"
                       href="#"
                       @click.prevent="copy()">
                        <document-duplicate-icon class="h-5 w-5" />
                    </a>
                </div>
            </div>
            <inertia-link class="btn btn-primary"
                          :href="`/command/refresh-token?id=${command.id}`"
                          method="put"
                          as="button">
                Refresh Token
            </inertia-link>
        </div>
        <div v-if="!isNew" class="bg-white p-8 mt-5">
            <card-title>
                <trash-icon class="h-6 w-6 mr-2" />
                <span>
                    Delete Command
                </span>
            </card-title>
            <div class="mb-6">
                Are you sure you want to delete the command?
            </div>
            <inertia-link class="btn-red"
                          :href="`/command/delete?id=${command.id}`"
                          method="delete"
                          as="button">
                Delete Command
            </inertia-link>
        </div>
    </div>
</template>

<script>
import {
    ChevronRightIcon,
    DocumentDuplicateIcon,
    TerminalIcon,
    TrashIcon,
    ViewGridIcon,
    ViewGridAddIcon
} from '@heroicons/vue/outline';

import { computed, ref, toRefs } from 'vue';
import { useForm } from '@inertiajs/inertia-vue3';
import Breadcrumb from '../../common/Breadcrumb.vue';
import CardTitle from '../../common/CardTitle.vue';
import Layout from '../../common/Layout.vue';

export default {
    components: {
        ChevronRightIcon,
        DocumentDuplicateIcon,
        TerminalIcon,
        TrashIcon,
        ViewGridIcon,
        ViewGridAddIcon,
        Breadcrumb,
        CardTitle
    },

    layout: Layout,

    props: {
        command: {
            type: Object,
            required: true
        },

        commandImages: {
            type: Array,
            required: true
        },

        commandPayload: {
            type: String,
            required: true
        },

        commandPath: {
            type: String,
            default: ''
        }
    },

    setup(props) {
        const { command } = toRefs(props);
        const isNew = computed(() => command.value.id === 0);
        const commandPathEl = ref();

        const subtitle = computed(() => (isNew.value
            ? 'Create Command'
            : 'Edit Command'));

        const iconName = computed(() => (isNew.value
            ? 'view-grid-add-icon'
            : 'view-grid-icon'));

        const url = computed(() => (isNew.value
            ? '/command/create'
            : `/command/edit?id=${props.command.id}`));

        const links = computed(() => [
            { name: 'Commands', href: '/' },
            { name: subtitle.value }
        ]);

        const form = useForm({
            name: props.command.name,
            image: props.command.image,
            value: props.command.value
        });

        const copy = () => {
            commandPathEl.value.select();
            commandPathEl.value.setSelectionRange(0, 99999);

            document.execCommand('copy');

            document.dispatchEvent(new CustomEvent('flash', {
                detail: {
                    flash: 'Copied successfully.'
                }
            }));
        };

        return {
            isNew,
            commandPathEl,
            subtitle,
            links,
            form,
            iconName,
            url,
            copy
        };
    }
};
</script>
