import 'regenerator-runtime/runtime'
import router from './router-plugin'
import routerP from './router-plugin-p'
import 'view-design/dist/styles/iview.css'
import './locale/i18n'
import zhCN from '@/locale/i18n/zh-CN.json'
import enUS from '@/locale/i18n/en-US.json'
import { ValidationProvider } from 'vee-validate'
import './vee-validate-local-config'
window.component('ValidationProvider', ValidationProvider)

window.locale('zh-CN', zhCN)
window.locale('en-US', enUS)
const implicitRoute = {
  '/taskman/template-mgmt': {
    parentBreadcrumb: { 'zh-CN': 'TaskMan', 'en-US': 'TaskMan' },
    childBreadcrumb: { 'zh-CN': '模板', 'en-US': 'Template' }
  },
  '/taskman/template-group': {
    parentBreadcrumb: { 'zh-CN': 'TaskMan', 'en-US': 'TaskMan' },
    childBreadcrumb: { 'zh-CN': '模板组', 'en-US': 'Template Group' }
  },
  '/templateManagementIndex': {
    parentBreadcrumb: { 'zh-CN': 'TaskMan', 'en-US': 'TaskMan' },
    childBreadcrumb: { 'zh-CN': '模板管理', 'en-US': 'Template Management' }
  },
  '/taskman/request-mgmt': {
    parentBreadcrumb: { 'zh-CN': 'TaskMan', 'en-US': 'TaskMan' },
    childBreadcrumb: { 'zh-CN': '请求', 'en-US': 'Request' }
  },
  '/taskman/task-mgmt': {
    parentBreadcrumb: { 'zh-CN': 'TaskMan', 'en-US': 'TaskMan' },
    childBreadcrumb: { 'zh-CN': '任务', 'en-US': 'Task' }
  },
  '/taskMgmtIndex': {
    parentBreadcrumb: { 'zh-CN': 'TaskMan', 'en-US': 'TaskMan' },
    childBreadcrumb: { 'zh-CN': '任务管理', 'en-US': 'Task Management' }
  },
  '/requestManagementIndex': {
    parentBreadcrumb: { 'zh-CN': 'TaskMan', 'en-US': 'TaskMan' },
    childBreadcrumb: { 'zh-CN': '请求管理', 'en-US': 'Request Management' }
  }
}
window.addImplicitRoute(implicitRoute)
window.addRoutersWithoutPermission(routerP, 'taskman')
window.addRoutes && window.addRoutes(router, 'taskman')
