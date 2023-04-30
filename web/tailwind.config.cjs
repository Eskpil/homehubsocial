/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
    colors: {
      "white": "#FFFFFF",

      "primary-background": "#FFF3E9",

      "secondary-background": "#8C5538",
      "secondary-text": "#C8AE9F",

      "accent": "#392827",
      "accent-text": "#BE9868",
    }
  },
  plugins: [],
}
