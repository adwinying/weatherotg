import type { Config } from 'tailwindcss'
import plugin from 'tailwindcss/plugin'
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
    plugin(({ addVariant }) => {
      addVariant('htmx-settling', ['&.htmx-settling', '.htmx-settling &'])
      addVariant('htmx-request',  ['&.htmx-request',  '.htmx-request &'])
      addVariant('htmx-swapping', ['&.htmx-swapping', '.htmx-swapping &'])
      addVariant('htmx-added',    ['&.htmx-added',    '.htmx-added &'])
    }),
  ],
} satisfies Config
