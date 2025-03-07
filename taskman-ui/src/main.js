/*
 * @Author: wanghao7717 792974788@qq.com
 * @Date: 2023-10-12 17:16:35
 * @LastEditors: wanghao7717 792974788@qq.com
 * @LastEditTime: 2025-03-06 19:44:15
 */
/* eslint-disable camelcase */
import 'regenerator-runtime/runtime'
import Vue from 'vue'
import App from './App.vue'
import VueRouter from 'vue-router'
import { routes, dynaticRoutes } from './router'
import ViewUI from 'view-design'
import './styles/index.less'
import VueI18n from 'vue-i18n'
import { i18n } from './locale/i18n/index.js'
import viewDesignEn from 'view-design/dist/locale/en-US'
import viewDesignZh from 'view-design/dist/locale/zh-CN'
import { ValidationProvider } from 'vee-validate'
import './vee-validate-local-config'
import commonUI from 'wecube-common-ui'
import 'wecube-common-ui/lib/wecube-common-ui.css'
import { getCookie } from '@/pages/util/cookie'

import { implicitRoute, routerP } from './config.js'
Vue.component('ValidationProvider', ValidationProvider)
// 引用wecube公共组件
Vue.use(commonUI)

Vue.config.productionTip = false
Vue.prototype.$bus = new Vue()

Vue.use(ViewUI, {
  transfer: true,
  size: 'default',
  VueI18n,
  locale: i18n.locale === 'en-US' ? viewDesignEn : viewDesignZh
})

// 判断是否在qiankun的运行环境下
if (window.__POWERED_BY_QIANKUN__) {
  // eslint-disable-next-line no-undef
  __webpack_public_path__ = window.__INJECTED_PUBLIC_PATH_BY_QIANKUN__
}

let instance = null
let router = null
function render (props = {}) {
  const { container } = props
  router = new VueRouter({
    base: window.__POWERED_BY_QIANKUN__ ? process.env.MICRO_APP_URL : '/',
    mode: 'hash',
    routes
  })
  router.beforeEach((to, from, next) => {
    const ls = getCookie('accessToken')
    if (ls || to.path === '/login') {
      next()
    } else {
      next('/login')
    }
  })
  instance = new Vue({
    router,
    render: h => h(App),
    i18n
  }).$mount(container ? container.querySelector('#app') : '#app', true)
  Vue.prototype.$qiankunProps = props
  Vue.prototype.$qiankunProps.setGlobalState({
    implicitRoute: implicitRoute,
    childRouters: routerP,
    routes: dynaticRoutes
  })
}

if (!window.__POWERED_BY_QIANKUN__) {
  console.log('独立运行')
  render()
}

// 各个生命周期，只会在微应用初始化的时候调用一次，下次进入微应用重新进入是会直接调用mount钩子，不会再重复调用bootstrap
export async function bootstrap () {
  console.log('[vue] vue app bootstraped')
}
// 应用每次进入都会调用mount方法，通常在这里触发应用的渲染方法
export async function mount (props) {
  // 接受主应用参数
  console.log('[vue] props from main framework', props)
  render(props)
}
// 应用每次切除/注销会调用的方法，在这里会注销微应用的应用实例
export async function unmount () {
  instance.$destroy()
  instance.$el.innerHTML = ''
  instance = null
}
