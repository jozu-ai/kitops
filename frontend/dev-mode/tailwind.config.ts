import type { Config } from 'tailwindcss'
import baseColors from 'tailwindcss/colors'
import plugin from 'tailwindcss/plugin'

// From https://github.com/tailwindlabs/tailwindcss/issues/4690#issuecomment-1046087220
delete (baseColors as Record<string, any>).lightBlue
delete (baseColors as Record<string, any>).warmGray
delete (baseColors as Record<string, any>).trueGray
delete (baseColors as Record<string, any>).coolGray
delete (baseColors as Record<string, any>).blueGray

/** @type {Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts}'
  ],

  theme: {
    extend: {
      fontFamily: {
        sans: ['Atkinson Hyperlegible']
      },
      colors: {
        'off-white': '#ECECEC',
        night: '#121212',
        gold: '#FFAF52',
        cornflower: '#7A8CFF',
        gray: {
          '01': '#2D2D2D',
          '02': '#4D4D4D',
          '03': '#585858',
          '04': '#606060',
          '05': '#6A6A6A',
          '06': '#AAAAAA',
          '07': '#E0E0E0',
          '08': '#E8E8E8',
          '09': '#F0F0F0',
          '10': '#F5F5F5'
        },
        elevation: {
          '0': 'rgb(18, 18, 18)',
          '01': '#1F1F1F',
          '02': '#292929',
          '03': '#333333',
          '04': '#3D3D3D',
          '05': '#464646'
        },
      },
      zIndex: {
        // @ts-ignore
        '1': 1
      },
      spacing: {
        '22': '5.5rem'
      }
    },
  },

  plugins: [
    // hocus variant = hover OR focus
    plugin(function({ addVariant }) {
      addVariant('hocus', ['&:hover', '&:focus'])
    }),

    // xs: media breakpoint for mobiles from 0 to `screens.sm`
    plugin(function ({ addVariant }) {
      addVariant('xs', "@media screen and (max-width: theme('screens.sm'))")
    }),
  ],

} satisfies Config

