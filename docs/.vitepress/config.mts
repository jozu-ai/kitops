import { URL, fileURLToPath } from 'node:url'

import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  appearance: 'force-dark',
  srcDir: 'src',
  title: 'KitOps',
  titleTemplate: 'KitOps',
  description: 'Documentation for KitOps',

  head: [
    // ['link', { rel: "apple-touch-icon", sizes: "180x180", href: "/favicons/apple-touch-icon.png"}],
    ['link', { rel: "icon", type: "image/png", sizes: "32x32", href: "/favicons/favicon-32x32.png"}],
    ['link', { rel: "icon", type: "image/png", sizes: "16x16", href: "/favicons/favicon-16x16.png"}],
    // ['link', { rel: "manifest", href: "/favicons/site.webmanifest"}],
    // ['link', { rel: "mask-icon", href: "/favicons/safari-pinned-tab.svg", color: "#3a0839"}],
    ['link', { rel: "shortcut icon", href: "/favicon.ico"}],
    // ['meta', { name: "msapplication-TileColor", content: "#3a0839"}],
    // ['meta', { name: "msapplication-config", content: "/favicons/browserconfig.xml"}],
    // ['meta', { name: "theme-color", content: "#ffffff"}],
  ],

  lastUpdated: true,

  // https://vitepress.dev/reference/default-theme-config
  themeConfig: {
    logo: '/logo.svg',

    externalLinkIcon: true,

    search: {
      provider: 'local'
    },

    // Top navigation
    nav: [
      { text: 'Why KitOps?', activeMatch: `^/#whykitops`, link: '/#whykitops' },
      { text: 'How does it work?', activeMatch: `^/#howdoesitwork`, link: '/#howdoesitwork' },
      { text: 'Docs', activeMatch: `^/docs`, link: '/docs/overview' },
      // { text: 'Blog', activeMatch: `^/blog`, link: '/blog' },
    ],

    // Sidebar nav
    sidebar: [
      {
        text: 'Getting started',
        items: [
          { text: 'Overview', link: '/docs/overview' },
          { text: 'Installation', link: '/docs/installation' },
        ]
      },
      {
        text: 'CLI',
        items: [
          { text: 'Download & Install', link: '/docs/cli/installation' },
          { text: 'Usage', link: '/docs/cli/usage' },
          { text: 'Flags', link: '/docs/cli/flags' },
        ]
      },
      {
        text: 'Manifest',
        items: [
          { text: 'Overview', link: '/docs/manifest/overview' },
          { text: 'Structure', link: '/docs/manifest/overview' },
          { text: 'Building and running', link: '/docs/manifest/building-running' },
          { text: 'Creating a new model', link: '/' },
          { text: 'Training a model', link: '/' },
          { text: 'Strategies', link: '/' },
        ]
      },
      {
        text: 'MLOps with Kitfile',
        items: [
          { text: 'Continuos integration and deployment', link: '/docs/mlops/ci-cd' },
          { text: 'Monitoring and logging', link: '/docs/mlops/ci-cd' },
          { text: 'Orchestration', link: '/docs/mlops/ci-cd' },
          { text: 'Scalability and resources', link: '/docs/mlops/ci-cd' },
        ]
      },
      {
        text: 'Advanced',
        items: [
          { text: 'Fine-tuning', link: '/docs/mlops/ci-cd' },
          { text: 'Multi-models and Multi-tasking', link: '/docs/mlops/ci-cd' },
          { text: 'Parallel training', link: '/docs/mlops/ci-cd' },
          { text: 'Reinforcement and deep RL', link: '/docs/mlops/ci-cd' },
        ]
      },
      {
        text: 'Contribute',
        items: [
          { text: 'Contribute to KitOps docs', link: '/' }
        ]
      },
      {
        text: 'Documentation Examples',
        items: [
          { text: 'Markdown Examples', link: '/markdown-examples' },
          { text: 'Runtime API Examples', link: '/api-examples' }
        ]
      }
    ],

    socialLinks: [
      {
        icon: 'discord',
        link: 'https://discord.gg/XzSmtPn3'
      },
      {
        icon: 'github',
        link: 'https://github.com/jozu-ai/kitops'
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
        }
      ]
    }
  }
})
