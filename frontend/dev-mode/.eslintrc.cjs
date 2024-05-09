/* eslint-env node */
require('@rushstack/eslint-patch/modern-module-resolution')

module.exports = {
  root: true,
  'extends': [
    'plugin:vue/vue3-essential',
    'eslint:recommended',
    '@vue/eslint-config-typescript'
  ],
  parserOptions: {
    ecmaVersion: 'latest'
  },
  plugins: [
    'import',
  ],
  rules: {
    quotes: ['error', 'single', { 'avoidEscape': true }],
    'no-unused-vars': 'off',
    '@typescript-eslint/no-unused-vars': ['error', { 'argsIgnorePattern': '^_' }],
    semi: ['error', 'never'],
    'prefer-promise-reject-errors': 'off',
    'vue/max-len': ['warn', {
      code: 120,
      template: 300
    }],
    'vue/component-tags-order': ['error', {
      order: ['script', 'template', 'style']
    }],
    'vue/multi-word-component-names': 'off',
    'vue/html-indent': ['error', 2, {
      attribute: 1,
      baseIndent: 0,
      closeBracket: 1,
      alignAttributesVertically: false,
      ignores: [],
    }],
    'vue/html-closing-bracket-spacing': 'error',
    'vue/html-closing-bracket-newline': ['error', {
      singleline: 'never',
      multiline: 'never'
    }],
    'vue/no-reserved-component-names': 'off',
    'import/order': ['error', {
      'newlines-between': 'always',
      alphabetize: {
        order: 'asc',
        caseInsensitive: true,
      },
    }],
  },
  globals: {
    // env globals
    process: true,
    // Vue Macros
    defineProps: 'readonly',
    defineEmits: 'readonly',
    defineExpose: 'readonly',
    withDefaults: 'readonly',
    // env vars
    config: 'readable',
  }
}
