module.exports = {
    root: true,
    env: {
        browser: true,
        node: true,
    },
    parser: "vue-eslint-parser",
    parserOptions: {
        project: './tsconfig.json', // Specify it only for TypeScript files
        "parser": {
            // Script parser for `<script>`
            "js": "espree",

            // Script parser for `<script lang="ts">`
            "ts": "@typescript-eslint/parser",

            // Script parser for vue directives (e.g. `v-if=` or `:attribute=`)
            // and vue interpolations (e.g. `{{variable}}`).
            // If not specified, the parser determined by `<script lang ="...">` is used.
            "<template>": "espree",
            project: './tsconfig.json'
        },
        tsconfigRootDir: __dirname,
        createDefaultProgram: true,
        ecmaVersion: 2020,
        sourceType: "module",
    },
    extends: [
        'plugin:nuxt/recommended',
        'plugin:vue/recommended',
        "eslint:recommended",
        "plugin:@typescript-eslint/recommended",
        "plugin:@typescript-eslint/recommended-requiring-type-checking",
    ],
    plugins: [
        "@typescript-eslint",
    ],
    overrides: [
        {
            files: ['*.ts', '*.tsx'], // Your TypeScript files extension

            // As mentioned in the comments, you should extend TypeScript plugins here,
            // instead of extending them outside the `overrides`.
            // If you don't want to extend any rules, you don't need an `extends` attribute.
            extends: [
                'plugin:@typescript-eslint/recommended',
                'plugin:@typescript-eslint/recommended-requiring-type-checking',
            ],
        }
    ],
    rules: {
        '@typescript-eslint/no-unsafe-call': 'off',
        'vue/multi-word-component-names': 'off',
        // allow-dangle (syncs w/ laravel).
        'comma-dangle': ['error', 'always-multiline'],

        // Only allow curly braces of multi-line statements.
        curly: ['error', 'multi'],

        // for sanity - allow people to write console.log - DO NOT CHECK IN
        'no-console': 'off',
        'no-debugger': 'off',

        // do not require default for properties
        'vue/require-default-prop': 'off',

        // lets allow content on the same line for short html
        'vue/singleline-html-element-content-newline': 0,

        // under discussion
        camelcase: 'off',

        // Allows for lexical declarations within case/default clauses.
        'no-case-declarations': 'off',

        // this restrictions functions and lib to be defined before they are used
        // this is impossible to have on with model-typer
        'no-use-before-define': 'off',

        // note you must disable the base rule as it can report incorrect errors
        'no-useless-constructor': 'off',

        // https://eslint.vuejs.org/rules/no-v-html.html
        'vue/no-v-html': 'off',
    },
}
