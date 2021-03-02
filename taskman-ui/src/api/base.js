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
        config.headers.Authorization = 'Bearer ' + 'eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJ1bWFkbWluIiwiaWF0IjoxNjE0Njc4NTU3LCJ0eXBlIjoiYWNjZXNzVG9rZW4iLCJjbGllbnRUeXBlIjoiVVNFUiIsImV4cCI6MTYxNDY3OTc1NywiYXV0aG9yaXR5IjoiW1NVQl9TWVNURU0sQ01EQl9BRE1JTixTVVBFUl9BRE1JTixJTVBMRU1FTlRBVElPTl9XT1JLRkxPV19FWEVDVVRJT04sSU1QTEVNRU5UQVRJT05fQVJUSUZBQ1RfTUFOQUdFTUVOVCxNT05JVE9SX01BSU5fREFTSEJPQVJELE1PTklUT1JfTUVUUklDX0NPTkZJRyxNT05JVE9SX0FMQVJNX0NPTkZJRyxNT05JVE9SX0FMQVJNX01BTkFHRU1FTlQsQ09MTEFCT1JBVElPTl9QTFVHSU5fTUFOQUdFTUVOVCxDT0xMQUJPUkFUSU9OX1dPUktGTE9XX09SQ0hFU1RSQVRJT04sQURNSU5fU1lTVEVNX1BBUkFNUyxBRE1JTl9SRVNPVVJDRVNfTUFOQUdFTUVOVCxBRE1JTl9VU0VSX1JPTEVfTUFOQUdFTUVOVCxBRE1JTl9DTURCX01PREVMX01BTkFHRU1FTlQsQ01EQl9BRE1JTl9CQVNFX0RBVEFfTUFOQUdFTUVOVCxBRE1JTl9RVUVSWV9MT0csTUVOVV9BRE1JTl9QRVJNSVNTSU9OX01BTkFHRU1FTlQsTUVOVV9JRENfUkVTT1VSQ0VfUExBTk5JTkcsTUVOVV9DTURCX0FETUlOX0JBU0VfREFUQV9NQU5BR0VNRU5ULE1FTlVfREVTSUdOSU5HX0NJX0lOVEVHUkFURURfUVVFUllfRVhFQ1VUSU9OLE1FTlVfQVBQTElDQVRJT05fREVQTE9ZTUVOVF9ERVNJR04sTUVOVV9ERVNJR05JTkdfQ0lfREFUQV9NQU5BR0VNRU5ULE1FTlVfREVTSUdOSU5HX0NJX0lOVEVHUkFURURfUVVFUllfTUFOQUdFTUVOVCxNRU5VX0lEQ19QTEFOTklOR19ERVNJR04sTUVOVV9BRE1JTl9RVUVSWV9MT0csTUVOVV9BUFBMSUNBVElPTl9BUkNISVRFQ1RVUkVfUVVFUlksTUVOVV9BUFBMSUNBVElPTl9BUkNISVRFQ1RVUkVfREVTSUdOLE1FTlVfREVTSUdOSU5HX0NJX0RBVEFfRU5RVUlSWSxDQVBBQ0lUWV9NT0RFTCxDQVBBQ0lUWV9GT1JFQ0FTVCxBRE1JTl9JVFNfREFOR0VST1VTX0NPTkZJRyxKT0JTX1NFUlZJQ0VfQ0FUQUxPR19NQU5BR0VNRU5ULE1PTklUT1JfQ1VTVE9NX0RBU0hCT0FSRCxKT0JTX1RFTVBMQVRFX0dST1VQX01BTkFHRU1FTlQsSk9CU19URU1QTEFURV9NQU5BR0VNRU5ULEpPQlNfUkVRVUVTVF9NQU5BR0VNRU5ULE1FTlVfQURNSU5fQ01EQl9NT0RFTF9NQU5BR0VNRU5ULEpPQlNfVEFTS19NQU5BR0VNRU5ULElNUExFTUVOVEFUSU9OX0JBVENIX0VYRUNVVElPTixBRE1JTl9URVJSQUZPUk1fQ09ORklHLElNUExFTUVOVEFUSU9OX1RFUk1JTkFMLEFETUlOX1RFUk1JTkFMX0NPTkZJR10ifQ.BDd7ViW5QTSpPinZ-dzChsaGauJESSXhPH3NUssUiYvrF2G4J3nujhsNX3sOVCf1FIMjBF9hzKIUmfwaeFRpjw'
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
