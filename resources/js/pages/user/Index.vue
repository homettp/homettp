<template>
    <div class="user__index">
        <ol class="breadcrumb">
            <li class="breadcrumb-item">
                <inertia-link href="/">
                    Home
                </inertia-link>
            </li>
            <li class="breadcrumb-item active">
                Users
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
            title: 'Users'
        };
    }
};
</script>
