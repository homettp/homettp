import { ref } from 'vue';

export default (key = 'id') => {
    const toggled = ref([]);
    const isShow = item => toggled.value.indexOf(item[key]) !== -1;

    const toggle = item => {
        const index = toggled.value.indexOf(item[key]);

        if (index === -1) {
            toggled.value.push(item[key]);
        } else {
            toggled.value.splice(index, 1);
        }
    };

    return {
        isShow,
        toggle
    };
};
