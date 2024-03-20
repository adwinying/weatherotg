import type { Config } from 'tailwindcss'
import daisyui from 'daisyui'
import animate from 'tailwindcss-animate'

export default {
  content: ['templates/**/*.templ'],
  theme: {
    extend: {},
  },
  plugins: [
    daisyui,
    animate,
  ],
} satisfies Config
