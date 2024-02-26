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

    // lastUpdated: {
    //   text: 'Updated at',
    //   formatOptions: {
    //     dateStyle: 'full',
    //     timeStyle: 'medium'
    //   }
    // },

    // Top navigation
    nav: [
      { text: 'Docs', activeMatch: `^/docs`, link: '/docs/overview' },
      { text: 'Guides', activeMatch: `^/guides`, link: '/guides/index' },
      { text: 'Feed', activeMatch: `^/feed`, link: '/feed' },
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
        icon: 'github',
        link: 'https://github.com/jozu-ai/kitops'
      },
      {
        icon: {
          svg: '<svg  role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 56 56" fill="none"><path d="M5.25 6.18281L27.8658 0L50.4817 6.18281V56C50.4817 56 38.4113 49.8172 27.8658 49.8172C17.3204 49.8172 5.25 56 5.25 56V6.18281Z" fill="#075550"></path><path d="M17.7325 25.3151C16.0439 25.3151 14.701 24.9033 13.7039 24.0796C12.7068 23.256 12.0393 22.1577 11.7016 20.7849L15.1995 19.501C15.296 19.9532 15.4408 20.3973 15.6338 20.8334C15.8428 21.2533 16.1162 21.6005 16.454 21.8751C16.8078 22.1496 17.2339 22.2869 17.7325 22.2869C18.5366 22.2869 19.1478 22.0366 19.5659 21.5359C20.0001 21.0191 20.2172 20.1147 20.2172 18.8226V8.84161H24.0288V18.8226C24.0288 20.8414 23.49 22.4323 22.4125 23.5951C21.351 24.7418 19.791 25.3151 17.7325 25.3151Z" fill="#3febe0"></path><path d="M35.6325 25.3346C33.9921 25.3346 32.5447 24.9873 31.2903 24.2928C30.0359 23.5984 29.0548 22.6213 28.3472 21.3615C27.6396 20.1018 27.2858 18.6321 27.2858 16.9524C27.2858 15.2728 27.6396 13.8031 28.3472 12.5433C29.0548 11.2836 30.0359 10.3065 31.2903 9.61202C32.5447 8.91755 33.9921 8.57031 35.6325 8.57031C37.289 8.57031 38.7445 8.91755 39.9989 9.61202C41.2533 10.3065 42.2344 11.2836 42.942 12.5433C43.6657 13.8031 44.0275 15.2728 44.0275 16.9524C44.0275 18.6321 43.6657 20.1018 42.942 21.3615C42.2344 22.6213 41.2533 23.5984 39.9989 24.2928C38.7445 24.9873 37.289 25.3346 35.6325 25.3346ZM35.6325 22.1852C36.4206 22.1852 37.1443 21.9914 37.8037 21.6038C38.463 21.2 38.9857 20.6105 39.3717 19.8353C39.7738 19.0601 39.9748 18.0991 39.9748 16.9524C39.9748 15.8057 39.7738 14.8448 39.3717 14.0696C38.9857 13.2943 38.463 12.7129 37.8037 12.3253C37.1443 11.9215 36.4206 11.7197 35.6325 11.7197C34.8606 11.7197 34.1449 11.9215 33.4855 12.3253C32.8262 12.7129 32.2954 13.2943 31.8934 14.0696C31.5074 14.8448 31.3144 15.8057 31.3144 16.9524C31.3144 18.0991 31.5074 19.0601 31.8934 19.8353C32.2954 20.6105 32.8262 21.2 33.4855 21.6038C34.1449 21.9914 34.8606 22.1852 35.6325 22.1852Z" fill="#3febe0"></path><path d="M12.2334 44.8154V41.8841L20.363 32H12.5952V28.6326H25.3083V31.5639L17.2511 41.4481H25.6702V44.8154H12.2334Z" fill="#3febe0"></path><path d="M35.6017 45.0735C33.0768 45.0735 31.2836 44.4679 30.2221 43.2566C29.1768 42.0453 28.6541 40.1961 28.6541 37.7089V28.6H32.4656V38.0238C32.4656 38.8152 32.538 39.5097 32.6827 40.1072C32.8436 40.6887 33.1491 41.1409 33.5994 41.4639C34.0497 41.7869 34.7172 41.9484 35.6017 41.9484C36.4862 41.9484 37.1456 41.7869 37.5798 41.4639C38.0301 41.1409 38.3277 40.6887 38.4724 40.1072C38.6332 39.5097 38.7136 38.8152 38.7136 38.0238V28.6H42.5252V37.7089C42.5252 40.1961 41.9944 42.0453 40.933 43.2566C39.8716 44.4679 38.0945 45.0735 35.6017 45.0735Z" fill="#3febe0"></path></svg>'
        },
        link: 'https://jozu.com',
        ariaLabel: 'JOzu Website'
      }
    ],

    footer: {
      license: {
        text: 'MIT License',
        link: 'https://opensource.org/licenses/MIT'
      },
      copyright: `Copyright © ${new Date().getFullYear()} Jozu`
    }
  }
})
