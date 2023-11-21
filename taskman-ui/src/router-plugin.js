import template from '@/pages/template'
import templateGroup from '@/pages/template-group'
import templateManagementIndex from '@/pages/temp-management/index'
import request from '@/pages/request'
import requestManagementIndex from '@/pages/request-management/index'
import task from '@/pages/task'
import taskManagementIndex from '@/pages/task-mgmt/index'
import requestCheck from '@/pages/request-management/request-check'
const router = [
  {
    path: '/taskman/template-mgmt',
    name: '/taskman/template-mgmt',
    component: template
  },
  {
    path: '/taskman/template-group',
    name: '/taskman/template-group',
    component: templateGroup
  },
  {
    path: '/templateManagementIndex',
    name: 'templateManagementIndex',
    component: templateManagementIndex
  },
  {
    path: '/taskman/request-mgmt',
    name: '/taskman/request-mgmt',
    component: request
  },
  {
    path: '/requestManagementIndex',
    name: 'requestManagementIndex',
    component: requestManagementIndex
  },
  {
    path: '/taskman/task-mgmt',
    name: '/taskman/task-mgmt',
    component: task
  },
  {
    path: '/taskMgmtIndex',
    name: 'taskMgmtIndex',
    component: taskManagementIndex
  },
  {
    path: '/requestCheck',
    name: 'requestCheck',
    component: requestCheck
  },
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
      }
    ]
  }
]
export default router
