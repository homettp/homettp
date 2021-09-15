<template>
    <!-- eslint-disable max-len -->
    <inertia-link class="flex items-center px-5 py-3.5 text-blueGray-300 hover:text-cyan-500 hover:bg-blueGray-700 hover:bg-opacity-40"
                  :class="{'text-cyan-500 bg-blueGray-700 bg-opacity-40': isActive}"
                  :href="href">
        <slot></slot>
    </inertia-link>
    <!-- eslint-enable max-len -->
</template>

<script>
export default {
    props: {
        isActive: {
            type: Boolean,
            required: true
        },

        href: {
            type: String,
            required: true
        }
    }
};
</script>
