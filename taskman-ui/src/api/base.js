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
        config.headers.Authorization = 'Bearer ' + 'eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJ1bWFkbWluIiwiaWF0IjoxNjA4ODg1OTUwLCJ0eXBlIjoiYWNjZXNzVG9rZW4iLCJjbGllbnRUeXBlIjoiVVNFUiIsImV4cCI6MTYwODg4NzE1MCwiYXV0aG9yaXR5IjoiW1NVUEVSX0FETUlOLElNUExFTUVOVEFUSU9OX1dPUktGTE9XX0VYRUNVVElPTixJTVBMRU1FTlRBVElPTl9CQVRDSF9FWEVDVVRJT04sSU1QTEVNRU5UQVRJT05fQVJUSUZBQ1RfTUFOQUdFTUVOVCxNT05JVE9SX01BSU5fREFTSEJPQVJELE1PTklUT1JfTUVUUklDX0NPTkZJRyxNT05JVE9SX0NVU1RPTV9EQVNIQk9BUkQsTU9OSVRPUl9BTEFSTV9DT05GSUcsTU9OSVRPUl9BTEFSTV9NQU5BR0VNRU5ULENPTExBQk9SQVRJT05fUExVR0lOX01BTkFHRU1FTlQsQ09MTEFCT1JBVElPTl9XT1JLRkxPV19PUkNIRVNUUkFUSU9OLEFETUlOX1NZU1RFTV9QQVJBTVMsQURNSU5fUkVTT1VSQ0VTX01BTkFHRU1FTlQsQURNSU5fVVNFUl9ST0xFX01BTkFHRU1FTlQsQURNSU5fQ01EQl9NT0RFTF9NQU5BR0VNRU5ULENNREJfQURNSU5fQkFTRV9EQVRBX01BTkFHRU1FTlQsQURNSU5fUVVFUllfTE9HLE1FTlVfQURNSU5fUEVSTUlTU0lPTl9NQU5BR0VNRU5ULEpPQlNfUkVRVUVTVF9NQU5BR0VNRU5ULE1FTlVfREVTSUdOSU5HX0NJX0RBVEFfRU5RVUlSWSxNRU5VX0RFU0lHTklOR19DSV9JTlRFR1JBVEVEX1FVRVJZX0VYRUNVVElPTixNRU5VX0RFU0lHTklOR19DSV9EQVRBX01BTkFHRU1FTlQsTUVOVV9ERVNJR05JTkdfQ0lfSU5URUdSQVRFRF9RVUVSWV9NQU5BR0VNRU5ULE1FTlVfSURDX1BMQU5OSU5HX0RFU0lHTixNRU5VX0lEQ19SRVNPVVJDRV9QTEFOTklORyxNRU5VX0NNREJfQURNSU5fQkFTRV9EQVRBX01BTkFHRU1FTlQsTUVOVV9BRE1JTl9RVUVSWV9MT0csTUVOVV9BUFBMSUNBVElPTl9BUkNISVRFQ1RVUkVfREVTSUdOLE1FTlVfQVBQTElDQVRJT05fQVJDSElURUNUVVJFX1FVRVJZLE1FTlVfQVBQTElDQVRJT05fREVQTE9ZTUVOVF9ERVNJR04sTUVOVV9BRE1JTl9DTURCX01PREVMX01BTkFHRU1FTlQsSk9CU19UQVNLX01BTkFHRU1FTlQsSk9CU19TRVJWSUNFX0NBVEFMT0dfTUFOQUdFTUVOVCxKT0JTX1RFTVBMQVRFX0dST1VQX01BTkFHRU1FTlQsSk9CU19URU1QTEFURV9NQU5BR0VNRU5UXSJ9.paKpbeu2AmVoxDN7gV0PiwL-ztsR6Y6wWzmK0SawugnP85Zxy7SpP4YX4zxdYtMUIWsVhENGej8vkqTy1cL10g'
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
