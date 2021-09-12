const defaultTheme = require('tailwindcss/defaultTheme');
const colors = require('tailwindcss/colors');

module.exports = {
    mode: 'jit',
    purge: [
        './resources/**/*.gohtml',
        './resources/**/*.js',
        './resources/**/*.vue'
    ],
    darkMode: 'media', // or 'media' or 'class'
    theme: {
        fontFamily: {
            sans: ['Sarabun', ...defaultTheme.fontFamily.sans]
        },
        extend: {
            colors: {
                blueGray: colors.blueGray,
                cyan: colors.cyan
            }
        }
    },
    variants: {
        extend: {}
    },
    plugins: [
        require('@tailwindcss/forms')
    ]
};
