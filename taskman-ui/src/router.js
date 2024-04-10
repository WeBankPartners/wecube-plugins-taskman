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
          component: () => import('@/pages/template'),
          meta: { en: 'Template Management', zh: '模板管理' }
        },
        {
          path: '/taskman/template-group',
          name: '/taskman/template-group',
          component: () => import('@/pages/template-group'),
          meta: { en: 'Template Grou', zh: '模板组' }
        },
        {
          path: '/templateManagementIndex',
          name: 'templateManagementIndex',
          component: () => import('@/pages/temp-management/index'),
          meta: { en: 'Template Management', zh: '模板管理' }
        },
        {
          path: '/taskman/requestAudit',
          name: '/taskman/requestAudit',
          component: () => import('@/pages/workbench/request-audit.vue'),
          meta: { en: 'Request Audit', zh: '请求审计' }
        },
        // 工作台
        {
          path: '/taskman/workbench',
          name: '/taskman/workbench',
          component: () => import('@/pages/workbench.vue'),
          redirect: '/taskman/workbench/dashboard',
          meta: { en: 'Dashboard', zh: '工作台' },
          children: [
            {
              path: '/taskman/workbench/dashboard',
              name: '/taskman/workbench/dashboard',
              component: () => import('@/pages/workbench/index.vue')
            },
            {
              path: '/taskman/workbench/template',
              name: '/taskman/workbench/template',
              component: () => import('@/pages/workbench/template/index'),
              meta: { en: 'Template selection', zh: '模板选择' }
            },
            {
              path: '/taskman/workbench/createPublish',
              name: '/taskman/workbench/createPublish',
              component: () => import('@/pages/workbench/publish/create'),
              meta: { en: 'New Publish', zh: '新建发布' }
            },
            {
              path: '/taskman/workbench/detailPublish',
              name: '/taskman/workbench/detailPublish',
              component: () => import('@/pages/workbench/publish/detail'),
              meta: { en: 'Publish Detail', zh: '发布详情' }
            },
            {
              path: '/taskman/workbench/publishHistory',
              name: '/taskman/workbench/publishHistory',
              component: () => import('@/pages/workbench/publish/list'),
              meta: { en: 'History(Group)', zh: '历史(本组)' }
            },
            {
              path: '/taskman/workbench/createRequest',
              name: '/taskman/workbench/createRequest',
              component: () => import('@/pages/workbench/request/create'),
              meta: { en: 'New Request', zh: '新建请求' }
            },
            {
              path: '/taskman/workbench/detailRequest',
              name: '/taskman/workbench/detailRequest',
              component: () => import('@/pages/workbench/request/detail'),
              meta: { en: 'Request Detail', zh: '请求详情' }
            },
            {
              path: '/taskman/workbench/requestHistory',
              name: '/taskman/workbench/requestHistory',
              component: () => import('@/pages/workbench/request/list'),
              meta: { en: 'History(Group)', zh: '历史(本组)' }
            },
            {
              path: '/taskman/workbench/createProblem',
              name: '/taskman/workbench/createProblem',
              component: () => import('@/pages/workbench/problem/create'),
              meta: { en: 'New Problem', zh: '新建问题' }
            },
            {
              path: '/taskman/workbench/detailProblem',
              name: '/taskman/workbench/detailProblem',
              component: () => import('@/pages/workbench/problem/detail'),
              meta: { en: 'Problem Detail', zh: '问题详情' }
            },
            {
              path: '/taskman/workbench/problemHistory',
              name: '/taskman/workbench/problemHistory',
              component: () => import('@/pages/workbench/problem/list'),
              meta: { en: 'History(Group)', zh: '历史(本组)' }
            },
            {
              path: '/taskman/workbench/createEvent',
              name: '/taskman/workbench/createEvent',
              component: () => import('@/pages/workbench/event/create'),
              meta: { en: 'New Event', zh: '新建事件' }
            },
            {
              path: '/taskman/workbench/detailEvent',
              name: '/taskman/workbench/detailEvent',
              component: () => import('@/pages/workbench/event/detail'),
              meta: { en: 'Event Detail', zh: '事件详情' }
            },
            {
              path: '/taskman/workbench/eventHistory',
              name: '/taskman/workbench/eventHistory',
              component: () => import('@/pages/workbench/event/list'),
              meta: { en: 'History(Group)', zh: '历史(本组)' }
            },
            {
              path: '/taskman/workbench/createChange',
              name: '/taskman/workbench/createChange',
              component: () => import('@/pages/workbench/change/create'),
              meta: { en: 'New Change', zh: '新建变更' }
            },
            {
              path: '/taskman/workbench/detailChange',
              name: '/taskman/workbench/detailChange',
              component: () => import('@/pages/workbench/change/detail'),
              meta: { en: 'Change Detail', zh: '变更详情' }
            },
            {
              path: '/taskman/workbench/changeHistory',
              name: '/taskman/workbench/changeHistory',
              component: () => import('@/pages/workbench/change/list'),
              meta: { en: 'History(Group)', zh: '历史(本组)' }
            }
          ]
        }
      ]
    }
  ]
})
