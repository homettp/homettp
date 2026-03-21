<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <div class="bg-white dark:bg-slate-700 p-8">
            <card-title>
                <component :is="iconName" class="h-6 w-6 sm:mr-2" />
                <span>
                    {{ subtitle }}
                </span>
            </card-title>
            <!-- eslint-disable max-len -->
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
                        <div v-if="form.errors.name" class="mt-2 text-sm text-red-500 dark:text-red-400 lg:ml-64">
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
                        <div v-if="form.errors.image" class="mt-2 text-sm text-red-500 dark:text-red-400 lg:ml-64">
                            {{ form.errors.image[0] }}
                        </div>
                    </label>
                    <!-- eslint-disable max-len -->
                    <label class="block">
                        <div class="lg:flex lg:items-center">
                            <div class="lg:w-64">
                                Value
                                <span class="text-cyan-500">
                                    *
                                </span>
                            </div>
                            <input v-model="form.value"
                                   class="form-input mt-3 block w-full font-mono text-sm lg:mt-0 lg:flex-1"
                                   :class="{'form-invalid': form.errors.value}"
                                   type="text"
                                   name="name"
                                   placeholder="myapp --toggle light1 --input %p"
                                   required>
                        </div>
                        <div v-if="form.errors.value" class="mt-2 text-sm text-red-500 dark:text-red-400 lg:ml-64">
                            {{ form.errors.value[0] }}
                        </div>
                        <div class="text-sm mt-3 lg:ml-64">
                            Payload Variable:
                            <code class="text-cyan-500">{{ commandPayload }}</code>
                        </div>
                    </label>
                    <!-- eslint-enable max-len -->
                    <div class="lg:ml-64">
                        <button class="btn" type="submit">
                            {{ subtitle }}
                        </button>
                    </div>
                </div>
            </form>
            <!-- eslint-enable max-len -->
        </div>
        <div v-if="!isNew" class="bg-white dark:bg-slate-700 p-8 mt-5">
            <card-title>
                <command-line-icon class="h-6 w-6 sm:mr-2" />
                <span class="mr-auto">
                    Call Command
                </span>
                <inertia-link class="link flex items-center sm:ml-2"
                              :href="callHistoryPath(command.id)">
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
                       @click.prevent="copy(commandPathEl.value)">
                        <document-duplicate-icon class="h-5 w-5" />
                    </a>
                </div>
            </div>
            <inertia-link class="btn btn-primary"
                          :href="commandRefreshTokenPath(command.id)"
                          method="put"
                          as="button">
                Refresh Token
            </inertia-link>
        </div>
        <div v-if="!isNew" class="bg-white dark:bg-slate-700 p-8 mt-5">
            <card-title>
                <trash-icon class="h-6 w-6 sm:mr-2" />
                <span>
                    Delete Command
                </span>
            </card-title>
            <div class="mb-6">
                Are you sure you want to delete the command?
            </div>
            <inertia-link class="btn-red"
                          :href="commandDeletePath(command.id)"
                          method="delete"
                          as="button">
                Delete Command
            </inertia-link>
        </div>
    </div>
</template>

<script setup lang="ts">
import {
    ChevronRightIcon,
    CommandLineIcon,
    DocumentDuplicateIcon,
    PlusIcon,
    TrashIcon
} from '@heroicons/vue/24/outline';

import { computed, ref } from 'vue';
import { useForm } from '@inertiajs/vue3';
import type { CommandImage } from '../../types';
import useCopy from '../../use/useCopy';
import usePaths from '../../use/usePaths';
import Breadcrumb from '../../base/Breadcrumb.vue';
import CardTitle from '../../base/CardTitle.vue';
import Layout from '../../base/Layout.vue';

const {
    command,
    commandImages,
    commandPayload,
    commandPath = ''
} = defineProps<{
    command: Record<string, any>
    commandImages: CommandImage[]
    commandPayload: string
    commandPath?: string
}>();

defineOptions({
    layout: Layout
});

const { copy } = useCopy();
const {
    callHistoryPath,
    commandCreatePath,
    commandDeletePath,
    commandEditPath,
    commandIndexPath,
    commandRefreshTokenPath
} = usePaths();

const isNew = computed(() => command.id === 0);
const commandPathEl = ref<HTMLInputElement>();

const subtitle = computed(() => (isNew.value
    ? 'Create Command'
    : 'Edit Command'));

const iconName = computed(() => (isNew.value
    ? PlusIcon
    : CommandLineIcon));

const url = computed(() => (isNew.value
    ? commandCreatePath()
    : commandEditPath(command.id)));

const links = computed(() => [
    { name: 'Commands', href: commandIndexPath() },
    { name: subtitle.value }
]);

const form = useForm({
    name: command.name,
    image: command.image,
    value: command.value
});
</script>
