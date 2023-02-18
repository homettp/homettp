import { DateTime } from 'luxon';
import { createApp, h } from 'vue';
import { createInertiaApp, Link } from '@inertiajs/vue3';
import AppTitle from './common/AppTitle.vue';

createInertiaApp({
    // eslint-disable-next-line import/no-dynamic-require
    resolve: name => require(`./pages/${name}`),
    setup({
        el, App, props, plugin
    }) {
        createApp({ render: () => h(App, props) })
            .use(plugin)
            .component('InertiaLink', Link)
            .component('AppTitle', AppTitle)
            .mixin({
                methods: {
                    copy(inputEl) {
                        navigator.clipboard.writeText(inputEl.value).then(() => {
                            document.dispatchEvent(new CustomEvent('flash', {
                                detail: {
                                    flash: 'Copied successfully.'
                                }
                            }));
                        });
                    },

                    date(value) {
                        const date = DateTime.fromISO(value);

                        if (!date.isValid) {
                            return value;
                        }

                        return date.toLocaleString(DateTime.DATETIME_MED);
                    }
                }
            })
            .mount(el);
    }
});
