/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
  theme: {
    extend: {
      colors: {
        dark_primary: '#27282d',
      },
    },
  },
  plugins: [],
  darkMode: 'class',
}
