import type MarkdownIt from 'markdown-it'
import type StateBlock from 'markdown-it/lib/rules_block/state_block'

import { URL, fileURLToPath } from 'node:url'
import { resolve } from 'path'

import { defineConfig, type HeadConfig } from 'vitepress'
import { getSidebarItemsFromMdFiles } from './utils.mts'

const inProd = process.env.NODE_ENV === 'production'

const head = [
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
  ['script', { async: '', src: 'https://www.googletagmanager.com/gtag/js?id=G-QTDTMG01Z5'}],
  ['script', {}, "window.dataLayer = window.dataLayer || [];\nfunction gtag(){dataLayer.push(arguments);}\ngtag('js', new Date());\ngtag('config', 'G-QTDTMG01Z5');"],
  ['script', {}, "(function(w,d,s,l,i){w[l]=w[l]||[];w[l].push({'gtm.start':new Date().getTime(),event:'gtm.js'});\nvar f=d.getElementsByTagName(s)[0],j=d.createElement(s),dl=l!='dataLayer'?'&l='+l:'';\nj.async=true;j.src='https://www.googletagmanager.com/gtm.js?id='+i+dl;f.parentNode.insertBefore(j,f);\n})(window,document,'script','dataLayer','GTM-TFFZXCQW');"],
  ['script', { async: '', defer: '', src: 'https://buttons.github.io/buttons.js' }],
] as HeadConfig[]

// Prod only scripts
if (inProd) {
//   // Add REB2B
//   head.push(
//     ['script', {}, '!function () {var reb2b = window.reb2b = window.reb2b || [];if (reb2b.invoked) return;reb2b.invoked = true;reb2b.methods = ["identify", "collect"];reb2b.factory = function (method) {return function () {var args = Array.prototype.slice.call(arguments);args.unshift(method);reb2b.push(args);return reb2b;};};for (var i = 0; i < reb2b.methods.length; i++) {var key = reb2b.methods[i];reb2b[key] = reb2b.factory(key);}reb2b.load = function (key) {var script = document.createElement("script");script.type = "text/javascript";script.async = true;script.src = "https://s3-us-west-2.amazonaws.com/b2bjsstore/b/" + key + "/reb2b.js.gz";var first = document.getElementsByTagName("script")[0];first.parentNode.insertBefore(script, first);};reb2b.SNIPPET_VERSION = "1.0.1";reb2b.load("L9NMMZHE4PNW");}();'],
//   )

  // Posthog
  head.push(
    ['script', {}, '!function(t,e){var o,n,p,r;e.__SV||(window.posthog=e,e._i=[],e.init=function(i,s,a){function g(t,e){var o=e.split(".");2==o.length&&(t=t[o[0]],e=o[1]),t[e]=function(){t.push([e].concat(Array.prototype.slice.call(arguments,0)))}}(p=t.createElement("script")).type="text/javascript",p.crossOrigin="anonymous",p.async=!0,p.src=s.api_host.replace(".i.posthog.com","-assets.i.posthog.com")+"/static/array.js",(r=t.getElementsByTagName("script")[0]).parentNode.insertBefore(p,r);var u=e;for(void 0!==a?u=e[a]=[]:a="posthog",u.people=u.people||[],u.toString=function(t){var e="posthog";return"posthog"!==a&&(e+="."+a),t||(e+=" (stub)"),e},u.people.toString=function(){return u.toString(1)+".people (stub)"},o="init capture register register_once register_for_session unregister unregister_for_session getFeatureFlag getFeatureFlagPayload isFeatureEnabled reloadFeatureFlags updateEarlyAccessFeatureEnrollment getEarlyAccessFeatures on onFeatureFlags onSessionId getSurveys getActiveMatchingSurveys renderSurvey canRenderSurvey getNextSurveyStep identify setPersonProperties group resetGroups setPersonPropertiesForFlags resetPersonPropertiesForFlags setGroupPropertiesForFlags resetGroupPropertiesForFlags reset get_distinct_id getGroups get_session_id get_session_replay_url alias set_config startSessionRecording stopSessionRecording sessionRecordingStarted captureException loadToolbar get_property getSessionProperty createPersonProfile opt_in_capturing opt_out_capturing has_opted_in_capturing has_opted_out_capturing clear_opt_in_out_capturing debug".split(" "),n=0;n<o.length;n++)g(u,o[n]);e._i.push([i,s,a])},e.__SV=1)}(document,window.posthog||[]); posthog.init("phc_1BXud8NPWtXIaPkE5NhoWUmW5BTXjuRdDgjeCzoPlYU",{api_host:"https://us.i.posthog.com", person_profiles: "identified_only" })']
  )

  // CookieYes
  head.push(
    ['script', { id: 'cookieyes', type: 'text/javascript', src: 'https://cdn-cookieyes.com/client_data/43ed933a5c31287e10124caf/script.js' }]
  )
}

