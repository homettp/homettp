<template>
    <div class="command__form layout__form">
        <ol class="breadcrumb">
            <li class="breadcrumb-item">
                <inertia-link href="/">
                    Home
                </inertia-link>
            </li>
            <li class="breadcrumb-item">
                <inertia-link href="/">
                    Commands
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
                        <label class="col-sm-2 col-form-label required" for="name">
                            Name
                        </label>
                        <div class="col-sm-10">
                            <input id="name"
                                   v-model="form.name"
                                   class="form-control"
                                   :class="{'is-invalid': form.errors.name}"
                                   type="text"
                                   name="name"
                                   placeholder="Name"
                                   required>
                            <span v-if="form.errors.name" class="invalid-feedback">
                                {{ form.errors.name[0] }}
                            </span>
                        </div>
                    </div>
                    <div class="form-group row">
                        <label class="col-sm-2 col-form-label required" for="image">
                            Image
                        </label>
                        <div class="col-sm-10">
                            <select id="image"
                                    v-model="form.image"
                                    class="form-control"
                                    :class="{'is-invalid': form.errors.image}"
                                    name="image"
                                    required>
                                <option v-for="commandImage in commandImages"
                                        :key="commandImage.value"
                                        :value="commandImage.value">
                                    {{ commandImage.name }}
                                </option>
                            </select>
                            <span v-if="form.errors.image" class="invalid-feedback">
                                {{ form.errors.image[0] }}
                            </span>
                        </div>
                    </div>
                    <div class="form-group row mb-0">
                        <label class="col-sm-2 col-form-label required" for="value">
                            Value
                        </label>
                        <div class="col-sm-10">
                            <input id="value"
                                   v-model="form.value"
                                   class="form-control"
                                   :class="{'is-invalid': form.errors.value}"
                                   type="text"
                                   name="value"
                                   placeholder="Value"
                                   required>
                            <span v-if="form.errors.value" class="invalid-feedback">
                                {{ form.errors.value[0] }}
                            </span>
                            <small class="form-text text-muted">
                                Payload Variable:
                                <code class="text-primary">
                                    {{ commandPayload }}
                                </code>
                            </small>
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
        <div v-if="!isNew" class="form__secondary card">
            <div class="card-header">
                <svg class="bi"
                     width="1em"
                     height="1em"
                     fill="currentColor">
                    <use :xlink:href="icon('terminal')" />
                </svg>
                <span class="mr-auto">
                    Call Command
                </span>
                <inertia-link class="command__call ml-3" :href="`/call/history?id=${command.id}`">
                    Call History
                    <svg class="bi"
                         width="1em"
                         height="1em"
                         fill="currentColor">
                        <use :xlink:href="icon('chevron-right')" />
                    </svg>
                </inertia-link>
            </div>
            <div class="card-body">
                <p>
                    Making a
                    <code class="text-primary">GET</code>
                    or
                    <code class="text-primary">POST</code>
                    request to this URL will call your command:
                </p>
                <div class="input-group">
                    <input ref="commandPath"
                           class="form-control"
                           type="text"
                           :value="commandPath"
                           readonly>
                    <div class="input-group-append">
                        <a class="btn btn-outline-primary"
                           href="#"
                           @click.prevent="copy()">
                            <svg class="bi"
                                 width="1em"
                                 height="1em"
                                 fill="currentColor">
                                <use :xlink:href="icon('files')" />
                            </svg>
                        </a>
                    </div>
                </div>
            </div>
            <div class="card-footer">
                <inertia-link class="btn btn-primary"
                              :href="`/command/refresh-token?id=${command.id}`"
                              method="put"
                              as="button">
                    Refresh Token
                </inertia-link>
            </div>
        </div>
        <div v-if="!isNew" class="form__secondary card">
            <div class="card-header">
                <svg class="bi"
                     width="1em"
                     height="1em"
                     fill="currentColor">
                    <use :xlink:href="icon('trash')" />
                </svg>
                <span>
                    Delete Command
                </span>
            </div>
            <div class="card-body">
                Are you sure you want to delete the command?
            </div>
            <div class="card-footer">
                <inertia-link class="btn btn-danger"
                              :href="`/command/delete?id=${command.id}`"
                              method="delete"
                              as="button">
                    Delete Command
                </inertia-link>
            </div>
        </div>
    </div>
</template>

<script>
import Layout from '../../common/Layout.vue';

export default {
    layout: Layout,

    props: {
        command: {
            type: Object,
            required: true
        },

        commandImages: {
            type: Array,
            required: true
        },

        commandPayload: {
            type: String,
            required: true
        },

        commandPath: {
            type: String,
            default: ''
        }
    },

    metaInfo() {
        return {
            title: this.isNew
                ? 'Create Command'
                : 'Edit Command'
        };
    },

    data() {
        return {
            form: this.$inertia.form({
                name: this.command.name,
                image: this.command.image,
                value: this.command.value
            })
        };
    },

    computed: {
        isNew() {
            return this.command.id === 0;
        },

        iconName() {
            return this.isNew
                ? 'plus-circle'
                : 'command';
        },

        url() {
            return this.isNew
                ? '/command/create'
                : `/command/edit?id=${this.command.id}`;
        }
    },

    methods: {
        copy() {
            this.$refs.commandPath.select();
            this.$refs.commandPath.setSelectionRange(0, 99999);

            document.execCommand('copy');

            this.$parent.$emit('flash', 'Copied successfully.');
        }
    }
};
</script>
