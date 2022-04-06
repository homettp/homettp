<template>
    <Head :title="title ? `${title} - ${$page.props.title}` : $page.props.title">
        <slot></slot>
    </Head>
</template>

<script>
import { Head } from '@inertiajs/inertia-vue3';

export default {
    components: {
        Head
    },

    props: {
        title: {
            type: String,
            default: ''
        }
    }
};
</script>
