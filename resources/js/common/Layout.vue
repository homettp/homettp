<template>
    <flash-messages />
    <div v-show="isSidebarOpen"
         class="bg-black bg-opacity-50 fixed inset-0 z-10 md:hidden"
         @click.prevent="isSidebarOpen = false"></div>
    <!-- eslint-disable max-len -->
    <div class="fixed inset-y-0 left-0 bg-gray-800 w-60 z-20 transform transition-transform md:translate-x-0"
         :class="{'translate-x-0': isSidebarOpen, '-translate-x-full': !isSidebarOpen}">
        <div class="flex h-20 bg-gray-900">
            <inertia-link class="flex items-center m-auto text-white text-xl"
                          href="/">
                <home-icon class="h-7 w-7 mr-2" />
                <span>
                    {{ $page.props.title }}
                </span>
            </inertia-link>
        </div>
        <div class="sidebar overflow-y-auto">
            <sidebar-title>
                Commands
            </sidebar-title>
            <sidebar-link :is-active="!!$page.props.isCommandsActive"
                          href="/">
                <command-line-icon class="h-5 w-5 mr-2" />
                <span>
                    Commands
                </span>
            </sidebar-link>
            <sidebar-link :is-active="!!$page.props.isCreateCommandActive"
                          href="/command/create">
                <plus-icon class="h-5 w-5 mr-2" />
                <span>
                    Create Command
                </span>
            </sidebar-link>
            <sidebar-title>
                Users
            </sidebar-title>
            <sidebar-link :is-active="!!$page.props.isUsersActive"
                          href="/user">
                <user-icon class="h-5 w-5 mr-2" />
                <span>
                    Users
                </span>
            </sidebar-link>
            <sidebar-link :is-active="!!$page.props.isCreateUserActive"
                          href="/user/create">
                <user-plus-icon class="h-5 w-5 mr-2" />
                <span>
                    Create User
                </span>
            </sidebar-link>
        </div>
        <div class="flex h-20 bg-slate-700 bg-opacity-40 text-sm text-slate-300">
            <span class="m-auto">
                &copy; {{ year }}
                <span class="text-cyan-500">
                    {{ $page.props.title }}
                </span>
            </span>
        </div>
    </div>
    <div class="md:ml-60">
        <header class="flex items-center bg-white h-20 shadow-sm px-5">
            <a class="md:hidden"
               href="#"
               @click.prevent="isSidebarOpen = true">
                <bars3-icon class="h-6 w-6" />
            </a>
            <Menu as="div" class="relative ml-auto">
                <div>
                    <MenuButton class="w-11">
                        <img class="rounded-full max-w-full" :src="$page.props.auth.gravatar">
                    </MenuButton>
                </div>
                <transition enter-active-class="transition ease-out duration-100"
                            enter-from-class="transform opacity-0 scale-95"
                            enter-to-class="transform opacity-100 scale-100"
                            leave-active-class="transition ease-in duration-75"
                            leave-from-class="transform opacity-100 scale-100"
                            leave-to-class="transform opacity-0 scale-95">
                    <MenuItems class="origin-top-right absolute right-0 mt-2 w-56 rounded-sm shadow-lg bg-white ring-1 ring-black ring-opacity-5 focus:outline-none">
                        <div class="py-1">
                            <MenuItem v-slot="{ active }">
                                <inertia-link :class="[active ? 'bg-slate-50 text-slate-800' : 'text-slate-600', 'block px-4 py-2 text-sm']"
                                              :href="`/user/edit?id=${$page.props.auth.user.id}`">
                                    Edit Profile
                                </inertia-link>
                            </MenuItem>
                            <MenuItem v-slot="{ active }">
                                <inertia-link :class="[active ? 'bg-slate-50 text-slate-800' : 'text-slate-600', 'block px-4 py-2 text-sm']"
                                              href="/logout">
                                    Logout
                                </inertia-link>
                            </MenuItem>
                        </div>
                    </MenuItems>
                </transition>
            </Menu>
        </header>
        <main class="content overflow-y-auto" scroll-region>
            <slot></slot>
        </main>
    </div>
    <!-- eslint-enable max-len -->
</template>

<script>
import {
    Menu,
    MenuButton,
    MenuItem,
    MenuItems
} from '@headlessui/vue';

import {
    Bars3Icon,
    CommandLineIcon,
    HomeIcon,
    PlusIcon,
    UserIcon,
    UserPlusIcon
} from '@heroicons/vue/24/outline';

import { ref, onUnmounted } from 'vue';
import { Inertia } from '@inertiajs/inertia';
import FlashMessages from './FlashMessages.vue';
import SidebarTitle from './SidebarTitle.vue';
import SidebarLink from './SidebarLink.vue';

export default {
    components: {
        Menu,
        MenuButton,
        MenuItem,
        MenuItems,
        Bars3Icon,
        CommandLineIcon,
        HomeIcon,
        PlusIcon,
        UserIcon,
        UserPlusIcon,
        FlashMessages,
        SidebarTitle,
        SidebarLink
    },

    setup() {
        const isSidebarOpen = ref(false);
        const year = ref(new Date().getFullYear());

        onUnmounted(
            Inertia.on('navigate', () => {
                isSidebarOpen.value = false;
            })
        );

        return {
            isSidebarOpen,
            year
        };
    }
};
</script>
