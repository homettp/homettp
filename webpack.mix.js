const mix = require('laravel-mix');

/*
 |--------------------------------------------------------------------------
 | Mix Asset Management
 |--------------------------------------------------------------------------
 |
 | Mix provides a clean, fluent API for defining some Webpack build steps
 | for your Laravel application. By default, we are compiling the Sass
 | file for the application as well as bundling up all the JS files.
 |
 */

mix.options({ processCssUrls: false })
    .js('resources/js/app.js', 'public/js')
    .postCss('resources/css/app.css', 'css', [
        require('tailwindcss')
    ])
    .copy('resources/images/favicon.ico', 'public/favicon.ico')
    .vue()
    .version()
    .sourceMaps()
    .setPublicPath('public');
