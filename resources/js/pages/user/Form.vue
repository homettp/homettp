<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <div class="bg-white p-8">
            <card-title>
                <component :is="iconName" class="h-6 w-6 sm:mr-2" />
                <span>
                    {{ subtitle }}
                </span>
            </card-title>
            <form @submit.prevent="submit()">
                <div class="grid grid-cols-1 gap-6">
                    <label class="block">
                        <div class="lg:flex lg:items-center">
                            <div class="lg:w-64">
                                Username
                                <span class="text-cyan-500">
                                    *
                                </span>
                            </div>
                            <input v-model="form.username"
                                   class="form-input mt-3 block w-full lg:mt-0 lg:flex-1"
                                   :class="{'form-invalid': form.errors.username}"
                                   type="text"
                                   name="username"
                                   placeholder="Username"
                                   required>
                        </div>
                        <div v-if="form.errors.username" class="mt-2 text-sm text-red-500 lg:ml-64">
                            {{ form.errors.username[0] }}
                        </div>
                    </label>
                    <label class="block">
                        <div class="lg:flex lg:items-center">
                            <div class="lg:w-64">
                                E-mail
                                <span class="text-cyan-500">
                                    *
                                </span>
                            </div>
                            <input v-model="form.email"
                                   class="form-input mt-3 block w-full lg:mt-0 lg:flex-1"
                                   :class="{'form-invalid': form.errors.email}"
                                   type="email"
                                   name="email"
                                   placeholder="E-mail"
                                   required>
                        </div>
                        <div v-if="form.errors.email" class="mt-2 text-sm text-red-500 lg:ml-64">
                            {{ form.errors.email[0] }}
                        </div>
                    </label>
                    <label class="block">
                        <div class="lg:flex lg:items-center">
                            <div class="lg:w-64">
                                Password
                                <span v-if="isNew"
                                      class="text-cyan-500">
                                    *
                                </span>
                            </div>
                            <input v-model="form.password"
                                   class="form-input mt-3 block w-full lg:mt-0 lg:flex-1"
                                   :class="{'form-invalid': form.errors.password}"
                                   type="password"
                                   name="password"
                                   placeholder="Password"
                                   :required="isNew">
                        </div>
                        <div v-if="form.errors.password" class="mt-2 text-sm text-red-500 lg:ml-64">
                            {{ form.errors.password[0] }}
                        </div>
                    </label>
                    <label class="flex items-center lg:ml-64">
                        <input v-model="form.is_enabled"
                               class="form-checkbox"
                               type="checkbox">
                        <span class="ml-2">
                            Is Enabled?
                        </span>
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
                <trash-icon class="h-6 w-6 sm:mr-2" />
                <span>
                    Delete User
                </span>
            </card-title>
            <div class="mb-6">
                Are you sure you want to delete the user?
            </div>
            <inertia-link class="btn-red"
                          :href="`/user/delete?id=${user.id}`"
                          method="delete"
                          as="button">
                Delete User
            </inertia-link>
        </div>
    </div>
</template>

<script>
import {
    TrashIcon,
    UserIcon,
    UserAddIcon
} from '@heroicons/vue/outline';

import { computed, toRefs } from 'vue';
import { useForm } from '@inertiajs/inertia-vue3';
import Breadcrumb from '../../common/Breadcrumb.vue';
import CardTitle from '../../common/CardTitle.vue';
import Layout from '../../common/Layout.vue';

export default {
    components: {
        TrashIcon,
        UserIcon,
        UserAddIcon,
        Breadcrumb,
        CardTitle
    },

    layout: Layout,

    props: {
        user: {
            type: Object,
            required: true
        }
    },

    setup(props) {
        const { user } = toRefs(props);
        const isNew = computed(() => user.value.id === 0);

        const subtitle = computed(() => (isNew.value
            ? 'Create User'
            : 'Edit User'));

        const iconName = computed(() => (isNew.value
            ? 'user-add-icon'
            : 'user-icon'));

        const url = computed(() => (isNew.value
            ? '/user/create'
            : `/user/edit?id=${user.value.id}`));

        const links = computed(() => [
            { name: 'Users', href: '/user' },
            { name: subtitle.value }
        ]);

        const form = useForm({
            username: user.value.username,
            email: user.value.email,
            password: '',
            is_enabled: user.value.is_enabled
        });

        const submit = () => {
            form.post(url.value, {
                onSuccess: () => form.reset('password')
            });
        };

        return {
            isNew,
            subtitle,
            iconName,
            url,
            links,
            form,
            submit
        };
    }
};
</script>
