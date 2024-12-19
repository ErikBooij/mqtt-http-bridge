import globals from 'globals';
import pluginJs from '@eslint/js';
import tseslint from 'typescript-eslint';
import pluginReact from 'eslint-plugin-react';


/** @type {import('eslint').Linter.Config[]} */
export default [
  { files: [ '**/*.{js,mjs,cjs,ts,jsx,tsx}' ] },
  { languageOptions: { globals: globals.browser } },
  pluginJs.configs.recommended,
  ...tseslint.configs.recommended,
  pluginReact.configs.flat.recommended,
  {
    rules: {
      // Enforce spaces within curly braces
      'object-curly-spacing': [ 'error', 'always' ],
      // Enforce spaces within square brackets
      'array-bracket-spacing': [ 'error', 'always' ],
      // Enforce single quotes wherever possible
      'quotes': [ 'error', 'single', { 'avoidEscape': true, 'allowTemplateLiterals': true } ],
    },
  },
];
