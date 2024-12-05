// https://vitepress.dev/guide/custom-theme
import { h, ref } from 'vue'
import { type Theme, inBrowser } from 'vitepress'
import DefaultTheme from 'vitepress/theme'
import PlatformSelect from './components/PlatformSelect.vue'
import PlatformSnippet from './components/PlatformSnippet.vue'
import GithubStartButton from './components/GithubStartButton.vue'
import './assets/css/fonts.css'
import './assets/css/tailwind.css'
import './style.css'
import DiscordModal from './components/DiscordModal.vue'

const isPlatformModalOpen = ref(false)

export default {
  extends: DefaultTheme,
  Layout: () => {
    if (inBrowser) {
      // Google tag manager
      const noscript = document.createElement('noscript')
      const ifr = document.createElement('iframe')
      ifr.src = 'https://www.googletagmanager.com/ns.html?id=GTM-TFFZXCQW';
      ifr.height = 0;
      ifr.width = 0;
      ifr.style = 'display:none;visibility:hidden;';
      noscript.appendChild(ifr)
      document.body.insertBefore(noscript, document.body.firstElementChild)
    }

    return h(DefaultTheme.Layout, null, {
      // https://vitepress.dev/guide/extending-default-theme#layout-slots
      'sidebar-nav-after': () => h(PlatformSelect),
      'nav-bar-content-after': () => h(GithubStartButton, {
        class: 'ml-4 pt-2'
      }),
      'doc-bottom': () => h(DiscordModal)
    })
  },
  enhanceApp({ app, router, siteData }) {
    app.component('PlatformSnippet', PlatformSnippet)
    app.provide('isPlatformModalOpen', isPlatformModalOpen)
  }
} satisfies Theme
