import 'regenerator-runtime/runtime'
import Vue from 'vue'
import App from './App.vue'
import router from './router'
import ViewUI from 'view-design'
import VueI18n from 'vue-i18n'
import './styles/index.less'
import { i18n } from './locale/i18n/index.js'
import viewDesignEn from 'view-design/dist/locale/en-US'
import viewDesignZh from 'view-design/dist/locale/zh-CN'
import { ValidationProvider } from 'vee-validate'
import './vee-validate-local-config'
import commonUI from 'wecube-common-ui'
import 'wecube-common-ui/lib/wecube-common-ui.css'
import { getCookie } from '@/pages/util/cookie'
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

router.beforeEach((to, from, next) => {
  const ls = getCookie('accessToken')
  if (ls || to.path === '/login') {
    next()
  } else {
    next('/login')
  }
})

new Vue({
  router,
  render: h => h(App),
  i18n
}).$mount('#app')
