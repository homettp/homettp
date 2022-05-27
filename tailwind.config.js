const defaultTheme = require('tailwindcss/defaultTheme');

module.exports = {
    content: [
        './resources/**/*.gohtml',
        './resources/**/*.js',
        './resources/**/*.vue'
    ],
    theme: {
        extend: {
            fontFamily: {
                sans: ['Sarabun', ...defaultTheme.fontFamily.sans]
            }
        }
    },
    plugins: [
        require('@tailwindcss/forms')
    ]
};
