import Vue from 'vue'
// import VueI18n from 'vue-i18n'

import zh from 'view-design/dist/locale/zh-CN'
import en from 'view-design/dist/locale/en-US'

import ZH from './zh-CN.json'
import EN from './en-US.json'

// Vue.use(VueI18n)

// const messages = {
//   'zh-CN': Object.assign(ZH, zh),
//   'en-US': Object.assign(EN, en)
// }

// export const i18n = new VueI18n({
//   locale:
//     localStorage.getItem('lang') || (navigator.language || navigator.userLanguage === 'zh-CN' ? 'zh-CN' : 'en-US'),
//   messages
// })

// Vue.locale('zh-CN', Object.assign(zh, ZH))
// Vue.locale('en-US', Object.assign(en, EN))

const navLang = navigator.language

const localLang = navLang === 'zh-CN' || navLang === 'en-US' ? navLang : false

const lang = window.localStorage.getItem('lang') || localLang || 'en-US'
Vue.config.lang = lang
