// 防抖函数
// export const debounce = (fn, delay) => {
//   let timer = null
//   let that = this
//   return (...args) => {
//     timer && clearTimeout(timer)
//     timer = setTimeout(() => {
//       fn.apply(that, args)
//     }, delay)
//   }
// }
export function debounce (fn, delay = 500) {
  let timer = null
  return function () {
    const args = arguments
    if (timer) {
      clearTimeout(timer)
    }
    timer = setTimeout(() => {
      fn.apply(this, [...args])
    }, delay)
  }
}

// 截流函数
export const throttle = (fn, delay) => {
  let timer = null
  let that = this
  return args => {
    if (timer) return
    timer = setTimeout(() => {
      fn.apply(that, args)
      timer = null
    }, delay)
  }
}
