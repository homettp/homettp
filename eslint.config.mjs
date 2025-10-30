import { defineConfig } from 'eslint/config';
import pluginVue from 'eslint-plugin-vue';
import globals from 'globals';
import vueParser from 'vue-eslint-parser';
import parser from '@babel/eslint-parser';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
// eslint-disable-next-line import/no-extraneous-dependencies
import js from '@eslint/js';
// eslint-disable-next-line import/no-extraneous-dependencies
import { FlatCompat } from '@eslint/eslintrc';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const compat = new FlatCompat({
    baseDirectory: __dirname,
    recommendedConfig: js.configs.recommended,
    allConfig: js.configs.all
});

export default defineConfig([
    ...pluginVue.configs['flat/recommended'],
    {
        extends: compat.extends('airbnb-base'),

        languageOptions: {
            globals: {
                ...globals.browser,
                ...globals.commonjs,
                ...globals.jquery,
                ...globals.node
            },

            parser: vueParser,
            sourceType: 'module',

            parserOptions: {
                parser,
                requireConfigFile: false
            }
        },

        rules: {
            'arrow-parens': ['error', 'as-needed'],
            'class-methods-use-this': 'off',
            'comma-dangle': ['error', 'never'],
            'global-require': 'off',
            'function-paren-newline': ['error', 'consistent'],
            'import/extensions': ['error', 'never'],
            'import/no-unresolved': 'off',
            'import/no-extraneous-dependencies': ['error', {
                devDependencies: true
            }],
            indent: ['error', 4],
            'keyword-spacing': ['error', {
                overrides: {
                    this: {
                        before: false
                    }
                }
            }],
            'object-shorthand': ['error', 'always', {
                avoidQuotes: false
            }],
            'no-console': 'off',
            'no-empty': ['error', {
                allowEmptyCatch: true
            }],
            'no-multi-assign': 'off',
            'no-param-reassign': ['error', {
                props: false
            }],
            'no-plusplus': ['error', {
                allowForLoopAfterthoughts: true
            }],
            'no-underscore-dangle': 'off',
            'no-unused-vars': 'off',
            'vue/first-attribute-linebreak': ['error', {
                multiline: 'beside'
            }],
            'vue/html-closing-bracket-newline': 'off',
            'vue/html-indent': ['error', 4],
            'vue/html-self-closing': ['error', {
                html: {
                    normal: 'never'
                }
            }],
            'vue/max-attributes-per-line': ['error', {
                singleline: 2
            }],
            'vue/multi-word-component-names': 'off',
            'vue/no-v-html': 'off'
        }
    }
]);
