// 防抖函数
export const debounce1 = (fn, delay) => {
  let timer = null
  let that = this
  return (...args) => {
    timer && clearTimeout(timer)
    timer = setTimeout(() => {
      fn.apply(that, args)
    }, delay)
  }
}
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

// 深拷贝
export const deepClone = obj => {
  let objClone = Array.isArray(obj) ? [] : {}
  if (obj && typeof obj === 'object') {
    for (let key in obj) {
      if (obj.hasOwnProperty(key)) {
        if (obj[key] && typeof obj[key] === 'object') {
          objClone[key] = deepClone(obj[key])
        } else {
          objClone[key] = obj[key]
        }
      }
    }
  }
  return objClone
}

// 数组对象数据去重，并返回重复的数据
export const uniqueArr = arr => {
  // 数组去重操作
  let newArrId = []
  let newArrObj = [] // 去重后的数组
  arr.forEach(item => {
    if (newArrId.indexOf(item.name) === -1) {
      newArrId.push(item.name)
      newArrObj.push(item)
    }
  })
  // 获取数组中重复的元素
  const arrId = arr.map(item => item.name)
  let sameArrId = []
  let sameArrObj = [] // 数组中重复的数据
  arrId.forEach(item => {
    if (arrId.indexOf(item) !== arrId.lastIndexOf(item) && sameArrId.indexOf(item) === -1) {
      sameArrId.push(item)
    }
  })
  sameArrId.forEach(i => {
    try {
      arr.forEach(j => {
        if (i === j.name) {
          sameArrObj.push(j)
          throw new Error('end')
        }
      })
    } catch (e) {}
  })
  return { arr: newArrObj, sameArr: sameArrObj }
}

// 返回数组中第一个重复元素的索引
export const findFirstDuplicateIndex = arr => {
  const map = new Map()
  for (let i = 0; i < arr.length; i++) {
    if (map.has(arr[i]) && map.get(arr[i]) > 0) {
      return i
    }
    map.set(arr[i], (map.get(arr[i]) || 0) + 1)
  }
  return -1 // 如果没有重复元素，则返回-1
}
