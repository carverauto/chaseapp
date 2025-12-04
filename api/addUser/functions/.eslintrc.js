module.exports = {
  root: true,
  env: {
    browser: true,
    node: true,
  },
  parserOptions: {
    parser: 'babel-eslint',
		ecmaVersion: 8,
  },
  plugins: [
  ],
  // add your custom rules here
  rules: {
    'comma-dangle': [ 'error', 'always-multiline' ],
    curly: [ 'error', 'multi' ],
    'no-console': 'off',
    camelcase: 'off',
    'array-bracket-spacing': [ 'error', 'always' ],
  },
}
