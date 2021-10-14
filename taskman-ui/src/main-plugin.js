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
const implicitRoute = {}
window.addImplicitRoute(implicitRoute)
window.addRoutersWithoutPermission(routerP, 'taskman')
window.addRoutes && window.addRoutes(router, 'taskman')