// https://vitepress.dev/reference/site-config
export default defineConfig({
  appearance: 'force-dark',
  srcDir: 'src',
  title: 'KitOps',
  description: 'Documentation for KitOps',

  head,

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
      { text: 'Get Started', activeMatch: '^/#getstarted', link: '/docs/get-started.html' },
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
          { text: 'Get Started', link: '/docs/get-started' },
          { text: 'Next Steps', link: '/docs/next-steps' },
          { text: 'Deploy ModelKits', link: '/docs/deploy' },
          { text: 'Local LLM Dev', link: '/docs/dev-mode' },
          { text: 'Why KitOps?', link: '/docs/why-kitops' },
          { text: 'How it is Used', link: '/docs/use-cases' },
          { text: 'KitOps versus...', link: '/docs/versus' },
        ]
      },
      {
        text: 'ModelKit',
        items: [
          { text: 'Overview', link: '/docs/modelkit/intro' },
          { text: 'Specification', link: '/docs/modelkit/spec' },
          { text: 'ModelKit Quick Starts', link: 'https://jozu.ml/organization/jozu-quickstarts' },
          { text: 'Compatibility', link: '/docs/modelkit/compatibility' },
        ]
      },
      {
        text: 'Kitfile',
        items: [
          { text: 'Overview', link: '/docs/kitfile/kf-overview' },
          { text: 'Format', link: '/docs/kitfile/format' }
        ]
      },
      {
        text: 'Python Library',
        items: [
          { text: 'Overview', link: '/docs/pykitops.html' },
          { text: 'How-to Guides', link: '/docs/pykitops/how-to-guides.html' },
          { text: 'Reference', link: '/docs/pykitops/reference.html' },
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
        link: 'https://discord.gg/Tapeh8agYy'
      },
    ],

    footer: {
      message: 'Made with <3 by Jozu',
      copyright: `Copyright © ${new Date().getFullYear()} Jozu`
    }
  },

  transformPageData(pageData, { siteConfig }) {
    // Generate the canonical url's on each page, considering the cleanUrl config
    const canonicalUrl = `https://kitops.ml/${pageData.relativePath}`
      .replace(/index\.md$/, '')
      .replace(/\.md$/, siteConfig.cleanUrls ? '' : '.html')

    pageData.frontmatter.head ??= []
    pageData.frontmatter.head.push([
      'link',
      { rel: 'canonical', href: canonicalUrl }
    ])
  },

  sitemap: {
    hostname: 'https://kitops.ml'
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

  ignoreDeadLinks: true,

  cleanUrls: true,

  markdown: {
    config: (md) => {
      md.use(discordBannerPlugin);
    },
  }
})

// Custom Markdown-it plugin
function discordBannerPlugin(md: MarkdownIt) {
  const marker = '[ discord banner ]'

  md.block.ruler.before(
    'fence',
    'discord-banner',
    (state: StateBlock, startLine: number, endLine: number, silent: boolean): boolean => {
    const start = state.bMarks[startLine] + state.tShift[startLine];
    const max = state.eMarks[startLine];

    // Match the custom marker
    if (state.src.slice(start, max).trim() !== marker) {
      return false;
    }

    if (silent) return true;

    // Create a token for the banner component
    const token = state.push('discord-banner', 'div', 0);
    token.map = [startLine, startLine + 1];
    state.line = startLine + 1;

    return true;
  });

  // Render the component
  md.renderer.rules['discord-banner'] = () => '<DiscordBanner />'
}

