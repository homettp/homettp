<template>
    <div class="bg-white p-8 dark:bg-slate-700 dark:text-slate-300">
        <card-title>
            <slot name="title" />
        </card-title>
        <div class="relative overflow-hidden overflow-x-auto">
            <table class="w-full">
                <thead>
                    <tr>
                        <th v-for="column in props.columns"
                            :key="column"
                            class="px-6 py-3 text-left">
                            {{ column }}
                        </th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="!props.items || !props.items.length"
                        class="bg-white p-8 dark:bg-slate-700">
                        <td class="border-t border-gray-100 px-6 py-3 dark:border-slate-600"
                            :colspan="colspan">
                            No items.
                        </td>
                    </tr>
                    <tr v-for="item in props.items"
                        :key="item.id"
                        class="hover:bg-gray-50 hover:dark:bg-slate-800/50">
                        <slot name="content" :item="item" />
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, type PropType } from 'vue';
import CardTitle from './CardTitle.vue';

type RowItem = {
    id: number
};

const props = defineProps({
    columns: {
        type: Array as PropType<string[]>,
        required: true
    },

    items: {
        type: Array as PropType<RowItem[]>,
        default: () => []
    }
});

type TableViewSlots = {
    title?: () => any
    content?: (props: { item: any }) => any
};

defineSlots<TableViewSlots>();

const colspan = computed(() => props.columns.length + 1);
</script>
