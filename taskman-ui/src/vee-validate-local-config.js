import Vue from 'vue'
import { extend, localize } from 'vee-validate'
import { required, email, min } from 'vee-validate/dist/rules'
import zh from 'vee-validate/dist/locale/zh_CN.json'
import en from 'vee-validate/dist/locale/en.json'

// Install required rule.
extend('required', required)

// Install email rule.
extend('email', email)

// Install min rule.
extend('min', min)

// Install English and Arabic localizations.
localize({
  zh_CN: {
    messages: zh.messages
    // names: {
    //   email: '邮箱地址',
    //   password: '密码'
    // },
    // fields: {
    //   password: {
    //     min: '{_field_} 太短'
    //   }
    // }
  },
  en: {
    messages: en.messages
    // names: {
    //   email: 'E-mail Address',
    //   password: 'Password'
    // },
    // fields: {
    //   password: {
    //     min: '{_field_} is too short, you want to get hacked?'
    //   }
    // }
  }
})

let LOCALE = 'zh_CN'
// this.locale = 'zh_CN'
localize('zh_CN')

// A simple get/set interface to manage our locale in components.
// This is not reactive, so don't create any computed properties/watchers off it.
Object.defineProperty(Vue.prototype, 'locale', {
  get () {
    return LOCALE
  },
  set (val) {
    LOCALE = val
    localize(val)
  }
})
