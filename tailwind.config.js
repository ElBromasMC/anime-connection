const defaultTheme = require('tailwindcss/defaultTheme')

/** @type {import('tailwindcss').Config} */
export const content = [
  './view/**/*.templ',
];
export const theme = {
  extend: {
    fontFamily: {
      mono: ['Courier Prime', 'monospace'],
      sans: ['"Nunito Sans"', ...defaultTheme.fontFamily.sans],
    },
    colors: {
      navy: '#1d2747',
      azure: '#39a0ed',
      chalky: '#f9f9f9',
      livid: '#4c6085',
      darkblue: "#0f151b"
    },
  },
};
export const plugins = [
];
export const corePlugins = { preFlight: true };
