import { URL, fileURLToPath } from 'node:url'

import { defineConfig } from 'vitepress'
import { getSidebarItemsFromMdFiles } from './utils.mts'
import { resolve } from 'path'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  appearance: 'force-dark',
  srcDir: 'src',
  title: 'KitOps',
  description: 'Documentation for KitOps',

  head: [
    ['link', { rel: "apple-touch-icon", sizes: "180x180", href: "/favicons/apple-touch-icon.png"}],
    ['link', { rel: "icon", type: "image/png", sizes: "32x32", href: "/favicons/favicon-32x32.png"}],
    ['link', { rel: "icon", type: "image/png", sizes: "16x16", href: "/favicons/favicon-16x16.png"}],
    ['link', { rel: "icon", type: "image/png", sizes: "16x16", href: "/favicons/favicon-16x16.png"}],
    ['link', { rel: "manifest", href: "/favicons/site.webmanifest"}],
    ['link', { rel: "mask-icon", href: "/favicons/safari-pinned-tab.svg", color: "#000000"}],
    ['link', { rel: "shortcut icon", href: "/favicon.ico"}],
    ['meta', { name: "msapplication-TileColor", content: "#000000"}],
    ['meta', { name: "msapplication-config", content: "/favicons/browserconfig.xml"}],
    ['meta', { name: "theme-color", content: "#000000"}],
    ['script', { async: '', src: 'https://www.googletagmanager.com/gtag/js?id=G-QTDTMG01Z5' } ],
    ['script', {}, "window.dataLayer = window.dataLayer || [];\nfunction gtag(){dataLayer.push(arguments);}\ngtag('js', new Date());\ngtag('config', 'G-QTDTMG01Z5');"],
    ['script', {}, "(function(w,d,s,l,i){w[l]=w[l]||[];w[l].push({'gtm.start':new Date().getTime(),event:'gtm.js'});\nvar f=d.getElementsByTagName(s)[0],j=d.createElement(s),dl=l!='dataLayer'?'&l='+l:'';\nj.async=true;j.src='https://www.googletagmanager.com/gtm.js?id='+i+dl;f.parentNode.insertBefore(j,f);\n})(window,document,'script','dataLayer','GTM-TFFZXCQW');"],
    ['script', { async: '', defer: '', src: 'https://buttons.github.io/buttons.js' }]
  ],

  lastUpdated: true,

  // https://vitepress.dev/reference/default-theme-config
  themeConfig: {
    outline: [2, 4],

    logo: '/logo.svg',

    externalLinkIcon: true,

    search: {
      provider: 'local'
    },

    // Top navigation
    nav: [
      { text: 'Why Kit?', activeMatch: `^/#whykitops`, link: '/#whykitops' },
      { text: 'How does it work?', activeMatch: `^/#howdoesitwork`, link: '/#howdoesitwork' },
      { text: 'Docs', activeMatch: `^/docs`, link: '/docs/overview' },
      { text: 'Blog', activeMatch: `^/blog`, link: '/blog' },
    ],

    // Sidebar nav
    sidebar: [
      {
        text: 'Getting started',
        items: [
          { text: 'Overview', link: '/docs/overview' },
          { text: 'Quick Start', link: '/docs/quick-start.md' },
          { text: 'Next Steps', link: '/docs/next-steps.md' },
          { text: 'Use Cases', link: '/docs/use-cases' },
          { text: 'Why KitOps?', link: '/docs/why-kitops' },
        ]
      },
      {
        text: 'ModelKit',
        items: [
          { text: 'Introduction', link: '/docs/modelkit/intro' },
          { text: 'Specification', link: '/docs/modelkit/spec' },
          { text: 'Compatibility', link: '/docs/modelkit/compatibility' },
        ]
      },
      {
        text: 'Kitfile',
        items: [
          { text: 'Overview', link: '/docs/kitfile/kf-overview' },
          { text: 'Format', link: '/docs/kitfile/format' },
        ]
      },
      {
        text: 'CLI',
        items: getSidebarItemsFromMdFiles('docs/cli', {
            replacements: {
              'cli-reference': 'Command Reference' ,
              'installation': 'Download & Install'
            },
            textFormat: (text) => text.replaceAll('cli-', '')
          })
      },
      {
        text: 'Contribute',
        items: [
          { text: 'Contribute to KitOps docs', link: '/contributing' }
        ]
      },
    ],

    socialLinks: [
      {
        icon: 'discord',
        link: 'https://discord.gg/3eDb4yAN'
      },
    ],
    footer: {
      license: {
        text: 'MIT License',
        link: 'https://opensource.org/licenses/MIT'
      },
      copyright: `Copyright © ${new Date().getFullYear()} Jozu`
    }
  },

  vite: {
    resolve: {
      alias: [
        // Override the footer with out custom footer
        {
          find: /^.*\/VPFooter\.vue$/,
          replacement: fileURLToPath(
            new URL('./theme/components/Footer.vue', import.meta.url)
          )
        },
        {
          find: '@',
          replacement: resolve(__dirname, '../src'),
        },
        {
          find: '$public',
          replacement: resolve(__dirname, '../src/public')
        }
      ]
    }
  },

  ignoreDeadLinks: true
})
