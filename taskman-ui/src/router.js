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
      path: '/request',
      name: 'request',
      component: () => import('@/pages/request')
    },
    {
      path: '/task',
      name: 'task',
      component: () => import('@/pages/task')
    },
    {
      path: '/taskMgmtIndex',
      name: 'taskMgmtIndex',
      component: () => import('@/pages/task-mgmt/index')
    },
    {
      path: '/requestManagementIndex',
      name: 'requestManagementIndex',
      component: () => import('@/pages/request-management/index')
    }
  ]
})
