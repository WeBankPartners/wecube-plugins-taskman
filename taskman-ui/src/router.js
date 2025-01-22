import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)
export const routes = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/login')
  },
  {
    path: '/',
    name: '/workbench',
    redirect: '/workbench',
    component: () => import('@/pages/index'),
    children: [
      {
        path: '/template-mgmt',
        name: '/template-mgmt',
        component: () => import('@/pages/template'),
        meta: { en: 'Template Management', zh: '模板管理' }
      },
      {
        path: '/template-group',
        name: '/template-group',
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
        path: '/requestAudit',
        name: '/requestAudit',
        component: () => import('@/pages/workbench/request-audit.vue'),
        meta: { en: 'Request Report', zh: '请求报表' }
      },
      // 工作台
      {
        path: '/workbench',
        name: '/workbench',
        component: () => import('@/pages/workbench.vue'),
        redirect: '/workbench/dashboard',
        meta: { en: 'Dashboard', zh: '工作台' },
        children: [
          {
            path: '/workbench/dashboard',
            name: '/workbench/dashboard',
            component: () => import('@/pages/workbench/index.vue')
          },
          {
            path: '/workbench/template',
            name: '/workbench/template',
            component: () => import('@/pages/workbench/template/index'),
            meta: { en: 'Template selection', zh: '模板选择' }
          },
          {
            path: '/workbench/createPublish',
            name: '/workbench/createPublish',
            component: () => import('@/pages/workbench/publish/create'),
            meta: { en: 'New Publish', zh: '新建发布' }
          },
          {
            path: '/workbench/detailPublish',
            name: '/workbench/detailPublish',
            component: () => import('@/pages/workbench/publish/detail'),
            meta: { en: 'Publish Detail', zh: '发布详情' }
          },
          {
            path: '/workbench/publishHistory',
            name: '/workbench/publishHistory',
            component: () => import('@/pages/workbench/publish/list'),
            meta: { en: 'History(Group)', zh: '历史(本组)' }
          },
          {
            path: '/workbench/createRequest',
            name: '/workbench/createRequest',
            component: () => import('@/pages/workbench/request/create'),
            meta: { en: 'New Request', zh: '新建请求' }
          },
          {
            path: '/workbench/detailRequest',
            name: '/workbench/detailRequest',
            component: () => import('@/pages/workbench/request/detail'),
            meta: { en: 'Request Detail', zh: '请求详情' }
          },
          {
            path: '/workbench/requestHistory',
            name: '/workbench/requestHistory',
            component: () => import('@/pages/workbench/request/list'),
            meta: { en: 'History(Group)', zh: '历史(本组)' }
          },
          {
            path: '/workbench/createProblem',
            name: '/workbench/createProblem',
            component: () => import('@/pages/workbench/problem/create'),
            meta: { en: 'New Problem', zh: '新建问题' }
          },
          {
            path: '/workbench/detailProblem',
            name: '/workbench/detailProblem',
            component: () => import('@/pages/workbench/problem/detail'),
            meta: { en: 'Problem Detail', zh: '问题详情' }
          },
          {
            path: '/workbench/problemHistory',
            name: '/workbench/problemHistory',
            component: () => import('@/pages/workbench/problem/list'),
            meta: { en: 'History(Group)', zh: '历史(本组)' }
          },
          {
            path: '/workbench/createEvent',
            name: '/workbench/createEvent',
            component: () => import('@/pages/workbench/event/create'),
            meta: { en: 'New Event', zh: '新建事件' }
          },
          {
            path: '/workbench/detailEvent',
            name: '/workbench/detailEvent',
            component: () => import('@/pages/workbench/event/detail'),
            meta: { en: 'Event Detail', zh: '事件详情' }
          },
          {
            path: '/workbench/eventHistory',
            name: '/workbench/eventHistory',
            component: () => import('@/pages/workbench/event/list'),
            meta: { en: 'History(Group)', zh: '历史(本组)' }
          },
          {
            path: '/workbench/createChange',
            name: '/workbench/createChange',
            component: () => import('@/pages/workbench/change/create'),
            meta: { en: 'New Change', zh: '新建变更' }
          },
          {
            path: '/workbench/detailChange',
            name: '/workbench/detailChange',
            component: () => import('@/pages/workbench/change/detail'),
            meta: { en: 'Change Detail', zh: '变更详情' }
          },
          {
            path: '/workbench/changeHistory',
            name: '/workbench/changeHistory',
            component: () => import('@/pages/workbench/change/list'),
            meta: { en: 'History(Group)', zh: '历史(本组)' }
          }
        ]
      }
    ]
  }
]

export default routes
