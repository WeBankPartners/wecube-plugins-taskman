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
      path: '/requestCheck',
      name: 'requestCheck',
      component: () => import('@/pages/request-management/request-check')
    },
    {
      path: '/requestManagementIndex',
      name: 'requestManagementIndex',
      component: () => import('@/pages/request-management/index')
    },
    // 工作台子页面路由需要以/taskman/workbench为前缀，有些判断条件是以/taskman/workbench写死判断的
    {
      path: '/taskman',
      name: 'taskman',
      component: () => import('@/pages/index'),
      redirect: '/taskman/workbench',
      children: [
        // 个人工作台
        {
          path: '/taskman/workbench',
          name: '/taskman/workbench',
          component: () => import('@/pages/workbench/index.vue')
        },
        // 模板选择
        {
          path: '/taskman/workbench/template',
          name: '/taskman/workbench/template',
          component: () => import('@/pages/workbench/template/index')
        },
        // 新建发布
        {
          path: '/taskman/workbench/createPublish',
          name: '/taskman/workbench/createPublish',
          component: () => import('@/pages/workbench/publish/create')
        },
        // 发布历史
        {
          path: '/taskman/workbench/publishHistory',
          name: '/taskman/workbench/publishHistory',
          component: () => import('@/pages/workbench/publish/list')
        },
        // 新建请求
        {
          path: '/taskman/workbench/createRequest',
          name: '/taskman/workbench/createRequest',
          component: () => import('@/pages/workbench/request/create')
        },
        // 请求历史
        {
          path: '/taskman/workbench/requestHistory',
          name: '/taskman/workbench/requestHistory',
          component: () => import('@/pages/workbench/request/list')
        },
        // 请求审计
        {
          path: '/taskman/requestAudit',
          name: '/taskman/requestAudit',
          component: () => import('@/pages/workbench/request-audit.vue')
        }
      ]
    }
  ]
})
