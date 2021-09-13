import 'regenerator-runtime/runtime'
import router from './router-plugin'
import 'view-design/dist/styles/iview.css'
import './locale/i18n'
import zhCN from '@/locale/i18n/zh-CN.json'
import enUS from '@/locale/i18n/en-US.json'

window.locale('zh-CN', zhCN)
window.locale('en-US', enUS)
const implicitRoute = {}
window.addImplicitRoute(implicitRoute)
window.addRoutes && window.addRoutes(router, 'dangerous')
