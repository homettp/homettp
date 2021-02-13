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
    .sass('resources/sass/app.scss', 'css')
    .copy('node_modules/bootstrap-icons/bootstrap-icons.svg', 'public/images')
    .copy('resources/images/favicon.ico', 'public/favicon.ico')
    .webpackConfig({ output: { chunkFilename: 'js/[name].js?id=[chunkhash]' } })
    .vue()
    .version()
    .sourceMaps()
    .setPublicPath('public');
