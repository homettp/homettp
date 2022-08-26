import { DateTime } from 'luxon';
import { createApp, h } from 'vue';
import { createInertiaApp, InertiaLink } from '@inertiajs/inertia-vue3';
import { InertiaProgress } from '@inertiajs/progress';
import AppTitle from './common/AppTitle.vue';

InertiaProgress.init();

createInertiaApp({
    // eslint-disable-next-line import/no-dynamic-require
    resolve: name => require(`./pages/${name}`),
    setup({
        el, App, props, plugin
    }) {
        createApp({ render: () => h(App, props) })
            .use(plugin)
            .component('InertiaLink', InertiaLink)
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
