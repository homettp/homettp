import { DateTime } from 'luxon';
import Vue from 'vue';
import VueMeta from 'vue-meta';
import { createInertiaApp, InertiaLink } from '@inertiajs/inertia-vue';
import { InertiaProgress } from '@inertiajs/progress';
import Tooltip from './common/Tooltip';

window._ = require('lodash');

try {
    window.Popper = require('popper.js').default;
    window.$ = window.jQuery = require('jquery');

    require('bootstrap');
} catch (e) {}

Vue.config.productionTip = false;
Vue.use(VueMeta);

InertiaProgress.init();
Vue.component('InertiaLink', InertiaLink);
Vue.directive('Tooltip', Tooltip);

Vue.filter('date', value => {
    const date = DateTime.fromISO(value);

    if (!date.isValid) {
        return value;
    }

    return date.toLocaleString(DateTime.DATETIME_MED);
});

Vue.mixin({
    methods: {
        icon(name) {
            return `${this.$page.props.icons}#${name}`;
        }
    }
});

createInertiaApp({
    // eslint-disable-next-line import/no-dynamic-require
    resolve: name => require(`./pages/${name}.vue`),
    setup({ el, app, props }) {
        new Vue({
            metaInfo: {
                titleTemplate: title => (title ? `${title} - ${props.title}` : props.title)
            },

            render: h => h(app, props)
        }).$mount(el);
    }
});
