<template>
    <inertia-head :title="title ? `${title} - ${$page.props.title}` : $page.props.title">
        <slot></slot>
    </inertia-head>
</template>

<script>
import { Head } from '@inertiajs/vue3';

export default {
    components: {
        InertiaHead: Head
    },

    props: {
        title: {
            type: String,
            default: ''
        }
    }
};
</script>
