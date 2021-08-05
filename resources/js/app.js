import { DateTime } from 'luxon';
import { createApp, h } from 'vue';
import { createInertiaApp, InertiaLink, InertiaHead } from '@inertiajs/inertia-vue3';
import { InertiaProgress } from '@inertiajs/progress';
import Tooltip from './common/Tooltip';

window._ = require('lodash');

try {
    window.Popper = require('popper.js').default;
    window.$ = window.jQuery = require('jquery');

    require('bootstrap');
} catch (e) {}

InertiaProgress.init();

createInertiaApp({
    // eslint-disable-next-line import/no-dynamic-require
    resolve: name => require(`./Pages/${name}`),
    setup({
        el, App, props, plugin
    }) {
        props.titleCallback = title => (title
            ? `${title} - ${props.initialPage.props.title}`
            : props.initialPage.props.title);

        createApp({ render: () => h(App, props) })
            .use(plugin)
            .component('InertiaLink', InertiaLink)
            .component('InertiaHead', InertiaHead)
            .directive('Tooltip', Tooltip)
            .mixin({
                methods: {
                    date(value) {
                        const date = DateTime.fromISO(value);

                        if (!date.isValid) {
                            return value;
                        }

                        return date.toLocaleString(DateTime.DATETIME_MED);
                    },

                    icon(name) {
                        return `${this.$page.props.icons}#${name}`;
                    }
                }
            })
            .mount(el);
    }
});
