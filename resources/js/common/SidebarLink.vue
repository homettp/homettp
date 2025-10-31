<template>
    <!-- eslint-disable max-len -->
    <inertia-link class="flex items-center px-5 py-3.5 text-slate-300 hover:text-cyan-500 hover:bg-slate-700/40"
                  :class="{'text-cyan-500 bg-slate-700/40': isActive}"
                  :href="href">
        <slot></slot>
    </inertia-link>
    <!-- eslint-enable max-len -->
</template>

<script setup>
import { defineProps } from 'vue';

defineProps({
    isActive: {
        type: Boolean,
        required: true
    },

    href: {
        type: String,
        required: true
    }
});
</script>
