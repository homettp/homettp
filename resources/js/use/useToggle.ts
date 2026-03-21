import { ref } from 'vue';

type ToggleKey<T> = keyof T | ((item: T) => string | number);

export default <T extends Record<string, any>>(key: ToggleKey<T> = 'id' as keyof T) => {
    const toggled = ref<Array<string | number>>([]);
    const itemKey = (item: T) => (typeof key === 'function'
        ? key(item)
        : item[key]) as string | number;
    const isShow = (item: T) => toggled.value.indexOf(itemKey(item)) !== -1;

    const toggle = (item: T) => {
        const value = itemKey(item);
        const index = toggled.value.indexOf(value);

        if (index === -1) {
            toggled.value.push(value);
        } else {
            toggled.value.splice(index, 1);
        }
    };

    return {
        isShow,
        toggle
    };
};
