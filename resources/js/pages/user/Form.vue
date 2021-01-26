<template>
    <div class="user__form">
        <ol class="breadcrumb">
            <li class="breadcrumb-item">
                <inertia-link href="/">
                    Home
                </inertia-link>
            </li>
            <li class="breadcrumb-item">
                <inertia-link href="/user">
                    Users
                </inertia-link>
            </li>
            <li class="breadcrumb-item active">
                {{ $metaInfo.title }}
            </li>
        </ol>
        <div class="card">
            <div class="card-header">
                <svg class="bi"
                     width="1em"
                     height="1em"
                     fill="currentColor">
                    <use :xlink:href="icon(iconName)" />
                </svg>
                <span>
                    {{ $metaInfo.title }}
                </span>
            </div>
            <form @submit.prevent="form.post(url)">
                <div class="card-body">
                    <div class="form-group row">
                        <label class="col-sm-2 col-form-label required" for="username">
                            Username
                        </label>
                        <div class="col-sm-10">
                            <input id="username"
                                   v-model="form.username"
                                   class="form-control"
                                   :class="{'is-invalid': form.errors.username}"
                                   type="text"
                                   name="username"
                                   placeholder="Username"
                                   required>
                            <span v-if="form.errors.username" class="invalid-feedback">
                                {{ form.errors.username[0] }}
                            </span>
                        </div>
                    </div>
                    <div class="form-group row">
                        <label class="col-sm-2 col-form-label required" for="email">
                            E-mail
                        </label>
                        <div class="col-sm-10">
                            <input id="email"
                                   v-model="form.email"
                                   class="form-control"
                                   :class="{'is-invalid': form.errors.email}"
                                   type="email"
                                   name="email"
                                   placeholder="E-mail"
                                   required>
                            <span v-if="form.errors.email" class="invalid-feedback">
                                {{ form.errors.email[0] }}
                            </span>
                        </div>
                    </div>
                    <div class="form-group row">
                        <label class="col-sm-2 col-form-label"
                               :class="{required: isNew}"
                               for="password">
                            Password
                        </label>
                        <div class="col-sm-10">
                            <input id="password"
                                   v-model="form.password"
                                   class="form-control"
                                   :class="{'is-invalid': form.errors.password}"
                                   type="password"
                                   name="password"
                                   placeholder="Password"
                                   :required="isNew">
                            <span v-if="form.errors.password" class="invalid-feedback">
                                {{ form.errors.password[0] }}
                            </span>
                        </div>
                    </div>
                    <div class="row">
                        <div class="col-sm-10 offset-sm-2">
                            <div class="form-check">
                                <input id="remember"
                                       v-model="form.is_enabled"
                                       class="form-check-input"
                                       type="checkbox">
                                <label class="form-check-label" for="remember">
                                    Is Enabled?
                                </label>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="card-footer">
                    <div class="row">
                        <div class="col-md-10 offset-md-2">
                            <button class="btn btn-primary" type="submit">
                                {{ $metaInfo.title }}
                            </button>
                        </div>
                    </div>
                </div>
            </form>
        </div>
        <div v-if="!isNew" class="form__delete card">
            <div class="card-header">
                <svg class="bi"
                     width="1em"
                     height="1em"
                     fill="currentColor">
                    <use :xlink:href="icon('trash')" />
                </svg>
                <span>
                    Delete User
                </span>
            </div>
            <form @submit.prevent="submitDelete">
                <div class="card-body">
                    Are you sure you want to delete the user?
                </div>
                <div class="card-footer">
                    <button class="btn btn-danger" type="submit">
                        Delete User
                    </button>
                </div>
            </form>
        </div>
    </div>
</template>

<script>
import Layout from '../../common/Layout.vue';

export default {
    props: {
        user: {
            type: Object,
            required: true
        }
    },

    layout: Layout,

    metaInfo() {
        return {
            title: this.isNew
                ? 'Create User'
                : 'Edit User'
        };
    },

    data() {
        return {
            form: this.$inertia.form({
                username: this.user.username,
                email: this.user.email,
                password: '',
                is_enabled: this.user.is_enabled
            })
        };
    },

    computed: {
        isNew() {
            return this.user.id === 0;
        },

        iconName() {
            return this.isNew
                ? 'person-plus'
                : 'person';
        },

        url() {
            return this.isNew
                ? '/user/create'
                : `/user/edit?id=${this.user.id}`;
        }
    },

    methods: {
        submitDelete() {
            this.$inertia.delete(`/user/delete?id=${this.user.id}`);
        }
    }
};
</script>
