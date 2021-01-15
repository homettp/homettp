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
export default {
    data() {
        return {
            message: ''
        };
    },

    watch: {
        '$page.props.flash'(val) {
            if (val) {
                this.show(val);
            }
        }
    },

    created() {
        this.$parent.$on('flash', message => {
            this.show(message);
        });
    },

    methods: {
        show(message) {
            this.message = message;
            $(this.$el).toast('show');
        }
    }
};
</script>
