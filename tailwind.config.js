const defaultTheme = require('tailwindcss/defaultTheme');
const colors = require('tailwindcss/colors');

module.exports = {
    content: [
        './resources/**/*.gohtml',
        './resources/**/*.js',
        './resources/**/*.vue'
    ],
    theme: {
        fontFamily: {
            sans: ['Sarabun', ...defaultTheme.fontFamily.sans]
        },
        extend: {
            colors: {
                slate: colors.slate,
                cyan: colors.cyan
            }
        }
    },
    plugins: [
        require('@tailwindcss/forms')
    ]
};
