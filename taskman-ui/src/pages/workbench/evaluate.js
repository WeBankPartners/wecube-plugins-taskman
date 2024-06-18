export const evaluateCondition = (item, inputVal) => {
  if (Array.isArray(inputVal)) {
    // 过滤掉空数据
    inputVal = inputVal.filter(i => i !== '')
  }
  const { operator, value } = item
  const emptyFlag = !inputVal || (Array.isArray(inputVal) && inputVal.length === 0)
  // 除了符号为空值类型，其它空值不参与条件计算
  if (emptyFlag && operator !== 'empty') return false
  switch (operator) {
    // 等于
    case 'eq':
      if (isDateTime(value)) {
        const time1 = new Date(inputVal).getTime()
        const time2 = new Date(value).getTime()
        return time1 === time2
      } else {
        return inputVal === value
      }
    // 不等于
    case 'neq':
      if (isDateTime(value)) {
        const time1 = new Date(inputVal).getTime()
        const time2 = new Date(value).getTime()
        return time1 !== time2
      } else {
        return inputVal !== value
      }
    // 小于
    case 'lt':
      if (isDateTime(value)) {
        const time1 = new Date(inputVal).getTime()
        const time2 = new Date(value).getTime()
        return time1 < time2
      } else if (!isNaN(value) && !isNaN(inputVal)) {
        return +inputVal < +value
      } else {
        return inputVal < value
      }
    // 大于
    case 'gt':
      if (isDateTime(value)) {
        const time1 = new Date(inputVal).getTime()
        const time2 = new Date(value).getTime()
        return time1 > time2
      } else if (!isNaN(value) && !isNaN(inputVal)) {
        return +inputVal > +value
      } else {
        return inputVal > value
      }
    // 包含
    case 'contains':
      return inputVal.includes(value)
    // 匹配开始
    case 'startsWith':
      return inputVal.startsWith(value)
    // 匹配结束
    case 'endsWith':
      return inputVal.endsWith(value)
    // 包含全部
    case 'containsAll':
      return evaluateContainsAll(value, inputVal)
    // 包含任意
    case 'containsAny':
      return evaluateContainsAny(value, inputVal)
    // 不包含
    case 'notContains':
      return evaluateNotContains(value, inputVal)
    case 'range':
      return evaluateRangeDateTime(value, inputVal)
    case 'empty':
      return isEmpty(inputVal)
    case 'notEmpty':
      return !isEmpty(inputVal)
    default:
      return false
  }
}

const isDateTime = val => {
  if (typeof val === 'string' && val.split(':') && val.split(':').length === 2) {
    return true
  } else {
    return false
  }
}

const isEmpty = val => {
  if (!val || (Array.isArray(val) && val.length === 0)) {
    return true
  } else {
    return false
  }
}

// 包含全部
const evaluateContainsAll = (value, inputVal) => {
  let conditionArr
  if (Array.isArray(value)) {
    conditionArr = value
  } else {
    conditionArr = value.split(',') || []
  }
  if (Array.isArray(inputVal)) {
    return conditionArr.every(i => inputVal.includes(i))
  } else {
    return conditionArr.every(i => inputVal === i)
  }
}

// 包含任意
const evaluateContainsAny = (value, inputVal) => {
  let conditionArr
  if (Array.isArray(value)) {
    conditionArr = value
  } else {
    conditionArr = value.split(',') || []
  }
  if (Array.isArray(inputVal)) {
    return conditionArr.some(i => inputVal.includes(i))
  } else {
    return conditionArr.some(i => inputVal === i)
  }
}

// 不包含
const evaluateNotContains = (value, inputVal) => {
  let conditionArr
  if (Array.isArray(value)) {
    conditionArr = value
  } else {
    conditionArr = value.split(',') || []
  }
  if (Array.isArray(inputVal)) {
    return conditionArr.every(i => !inputVal.includes(i))
  } else {
    return conditionArr.every(i => !inputVal === i)
  }
}

// 时间范围
const evaluateRangeDateTime = (value, inputVal) => {
  if (value && value.length > 0) {
    const timeStart = new Date(value[0]).getTime()
    const timeEnd = new Date(value[1]).getTime()
    const timeInput = new Date(inputVal).getTime()
    if (timeInput <= timeEnd && timeInput >= timeStart) {
      return true
    } else {
      return false
    }
  } else {
    return false
  }
}
