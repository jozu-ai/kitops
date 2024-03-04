/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/**/*.md',
    './.vitepress/**/*.{ts,vue,css}'
  ],
  theme: {
    extend: {
      fontFamily: {
        heading: ['Titillium Web', 'sans-serif'],
        brand: ['Major Mono Display', 'sans-serif'],
        sans: ['Atkinson Hyperlegible', 'sans-serif']
      },

      colors: {
        'off-white': '#ECECEC',
        gold: '#FFAF52',
        night: '#121212',
        salmon: '#FFA3AF',
        cornflower: '#7A8CFF',
        ash: '#BED8D4',
        aero: '#00C3E2',
        lavender: '#BC7EFF',
        mustard: '#D4BF40',
        redish: '#DB636A',
        zomp: '#659E8B',
        tea: '#DEF2C8',
        gray: {
          '02': '#4D4D4F',
          '05': '#6A6A6A'
        }
      },

      spacing: {
        22: '5.5rem'
      }
    },
  },
  plugins: [
    // Extract the colors to custom css variables for convenience with the theme config values
    function({ addBase, theme }) {
      function extractColorVars(colorObj, colorGroup = '') {
        return Object.keys(colorObj).reduce((vars, colorKey) => {
          const value = colorObj[colorKey];

          const newVars =
            typeof value === 'string'
              ? { [`--color${colorGroup}-${colorKey}`]: value }
              : extractColorVars(value, `-${colorKey}`);

          return { ...vars, ...newVars };
        }, {});
      }

      addBase({
        ':root': extractColorVars(theme('colors')),
      });
    },

    // hocus variant = hover OR focus
    function({ addVariant }) {
      addVariant('hocus', ['&:hover', '&:focus'])
    },

    // xs: media breakpoint for mobiles from 0 to `screens.sm`
    function ({ addVariant }) {
      addVariant('xs', "@media screen and (max-width: theme('screens.sm'))")
    },
  ],
}

