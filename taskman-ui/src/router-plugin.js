import template from '@/pages/template'
import templateGroup from '@/pages/template-group'
import templateManagementIndex from '@/pages/temp-management/index'
import request from '@/pages/request'
import requestManagementIndex from '@/pages/request-management/index'
import task from '@/pages/task'
import taskManagementIndex from '@/pages/task-mgmt/index'
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
    path: '/taskManagementIndex',
    name: 'taskManagementIndex',
    component: taskManagementIndex
  }
]
export default router
