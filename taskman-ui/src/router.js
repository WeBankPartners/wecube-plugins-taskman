import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)
export default new Router({
  routes: [
    {
      path: '/taskman/template-mgmt',
      name: '/taskman/template-mgmt',
      component: () => import('@/pages/template')
    },
    {
      path: '/taskman/template-group',
      name: '/taskman/template-group',
      component: () => import('@/pages/template-group')
    },
    {
      path: '/templateManagementIndex',
      name: 'templateManagementIndex',
      component: () => import('@/pages/temp-management/index')
    },
    {
      path: '/taskman/request-mgmt',
      name: '/taskman/request-mgmt',
      component: () => import('@/pages/request')
    },
    {
      path: '/taskman/task-mgmt',
      name: '/taskman/task-mgmt',
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
