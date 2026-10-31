<template>
    <inertia-head :title="title
        ? `${title} - ${$page.props.title}`
        : `${$page.props.title}`">
        <slot />
    </inertia-head>
</template>

<script setup lang="ts">
const {
    title = ''
} = defineProps<{
    title?: string
}>();
</script>
