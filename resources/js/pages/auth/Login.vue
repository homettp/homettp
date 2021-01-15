<template>
    <div class="auth">
        <div class="container">
            <header class="auth__header">
                <inertia-link class="auth__logo" href="/login">
                    <svg class="bi"
                         width="1em"
                         height="1em"
                         fill="currentColor">
                        <use :xlink:href="icon('house')" />
                    </svg>
                    <span>
                        {{ title }}
                    </span>
                </inertia-link>
            </header>
            <div class="card">
                <div class="card-header">
                    <svg class="bi"
                         width="1em"
                         height="1em"
                         fill="currentColor">
                        <use :xlink:href="icon('box-arrow-in-right')" />
                    </svg>
                    <span>
                        {{ $metaInfo.title }}
                    </span>
                </div>
                <form @submit.prevent="form.post('/login')">
                    <div class="card-body">
                        <div class="form-group">
                            <label class="form-control-label required" for="username_or_email">
                                Username / E-mail
                            </label>
                            <input id="username_or_email"
                                   v-model="form.username_or_email"
                                   class="form-control"
                                   :class="{'is-invalid': form.errors.username_or_email}"
                                   type="text"
                                   name="username_or_email"
                                   placeholder="Username / E-mail"
                                   required>
                            <span v-if="form.errors.username_or_email" class="invalid-feedback">
                                {{ form.errors.username_or_email[0] }}
                            </span>
                        </div>
                        <div class="form-group">
                            <label class="form-control-label required" for="password">
                                Password
                            </label>
                            <input id="password"
                                   v-model="form.password"
                                   class="form-control"
                                   :class="{'is-invalid': form.errors.password}"
                                   type="password"
                                   name="password"
                                   placeholder="Password"
                                   required>
                            <span v-if="form.errors.password" class="invalid-feedback">
                                {{ form.errors.password[0] }}
                            </span>
                        </div>
                        <div class="form-check">
                            <input id="remember"
                                   v-model="form.remember"
                                   class="form-check-input"
                                   type="checkbox">
                            <label class="form-check-label" for="remember">
                                Remember Me
                            </label>
                        </div>
                    </div>
                    <div class="card-footer">
                        <button class="btn btn-primary" type="submit">
                            {{ $metaInfo.title }}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    props: {
        title: {
            type: String,
            required: true
        }
    },

    metaInfo() {
        return {
            title: 'Login'
        };
    },

    data() {
        return {
            form: this.$inertia.form({
                username_or_email: '',
                password: '',
                remember: true
            })
        };
    }
};
</script>
