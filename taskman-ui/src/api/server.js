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

export const getManagementRoles = () => req.get('/taskman/api/v1/user/roles')
export const getUserRoles = () => req.get('/taskman/api/v1/role/list')
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

export const getTemplateNodesForTemp = requestTemplateId =>
  req.get(`/taskman/api/v1/process-nodes/${requestTemplateId}/template`)
export const getTemplateNodesForRequest = requestTemplateId =>
  req.get(`/taskman/api/v1/process-nodes/${requestTemplateId}/bind`)

export const saveTaskForm = (requestTemplateId, data) =>
  req.post(`/taskman/api/v1/task-template/${requestTemplateId}`, data)

export const confirmTemplate = requestTemplateId =>
  req.post(`/taskman/api/v1/request-template/confirm/${requestTemplateId}`)

export const getTemplateList = data => req.post('/taskman/api/v1/request-template/query', data)
export const deleteTemplate = data => req.delete('/taskman/api/v1/request-template', data)
export const forkTemplate = requestTemplateId => req.post(`/taskman/api/v1/request-template/fork/${requestTemplateId}`)
export const getRequestTemplateAttrs = requestTemplateId =>
  req.get(`/taskman/api/v1/request-template/${requestTemplateId}/attrs/get`)

export const getRequestFormTemplateData = requestTemplateId =>
  req.get(`/taskman/api/v1/request-form-template/${requestTemplateId}`)
export const getTaskFormDataByNodeId = (requestTemplateId, nodeId) =>
  req.get(`/taskman/api/v1/task-template/${requestTemplateId}/${nodeId}`)

export const getTemplateByUser = () => req.get('/taskman/api/v1/user/request-template')
export const createRequest = (requestId, data) => req.post('/taskman/api/v1/request', data)
export const updateRequest = (requestId, data) => req.put(`/taskman/api/v1/request/${requestId}`, data)
export const getRootEntity = params => req.get('/taskman/api/v1/entity/data', params)
export const getEntityData = params => req.get('/taskman/api/v1/request-data/preview', params)
export const saveEntityData = (requestId, params) =>
  req.post(`/taskman/api/v1/request-data/save/${requestId}/data`, params)
export const getBindData = requestId => req.get(`/taskman/api/v1/request-data/get/${requestId}/data`)
export const getBindRelate = requestId => req.get(`/taskman/api/v1/request-data/get/${requestId}/bing`)
export const saveRequest = (requestId, data) => req.post(`/taskman/api/v1/request-data/save/${requestId}/bing`, data)
export const updateRequestStatus = (requestId, status, params) =>
  req.post(`/taskman/api/v1/request-status/${requestId}/${status}`, params)
export const requestListForDraftInitiated = params => req.post(`/taskman/api/v1/user/request/use`, params)
export const requestListForHandle = params => req.post(`/taskman/api/v1/user/request/mgmt`, params)
export const deleteRequest = id => req.delete(`/taskman/api/v1/request/${id}`)
export const terminateRequest = id => req.post(`/taskman/api/v1/request/terminate/${id}`)
export const startRequest = (requestId, data) => req.post(`/taskman/api/v1/request/start/${requestId}`, data)
export const getRequestInfo = requestId => req.get(`/taskman/api/v1/request/${requestId}`)

export const getRefOptions = (requestId, attr, params) =>
  req.post(`/taskman/api/v1/request-data/reference/query/${attr}/${requestId}`, params)

export const taskList = params => req.post(`/taskman/api/v1/task/list`, params)
export const getTaskDetail = taskId => req.get(`/taskman/api/v1/task/detail/${taskId}`)
export const saveTaskData = (taskId, data) => req.post(`/taskman/api/v1/task/save/${taskId}`, data)
export const changeTaskStatus = (operation, taskId) => req.post(`/taskman/api/v1/task/status/${operation}/${taskId}`)
export const commitTaskData = (taskId, data) => req.post(`/taskman/api/v1/task/approve/${taskId}`, data)
export const getRequestDetail = requestId => req.get(`/taskman/api/v1/request/detail/${requestId}`)
export const getHandlerRoles = params => req.get(`/taskman/api/v1/role/user/list`, params)
export const getTemplateTags = requestTemplateGroup =>
  req.get(`/taskman/api/v1/request-template/tags/${requestTemplateGroup}`)
export const confirmUploadTemplate = confirmToken =>
  req.post(`/taskman/api/v1/request-template/import-confirm/${confirmToken}`)
export const deleteAttach = fileId => req.delete(`/taskman/api/v1/request/attach-file/remove/${fileId}`)
export const reRequest = fileId => req.post(`/taskman/api/v1/request/copy/${fileId}`)
export const requestParent = requestId => req.get(`/taskman/api/v1/request-parent/get?requestId=${requestId}`)
export const enableTemplate = templateId => req.post(`/taskman/api/v1/request-template/enable/${templateId}`)
export const disableTemplate = templateId => req.post(`/taskman/api/v1/request-template/disable/${templateId}`)
// taskman重构
// 选择模板列表
export const getTemplateTree = () => req.get('/taskman/api/v2/user/request-template')
// 模板收藏
export const collectTemplate = templateId => req.post(`/taskman/api/v1/user/template/collect/${templateId}`)
// 取消模板收藏
export const uncollectTemplate = templateId => req.delete(`/taskman/api/v1/user/template/collect/${templateId}`)
// 模板收藏列表
export const collectTemplateList = params => req.post(`/taskman/api/v1/user/template/collect/query`, params)
// 工作台看板数量
export const overviewData = () => req.post(`/taskman/api/v1/user/platform`)
// 工作台列表
export const getPlatformList = params => req.post(`/taskman/api/v1/user/platform/list`, params)
// 获取工作台筛选数据集合
export const getPlatformFilter = params => req.post(`/taskman/api/v1/user/platform/filter-item`, params)
// 工作台转给我
export const tansferToMe = templateId => req.post(`/taskman/api/v1/request/handler/${templateId}`)
// 新建发布-发布信息获取
export const getCreateInfo = params => req.post(`/taskman/api/v2/request`, params)
// 新建发布-请求进度
export const getProgressInfo = params => req.post(`/taskman/api/v1/request/progress`, params)
// 新建发布-保存数据
export const savePublishData = (requestId, params) =>
  req.post(`/taskman/api/v2/request-data/save/${requestId}/data`, params)
// 新建发布详情数据
export const getPublishInfo = requestId => req.get(`/taskman/api/v2/request/detail/${requestId}`)
// 发布历史页面
export const getPublishList = params => req.post(`/taskman/api/v1/request/history/list`, params)

// 查询流程图
export const getFlowByTemplateId = templateId => req.get(`/taskman/api/v1/request/process/definitions/${templateId}`)
export const getFlowByInstanceId = instanceId => req.get(`/taskman/api/v1/request/process/instances/${instanceId}`)
export const getNodeContextByNodeId = (instanceId, nodeId) =>
  req.post(`/taskman/api/v1/request/workflow/task_node/${instanceId}/${nodeId}`)
