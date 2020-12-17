import { setCookie, getCookie } from './cookie.js'
import Vue from 'vue'
import axios from 'axios'
export const baseURL = ''
export const req = axios.create({
  withCredentials: true,
  baseURL,
  timeout: 50000
})

const throwError = res => new Error(res.message || 'error')
req.interceptors.response.use(
  res => {
    if (res.status === 200) {
      if (res.data.status.startsWith('ERR')) {
        Vue.prototype.$Notice.error({
          title: 'Error',
          desc: res.data.message,
          duration: 0
        })
      }
      // if (!res.headers['username']) {
      //   window.location.href = '/wecmdb/logout'
      // }
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
  error => {
    const { response } = error
    Vue.prototype.$Notice.error({
      title: 'error',
      desc:
        (response.data &&
          'status:' +
            response.data.status +
            '<br/> error:' +
            response.data.error +
            '<br/> message:' +
            response.data.message) ||
        'error'
    })
    return new Promise((resolve, reject) => {
      resolve({
        data: throwError(error)
      })
    })
  }
)

let refreshRequest = null
req.interceptors.request.use(
  config => {
    return new Promise((resolve, reject) => {
        config.headers.Authorization = 'Bearer ' + 'eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiIxMTY4IiwiaWF0IjoxNTc4MzA1NzAyLCJ0eXBlIjoicmVmcmVzaFRva2VuIiwiY2xpZW50VHlwZSI6IlVTRVIiLCJleHAiOjE1NzgzMDc1MDJ9.dnCGb91Z9YDiUX6YBlpaZ7yakPsXNPVxSNAuT0LeM_2qpPkcztqdswBEe-01nnCNJlS_jMm1GPrHJrdaYRQSyQ'
        resolve(config)
    })
  },
  error => {
    return Promise.reject(error)
  }
)

function setHeaders (obj) {
  Object.keys(obj).forEach(key => {
    req.defaults.headers.common[key] = obj[key]
  })
}

export { setHeaders }
