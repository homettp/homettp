import { DateTime } from 'luxon';

export default () => {
    const date = (value: string) => {
        const d = DateTime.fromISO(value);

        if (!d.isValid) {
            return value;
        }

        return d.toLocaleString(DateTime.DATETIME_MED);
    };

    return {
        date
    };
};
