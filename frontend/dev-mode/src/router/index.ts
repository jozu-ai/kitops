import { createRouter, createWebHistory } from 'vue-router'

import NotFound from '../views/NotFound.vue'
import PromptsView from '../views/PromptsView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: PromptsView,
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: NotFound
    },
  ]
})

export default router
