module.exports = {
    parser: 'vue-eslint-parser',
    parserOptions: {
        parser: '@babel/eslint-parser',
        requireConfigFile: false
    },
    extends: [
        'airbnb-base',
        'plugin:vue/vue3-recommended'
    ],
    env: {
        browser: true,
        commonjs: true,
        es6: true,
        jquery: true,
        node: true
    },
    globals: {
        _: 'readonly'
    },
    settings: {
        'import/core-modules': ['webpack']
    },
    rules: {
        'arrow-parens': ['error', 'as-needed'],
        'comma-dangle': ['error', 'never'],
        'global-require': 'off',
        'import/extensions': ['error', 'never', { vue: 'always' }],
        'import/no-unresolved': 'off',
        'import/no-extraneous-dependencies': ['error', { devDependencies: true }],
        indent: ['error', 4],
        'object-shorthand': ['error', 'always', { avoidQuotes: false }],
        'no-console': 'off',
        'no-empty': ['error', { allowEmptyCatch: true }],
        'no-multi-assign': 'off',
        'no-param-reassign': ['error', { props: false }],
        'no-plusplus': ['error', { allowForLoopAfterthoughts: true }],
        'no-underscore-dangle': 'off',
        'no-unused-vars': 'off',
        'vue/html-closing-bracket-newline': 'off',
        'vue/html-indent': ['error', 4],
        'vue/html-self-closing': ['error', { html: { normal: 'never' } }],
        'vue/max-attributes-per-line': ['error', {
            singleline: 2,
            multiline: { allowFirstLine: true }
        }]
    }
};
