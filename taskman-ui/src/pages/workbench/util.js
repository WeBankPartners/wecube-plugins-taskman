// 校验entitytable表单必须有一条数据
export const noChooseCheck = (data, entityTable) => {
  let tabIndex = ''
  let result = true
  data.forEach((requestData, index) => {
    if (requestData.value && requestData.value.length === 0) {
      tabIndex = index
      result = false
    }
  })
  tabIndex && entityTable && entityTable.validTable(tabIndex)
  return result
}

// 校验表格数据必填项
export const requiredCheck = (data, entityTable) => {
  let tabIndex = ''
  let result = true
  data.forEach((requestData, index) => {
    let requiredName = []
    requestData.title.forEach(t => {
      if (t.required === 'yes') {
        requiredName.push(t.name)
      }
    })
    requestData.value.forEach(v => {
      requiredName.forEach(key => {
        let val = v.entityData[key]
        if (Array.isArray(val)) {
          if (val.length === 0) {
            result = false
            if (tabIndex === '') {
              tabIndex = index
            }
          }
        } else {
          if (val === '' || val === undefined) {
            result = false
            if (tabIndex === '') {
              tabIndex = index
            }
          }
        }
      })
    })
  })
  tabIndex && entityTable && entityTable.validTable(tabIndex)
  return result
}

// 审批任务流程必填校验
export const approvalCheck = data => {
  let result = true
  data.forEach(i => {
    if (i.handleTemplates && i.handleTemplates.length > 0) {
      i.handleTemplates.forEach(j => {
        if (!j.role || (!j.handler && !['system', 'claim'].includes(j.handlerType))) {
          result = false
        }
      })
    }
  })
  return result
}
