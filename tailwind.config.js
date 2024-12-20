const defaultTheme = require('tailwindcss/defaultTheme')

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [ './src/frontend/templates/**/*', './src/frontend/js/**/*', './src/frontend/index.html' ],
  theme: {
    extend: {
      fontFamily: {
        sans: [ 'Inter var', ...defaultTheme.fontFamily.sans ],
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}
