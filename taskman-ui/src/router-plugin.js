import pluginEntry from '@/plugin-entry.vue'
import template from '@/pages/template.vue'
import templateGroup from '@/pages/template-group.vue'
import templateManagementIndex from '@/pages/temp-management/index.vue'
import request from '@/pages/request.vue'
import requestManagementIndex from '@/pages/request-management/index.vue'
import task from '@/pages/task.vue'
import taskManagementIndex from '@/pages/task-mgmt/index.vue'
import requestCheck from '@/pages/request-management/request-check.vue'
import workbenchIndex from '@/pages/workbench/index.vue'
import workbenchTemplate  from '@/pages/workbench/template/index.vue'
import workbenchPublishCreate from '@/pages/workbench/publish/create'
import workbenchPublishDetail from '@/pages/workbench/publish/detail'
import workbenchPublishList  from '@/pages/workbench/publish/list'
import workbenchRequestCreate from '@/pages/workbench/request/create'
import workbenchRequestDetail from '@/pages/workbench/request/detail'
import workbenchRequestList from '@/pages/workbench/request/list'
import workbenchProblemCreate from '@/pages/workbench/problem/create'
import workbenchProblemDetail from '@/pages/workbench/problem/detail'
import workbenchProblemList from '@/pages/workbench/problem/list'
import workbenchEventCreate from '@/pages/workbench/event/create'
import workbenchEventDetail from '@/pages/workbench/event/detail'
import workbenchEventList from '@/pages/workbench/event/list'
import workbenchChangeCreate from '@/pages/workbench/change/create'
import workbenchChangeDetail from '@/pages/workbench/change/detail'
import workbenchChangeList from '@/pages/workbench/change/list'
import workbenchRequestAudit from '@/pages/workbench/request-audit.vue'
import workbenchVue  from '@/pages/workbench.vue'

const router = [
  {
    path: '/taskman',
    name: 'taskman',
    redirect: '/taskman/workbench',
    component: pluginEntry,
    children:[
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
        path: '/taskman/workbench',
        name: 'taskman/workbench',
        component: workbenchVue,
        redirect: '/taskman/workbench/dashboard',
        children: [
          // 个人工作台
          {
            path: '/taskman/workbench/dashboard',
            name: '/taskman/workbench/dashboard',
            component: workbenchIndex
          },
          // 模板选择
          {
            path: '/taskman/workbench/template',
            name: '/taskman/workbench/template',
            component: workbenchTemplate
          },
          // 新建发布
          {
            path: '/taskman/workbench/createPublish',
            name: '/taskman/workbench/createPublish',
            component: workbenchPublishCreate
          },
          // 发布详情
          {
            path: '/taskman/workbench/detailPublish',
            name: '/taskman/workbench/detailPublish',
            component: workbenchPublishDetail
          },
          // 发布历史
          {
            path: '/taskman/workbench/publishHistory',
            name: '/taskman/workbench/publishHistory',
            component: workbenchPublishList
          },
          // 新建请求
          {
            path: '/taskman/workbench/createRequest',
            name: '/taskman/workbench/createRequest',
            component: workbenchRequestCreate
          },
          // 请求详情
          {
            path: '/taskman/workbench/detailRequest',
            name: '/taskman/workbench/detailRequest',
            component: workbenchRequestDetail
          },
          // 请求历史
          {
            path: '/taskman/workbench/requestHistory',
            name: '/taskman/workbench/requestHistory',
            component: workbenchRequestList
          },
          // 新建问题
          {
            path: '/taskman/workbench/createProblem',
            name: '/taskman/workbench/createProblem',
            component: workbenchProblemCreate
          },
          // 问题详情
          {
            path: '/taskman/workbench/detailProblem',
            name: '/taskman/workbench/detailProblem',
            component: workbenchProblemDetail
          },
          // 问题历史
          {
            path: '/taskman/workbench/problemHistory',
            name: '/taskman/workbench/problemHistory',
            component: workbenchProblemList
          },
          // 新建事件
          {
            path: '/taskman/workbench/createEvent',
            name: '/taskman/workbench/createEvent',
            component: workbenchEventCreate
          },
          // 事件详情
          {
            path: '/taskman/workbench/detailEvent',
            name: '/taskman/workbench/detailEvent',
            component: workbenchEventDetail
          },
          // 事件历史
          {
            path: '/taskman/workbench/eventHistory',
            name: '/taskman/workbench/eventHistory',
            component: workbenchEventList
          },
          // 新建变更
          {
            path: '/taskman/workbench/createChange',
            name: '/taskman/workbench/createChange',
            component: workbenchChangeCreate
          },
          // 变更详情
          {
            path: '/taskman/workbench/detailChange',
            name: '/taskman/workbench/detailChange',
            component: workbenchChangeDetail
          },
          // 变更历史
          {
            path: '/taskman/workbench/changeHistory',
            name: '/taskman/workbench/changeHistory',
            component: workbenchChangeList
          }
        ]
      },
      // 请求审计
      {
        path: '/taskman/requestAudit',
        name: '/taskman/requestAudit',
        component: workbenchRequestAudit
      }
    ]
  }
]
export default router
