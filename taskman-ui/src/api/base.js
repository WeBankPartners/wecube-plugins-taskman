import Vue from 'vue'
import axios from 'axios'
// import exportFile from '@/const/export-file'
import { getCookie } from '../pages/util/cookie'

export const baseURL = ''
export const req = axios.create({
  withCredentials: false,
  baseURL,
  timeout: 500000
})

const throwError = res => {
  Vue.prototype.$Notice.warning({
    title: 'Error',
    desc: (res.data && 'status:' + res.data.status + '<br/> message:' + res.data.message) || 'error',
    duration: 6
  })
}

req.interceptors.request.use(
  config => {
    return new Promise((resolve, reject) => {
      const lang = localStorage.getItem('lang') || 'zh-CN'
      if (lang === 'zh-CN') {
        config.headers['Accept-Language'] = 'zh-CN,zh;q=0.9,en;q=0.8'
      } else {
        config.headers['Accept-Language'] = 'en-US,en;q=0.9,zh;q=0.8'
      }
      const accessToken = getCookie('accessToken')
      if (accessToken && config.url !== '/auth/v1/api/login') {
        config.headers.Authorization = 'Bearer ' + accessToken
        resolve(config)
      } else {
        resolve(config)
      }
    })
  },
  error => {
    return Promise.reject(error)
  }
)

req.interceptors.response.use(
  res => {
    if (res.status === 200) {
      if (res.data.statusCode && res.data.statusCode !== 'OK') {
        const errorMes = res.data.statusMessage
        Vue.prototype.$Notice.error({
          title: 'Error',
          desc: errorMes,
          duration: 6
        })
      }
      return {
        ...res.data,
        user: res.headers['username'] || ' - '
      }
    } else {
      return {
        data: throwError(res)
      }
    }
  },
  err => {
    console.log(err)
  }
)

function setHeaders (obj) {
  Object.keys(obj).forEach(key => {
    req.defaults.headers.common[key] = obj[key]
  })
}

export default req

export { setHeaders }
