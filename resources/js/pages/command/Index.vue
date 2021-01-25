<template>
    <div class="command__index">
        <ol class="breadcrumb">
            <li class="breadcrumb-item">
                <inertia-link href="/">
                    Home
                </inertia-link>
            </li>
            <li class="breadcrumb-item active">
                {{ $metaInfo.title }}
            </li>
        </ol>
    </div>
</template>

<script>
import Layout from '../../common/Layout.vue';

export default {
    layout: Layout,

    metaInfo() {
        return {
            title: 'Commands'
        };
    }
};
</script>
