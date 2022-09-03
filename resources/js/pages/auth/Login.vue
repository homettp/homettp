<template>
    <app-title :title="subtitle" />
    <div class="container mx-auto">
        <header class="flex justify-center py-8">
            <inertia-link class="flex items-center text-2xl text-gray-800" href="/login">
                <home-icon class="h-8 w-8 mr-2" />
                <span>
                    {{ title }}
                </span>
            </inertia-link>
        </header>
        <div class="bg-white max-w-sm mx-auto p-8">
            <card-title>
                <arrow-right-on-rectangle-icon class="h-6 w-6 sm:mr-2" />
                <span>
                    {{ subtitle }}
                </span>
            </card-title>
            <form @submit.prevent="form.post('/login')">
                <div class="grid grid-cols-1 gap-6">
                    <label class="block">
                        <div>
                            Username / E-mail
                            <span class="text-cyan-500">
                                *
                            </span>
                        </div>
                        <input v-model="form.username_or_email"
                               class="form-input mt-3 block w-full"
                               :class="{'form-invalid': form.errors.username_or_email}"
                               type="text"
                               name="username_or_email"
                               placeholder="Username / E-mail"
                               required>
                        <div v-if="form.errors.username_or_email" class="mt-2 text-sm text-red-500">
                            {{ form.errors.username_or_email[0] }}
                        </div>
                    </label>
                    <label class="block">
                        <div>
                            Password
                            <span class="text-cyan-500">
                                *
                            </span>
                        </div>
                        <input id="password"
                               v-model="form.password"
                               class="form-input mt-3 block w-full"
                               :class="{'form-invalid': form.errors.password}"
                               type="password"
                               name="password"
                               placeholder="Password"
                               required>
                        <div v-if="form.errors.password" class="mt-2 text-sm text-red-500">
                            {{ form.errors.password[0] }}
                        </div>
                    </label>
                    <label class="flex items-center">
                        <input v-model="form.remember"
                               class="form-checkbox"
                               type="checkbox">
                        <span class="ml-2">
                            Remember Me
                        </span>
                    </label>
                    <button class="btn block" type="submit">
                        {{ subtitle }}
                    </button>
                </div>
            </form>
        </div>
    </div>
</template>

<script>
import { ref } from 'vue';
import { ArrowRightOnRectangleIcon, HomeIcon } from '@heroicons/vue/24/outline';
import { useForm } from '@inertiajs/inertia-vue3';
import CardTitle from '../../common/CardTitle.vue';

export default {
    components: {
        ArrowRightOnRectangleIcon,
        HomeIcon,
        CardTitle
    },

    props: {
        title: {
            type: String,
            required: true
        }
    },

    setup() {
        const subtitle = ref('Login');

        const form = useForm({
            username_or_email: '',
            password: '',
            remember: true
        });

        return {
            subtitle,
            form
        };
    }
};
</script>
