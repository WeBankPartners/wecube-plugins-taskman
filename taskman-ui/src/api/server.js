import { req as request, baseURL } from './base'
import { pluginErrorMessage } from './base-plugin'
let req = request
if (window.request) {
  req = {
    post: (url, ...params) => pluginErrorMessage(window.request.post(baseURL + url, ...params)),
    get: (url, ...params) => pluginErrorMessage(window.request.get(baseURL + url, ...params)),
    delete: (url, ...params) => pluginErrorMessage(window.request.delete(baseURL + url, ...params)),
    put: (url, ...params) => pluginErrorMessage(window.request.put(baseURL + url, ...params)),
    patch: (url, ...params) => pluginErrorMessage(window.request.patch(baseURL + url, ...params))
  }
}

export const getTempGroupList = data => req.post('/taskman/api/v1/request-template-group/query', data)
export const createTempGroup = data => req.post('/taskman/api/v1/request-template-group', data)
export const updateTempGroup = data => req.put('/taskman/api/v1/request-template-group', data)
export const deleteTempGroup = data => req.delete('/taskman/api/v1/request-template-group', data)

export const getManagementRoles = () => req.get('/taskman/api/v1/role/list')
export const getUserRoles = () => req.get('/taskman/api/v1/user/roles')
export const getProcess = () => req.get('/taskman/api/v1/process/list')

export const createTemp = data => req.post('/taskman/api/v1/request-template', data)
export const updateTemp = data => req.put('/taskman/api/v1/request-template', data)
export const deleteTemp = data => req.delete('/taskman/api/v1/request-template', data)

export const getFormList = requestTemplateId =>
  req.get(`/taskman/api/v1/request-template/${requestTemplateId}/attrs/list`)

export const saveAttrs = (requestTemplateId, data) =>
  req.put(`/taskman/api/v1/request-template/${requestTemplateId}/attrs/update`, data)

export const getSelectedForm = requestTemplateId =>
  req.get(`/taskman/api/v1/request-template/${requestTemplateId}/attrs/get`)

export const saveRequsetForm = (requestTemplateId, data) =>
  req.post(`/taskman/api/v1/request-form-template/${requestTemplateId}`, data)

export const getTemplateNodes = requestTemplateId => req.get(`/taskman/api/v1/process-nodes/${requestTemplateId}`)

export const saveTaskForm = (requestTemplateId, data) =>
  req.post(`/taskman/api/v1/task-template/${requestTemplateId}`, data)

export const confirmTemplate = requestTemplateId =>
  req.post(`/taskman/api/v1/request-template/confirm/${requestTemplateId}`)
