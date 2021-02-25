import { DateTime } from 'luxon';
import Vue from 'vue';
import VueMeta from 'vue-meta';
import { App, plugin } from '@inertiajs/inertia-vue';
import { InertiaProgress } from '@inertiajs/progress';
import Tooltip from './common/Tooltip';

window._ = require('lodash');

try {
    window.Popper = require('popper.js').default;
    window.$ = window.jQuery = require('jquery');

    require('bootstrap');
} catch (e) {}

InertiaProgress.init();

Vue.config.productionTip = false;
Vue.use(plugin);
Vue.use(VueMeta);

Vue.directive('tooltip', Tooltip);

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

const el = document.getElementById('app');
const initialPage = JSON.parse(el.dataset.page);

new Vue({
    metaInfo: {
        titleTemplate: title => (title ? `${title} - ${initialPage.props.title}` : initialPage.props.title)
    },

    render: h => h(App, {
        props: {
            initialPage,
            resolveComponent: name => import(`./pages/${name}.vue`).then(module => module.default)
        }
    })
}).$mount(el);
