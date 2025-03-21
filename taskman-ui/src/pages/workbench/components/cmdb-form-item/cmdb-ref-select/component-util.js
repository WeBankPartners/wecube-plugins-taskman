export const finalDataForRequest = oriData => {
  if (Array.isArray(oriData)) {
    oriData.map(item => {
      const finalForArrayData = managementData(item)
      return finalForArrayData
    })
    return oriData
  } else {
    const finalForJsonData = managementData(oriData)
    return finalForJsonData
  }
}

const managementData = item => {
  const keys = Object.keys(item)
  keys.forEach(key => {
    if (Array.isArray(item[key])) {
      const guidGroup = item[key].map(arr => {
        return arr.guid || arr
      })
      item[key] = guidGroup.join(',')
    } else {
      if (isJSON(item[key])) {
        item[key] = item[key].guid
      }
    }
  })
  delete item.nextOperations
  delete item.isRowEditable
  delete item.weTableForm
  delete item.weTableRowId
  delete item.isNewAddedRow
  return item
}

const isJSON = jsons => {
  try {
    if (typeof jsons === 'object' && jsons) {
      return true
    } else {
      return false
    }
  } catch (e) {
    return false
  }
}

export const components = {
  number: {
    component: 'Input',
    type: 'number'
  },
  datetime: {
    component: 'DatePicker',
    type: 'datetimerange'
  },
  text: {
    component: 'Input',
    type: 'text'
  },
  select: {
    component: 'WeCMDBSelect',
    options: []
  },
  ref: {
    component: 'WeCMDBRefSelect',
    highlightRow: true
  },
  extRef: {
    component: 'WeCMDBRefSelect',
    highlightRow: true
  },
  multiSelect: {
    component: 'WeCMDBSelect',
    options: []
  },
  multiRef: {
    component: 'WeCMDBRefSelect'
  },
  textArea: {
    component: 'Input',
    type: 'text'
  },
  password: {
    component: 'WeCMDBCIPassword'
  }
}
