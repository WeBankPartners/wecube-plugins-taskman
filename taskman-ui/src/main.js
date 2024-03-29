import 'regenerator-runtime/runtime'
import Vue from 'vue'
import App from './App.vue'
import router from './router'
import ViewUI from 'view-design'
import 'view-design/dist/styles/iview.css'
import VueI18n from 'vue-i18n'
import locale from 'view-design/dist/locale/en-US'
import './locale/i18n'
import { ValidationProvider } from 'vee-validate'
import './vee-validate-local-config'
Vue.component('ValidationProvider', ValidationProvider)

Vue.config.productionTip = false

Vue.use(ViewUI, {
  transfer: true,
  size: 'default',
  VueI18n,
  locale
})

router.beforeEach((to, from, next) => {
  const ls = localStorage.getItem('taskman-accessToken')
  if (ls !== null || to.path === '/login') {
    next()
  } else {
    next('/login')
  }
})

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
