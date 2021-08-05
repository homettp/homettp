<template>
    <section class="d-flex">
        <div v-show="isSidebarOpen"
             class="sidebar__backdrop"
             @click.prevent="isSidebarOpen = false"></div>
        <div class="sidebar" :class="{'sidebar--open': isSidebarOpen}">
            <div class="sidebar__header">
                <inertia-link class="sidebar__logo" href="/">
                    <svg class="bi"
                         width="1em"
                         height="1em"
                         fill="currentColor">
                        <use :xlink:href="icon('house')" />
                    </svg>
                    <span>
                        {{ $page.props.title }}
                    </span>
                </inertia-link>
            </div>
            <div class="sidebar__content">
                <h1 class="sidebar__content--title">
                    Commands
                </h1>
                <ul class="nav flex-column">
                    <li class="nav-item">
                        <inertia-link class="nav-link"
                                      :class="{active: $page.props.isCommandsActive}"
                                      href="/">
                            <svg class="bi"
                                 width="1em"
                                 height="1em"
                                 fill="currentColor">
                                <use :xlink:href="icon('command')" />
                            </svg>
                            <span>
                                Commands
                            </span>
                        </inertia-link>
                        <inertia-link class="nav-link"
                                      :class="{active: $page.props.isCreateCommandActive}"
                                      href="/command/create">
                            <svg class="bi"
                                 width="1em"
                                 height="1em"
                                 fill="currentColor">
                                <use :xlink:href="icon('plus-circle')" />
                            </svg>
                            <span>
                                Create Command
                            </span>
                        </inertia-link>
                    </li>
                </ul>
                <h1 class="sidebar__content--title">
                    Users
                </h1>
                <ul class="nav flex-column">
                    <li class="nav-item">
                        <inertia-link class="nav-link"
                                      :class="{active: $page.props.isUsersActive}"
                                      href="/user">
                            <svg class="bi"
                                 width="1em"
                                 height="1em"
                                 fill="currentColor">
                                <use :xlink:href="icon('person')" />
                            </svg>
                            <span>
                                Users
                            </span>
                        </inertia-link>
                    </li>
                    <li class="nav-item">
                        <inertia-link class="nav-link"
                                      :class="{active: $page.props.isCreateUserActive}"
                                      href="/user/create">
                            <svg class="bi"
                                 width="1em"
                                 height="1em"
                                 fill="currentColor">
                                <use :xlink:href="icon('person-plus')" />
                            </svg>
                            <span>
                                Create User
                            </span>
                        </inertia-link>
                    </li>
                </ul>
            </div>
            <div class="sidebar__footer">
                <span class="m-auto">
                    &copy; {{ year }}
                    <span class="sidebar__highlight">
                        {{ $page.props.title }}
                    </span>
                </span>
            </div>
        </div>
        <div class="content">
            <header class="content__header">
                <a class="content__header--toggler"
                   href="#"
                   @click.prevent="isSidebarOpen = true">
                    <svg class="bi"
                         width="1em"
                         height="1em"
                         fill="currentColor">
                        <use :xlink:href="icon('list')" />
                    </svg>
                </a>
                <flash-message />
                <div class="content__header--user dropdown">
                    <a class="dropdown-toggle"
                       href="#"
                       data-toggle="dropdown">
                        <img class="img-fluid rounded-circle"
                             :src="$page.props.auth.gravatar">
                    </a>
                    <div class="dropdown-menu dropdown-menu-right">
                        <inertia-link class="dropdown-item"
                                      :href="`/user/edit?id=${$page.props.auth.user.id}`">
                            Edit Profile
                        </inertia-link>
                        <inertia-link class="dropdown-item" href="/logout">
                            Logout
                        </inertia-link>
                    </div>
                </div>
            </header>
            <main class="content__main">
                <slot></slot>
            </main>
        </div>
    </section>
</template>

<script>
import { ref, onUnmounted } from 'vue';
import { Inertia } from '@inertiajs/inertia';
import FlashMessage from './FlashMessage.vue';

export default {
    components: {
        FlashMessage
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
