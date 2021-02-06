<template>
    <div class="toast fade hide">
        <div class="toast-body d-flex align-items-center">
            <svg class="bi"
                 width="1em"
                 height="1em"
                 fill="currentColor">
                <use :xlink:href="icon('check-circle')" />
            </svg>
            <span>
                {{ message }}
            </span>
        </div>
    </div>
</template>

<script>
import { Inertia } from '@inertiajs/inertia';

export default {
    data() {
        return {
            message: ''
        };
    },

    created() {
        this.$parent.$on('flash', message => {
            this.show(message);
        });
    },

    mounted() {
        this.$once(
            'hook:destroyed',
            Inertia.on('success', () => { this.show(this.$page.props.flash); })
        );
    },

    methods: {
        show(message) {
            if (message) {
                this.message = message;
                $(this.$el).toast('show');
            }
        }
    }
};
</script>
