module.exports = {
  root: true,
  env: {
    browser: true,
    es2020: true
  },
  extends: ['plugin:vue/essential', 'eslint:recommended', '@vue/prettier'],
  parserOptions: {
    parser: 'babel-eslint'
  },
  rules: {
    // thanks Mohit  https://stackoverflow.com/a/66618201
    'prettier/prettier': [
      1,
      {
        trailingComma: 'es5',
        //to enable single quotes
        singleQuote: true,
        semi: false
      }
    ],
    quotes: 0,
    'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'warn' : 'off'
  },
  overrides: [
    {
      files: [
        '**/__tests__/*.{j,t}s?(x)',
        '**/tests/unit/**/*.spec.{j,t}s?(x)'
      ],
      env: {
        jest: true
      }
    }
  ]
}
