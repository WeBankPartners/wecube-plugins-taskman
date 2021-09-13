import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)
export default new Router({
  routes: [
    {
      path: '/template',
      name: 'template',
      component: () => import('@/pages/template')
    },
    {
      path: '/',
      name: 'templateGroup',
      component: () => import('@/pages/template-group')
    }
  ]
})
