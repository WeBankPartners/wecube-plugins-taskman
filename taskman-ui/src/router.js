import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)
export default new Router({
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/login')
    },
    {
      path: '/',
      name: '/taskman/workbench',
      redirect: '/taskman/workbench',
      component: () => import('@/pages/index'),
      children: [
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
        // 工作台路由
        {
          path: '/taskman/workbench',
          name: '/taskman/workbench',
          component: () => import('@/pages/workbench.vue'),
          redirect: '/taskman/workbench/dashboard',
          children: [
            // 工作台首页
            {
              path: '/taskman/workbench/dashboard',
              name: '/taskman/workbench/dashboard',
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
            // 发布详情
            {
              path: '/taskman/workbench/detailPublish',
              name: '/taskman/workbench/detailPublish',
              component: () => import('@/pages/workbench/publish/detail')
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
            // 请求详情
            {
              path: '/taskman/workbench/detailRequest',
              name: '/taskman/workbench/detailRequest',
              component: () => import('@/pages/workbench/request/detail')
            },
            // 请求历史
            {
              path: '/taskman/workbench/requestHistory',
              name: '/taskman/workbench/requestHistory',
              component: () => import('@/pages/workbench/request/list')
            },
            // 新建问题
            {
              path: '/taskman/workbench/createProblem',
              name: '/taskman/workbench/createProblem',
              component: () => import('@/pages/workbench/problem/create')
            },
            // 问题详情
            {
              path: '/taskman/workbench/detailProblem',
              name: '/taskman/workbench/detailProblem',
              component: () => import('@/pages/workbench/problem/detail')
            },
            // 问题历史
            {
              path: '/taskman/workbench/problemHistory',
              name: '/taskman/workbench/problemHistory',
              component: () => import('@/pages/workbench/problem/list')
            },
            // 新建事件
            {
              path: '/taskman/workbench/createEvent',
              name: '/taskman/workbench/createEvent',
              component: () => import('@/pages/workbench/event/create')
            },
            // 事件详情
            {
              path: '/taskman/workbench/detailEvent',
              name: '/taskman/workbench/detailEvent',
              component: () => import('@/pages/workbench/event/detail')
            },
            // 事件历史
            {
              path: '/taskman/workbench/eventHistory',
              name: '/taskman/workbench/eventHistory',
              component: () => import('@/pages/workbench/event/list')
            },
            // 新建变更
            {
              path: '/taskman/workbench/createChange',
              name: '/taskman/workbench/createChange',
              component: () => import('@/pages/workbench/change/create')
            },
            // 变更详情
            {
              path: '/taskman/workbench/detailChange',
              name: '/taskman/workbench/detailChange',
              component: () => import('@/pages/workbench/change/detail')
            },
            // 变更历史
            {
              path: '/taskman/workbench/changeHistory',
              name: '/taskman/workbench/changeHistory',
              component: () => import('@/pages/workbench/change/list')
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
    }
  ]
})
