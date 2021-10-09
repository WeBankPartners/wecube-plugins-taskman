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
      path: '/templateGroup',
      name: 'templateGroup',
      component: () => import('@/pages/template-group')
    },
    {
      path: '/templateManagementIndex',
      name: 'templateManagementIndex',
      component: () => import('@/pages/temp-management/index')
    },
    {
      path: '/',
      name: 'request',
      component: () => import('@/pages/request')
    },
    {
      path: '/requestManagementIndex',
      name: 'requestManagementIndex',
      component: () => import('@/pages/request-management/index')
    }
  ]
})
