// 防抖函数
export const debounce = (fn, delay) => {
  let timer = null
  let that = this
  return (...args) => {
    timer && clearTimeout(timer)
    timer = setTimeout(() => {
      fn.apply(that, args)
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
