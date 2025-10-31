<template>
    <flash-messages />
    <div v-show="isSidebarOpen"
         class="bg-black/50 fixed inset-0 z-10 md:hidden"
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
        <div class="flex h-20 bg-slate-700/40 text-sm text-slate-300">
            <span class="m-auto">
                &copy; {{ year }}
                <span class="text-cyan-500">
                    {{ $page.props.title }}
                </span>
            </span>
        </div>
    </div>
    <div class="md:ml-60">
        <header class="flex items-center bg-white dark:bg-slate-700 h-20 shadow-xs px-5 dark:bg-slate-700">
            <button class="md:hidden dark:text-slate-300"
                    type="button"
                    @click="isSidebarOpen = true">
                <bars3-icon class="h-6 w-6" />
            </button>
            <button v-if="theme === 'system'"
                    class="ml-auto dark:text-slate-300"
                    type="button"
                    @click="theme = 'dark'">
                <computer-desktop-icon class="h-6 w-6" />
            </button>
            <button v-else-if="theme === 'dark'"
                    class="ml-auto text-slate-300"
                    type="button"
                    @click="theme = 'light'">
                <moon-icon class="h-6 w-6" />
            </button>
            <button v-else-if="theme === 'light'"
                    class="ml-auto"
                    type="button"
                    @click="theme = 'system'">
                <sun-icon class="h-6 w-6" />
            </button>
            <MainMenu as="div" class="relative ml-4">
                <div>
                    <MenuButton class="w-11">
                        <img class="rounded-full max-w-full" :src="$page.props.auth.gravatar">
                    </MenuButton>
                </div>
                <transition enter-active-class="transition ease-out duration-100"
                            enter-from-class="opacity-0 scale-95"
                            enter-to-class="opacity-100 scale-100"
                            leave-active-class="transition ease-in duration-75"
                            leave-from-class="opacity-100 scale-100"
                            leave-to-class="opacity-0 scale-95">
                    <MenuItems class="origin-top-right absolute right-0 mt-2 w-56 rounded-xs shadow-lg bg-white dark:bg-slate-800 ring-1 ring-black/5 focus:outline-hidden">
                        <div class="py-1">
                            <MenuItem v-slot="{ active }">
                                <button :class="[active ? 'bg-slate-50 dark:hover:bg-slate-700/40 text-slate-800 dark:text-slate-300' : 'text-slate-600 dark:text-slate-400', 'flex w-full items-center px-4 py-2 text-sm']"
                                        @click="router.visit(`/user/edit?id=${$page.props.auth.user.id}`)">
                                    Edit Profile
                                </button>
                            </MenuItem>
                            <MenuItem v-slot="{ active }">
                                <button :class="[active ? 'bg-slate-50 dark:hover:bg-slate-700/40 text-slate-800 dark:text-slate-300' : 'text-slate-600 dark:text-slate-400', 'flex w-full items-center px-4 py-2 text-sm']"
                                        @click="router.visit('/logout')">
                                    Logout
                                </button>
                            </MenuItem>
                        </div>
                    </MenuItems>
                </transition>
            </MainMenu>
        </header>
        <main class="content overflow-y-auto" scroll-region>
            <slot></slot>
        </main>
    </div>
    <!-- eslint-enable max-len -->
</template>

<script setup>
import {
    Menu as MainMenu,
    MenuButton,
    MenuItem,
    MenuItems
} from '@headlessui/vue';

import {
    Bars3Icon,
    CommandLineIcon,
    ComputerDesktopIcon,
    HomeIcon,
    MoonIcon,
    PlusIcon,
    SunIcon,
    UserIcon,
    UserPlusIcon
} from '@heroicons/vue/24/outline';

import { ref, onUnmounted, watch } from 'vue';
import { router } from '@inertiajs/vue3';
import FlashMessages from './FlashMessages.vue';
import SidebarTitle from './SidebarTitle.vue';
import SidebarLink from './SidebarLink.vue';

const theme = ref(
    window.currentTheme()
);

window.setTheme(theme.value);

watch(theme, () => {
    if (theme.value === 'system') {
        localStorage.removeItem('theme');
    } else {
        localStorage.theme = theme.value;
    }

    window.setTheme(theme.value);
});

const isSidebarOpen = ref(false);
const year = ref(new Date().getFullYear());

onUnmounted(
    router.on('navigate', () => {
        isSidebarOpen.value = false;
    })
);
</script>
