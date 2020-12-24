import {req} from "./base";
const request = window.request ? window.request : req
// export const queryServiceRequest = data => request.post(`/service-mgmt/v1/service-requests/query`, data);
// export const getAllServiceRequest = () => request.get(`/service-mgmt/v1/service-requests/retrieve`);
// export const createServiceRequest = data => request.post(`/service-mgmt/v1/service-requests`, data);
// export const updateServiceRequest = data => request.put(`/service-mgmt/v1/service-requests/${data.id}/update`, data);
// export const getAllAvailableServiceTemplate = () => request.get(`/service-mgmt/v1/service-request-templates/available`);
// export const taskProcess = data => request.put(`/service-mgmt/v1/tasks/${data.taskId}/process`, data);
// export const queryMyTask = data => request.post(`/service-mgmt/v1/tasks/my-tasks/query`, data);
// export const taskTakeover = data => request.put(`/service-mgmt/v1/tasks/${data.taskId}/takeover`, data);
// export const getCurrentUserRoles = () => request.get(`/service-mgmt/v1/core-resources/users/current-user/roles`);

// export const getAllProcessDefinitionKeys = () => request.get(`/service-mgmt/v1/core-resources/workflow/process-definition-keys`);
// export const getAllAvailableServiceCatalogues = () => request.get(`/service-mgmt/v1/service-catalogues/available`);
// export const getServicePipelineByCatalogueId = id => request.get(`/service-mgmt/v1/service-pipelines/service-catalogues/${id}`);
// export const createServiceRequestTemplate = data => request.post(`/service-mgmt/v1/service-request-templates`, data);
// export const createServiceCatalogue = data => request.post(`/service-mgmt/v1/service-catalogues`, data);
// export const createServicePipeline = data => request.post(`/service-mgmt/v1/service-pipelines`, data);
// export const getAllRoles = () => request.get(`/service-mgmt/v1/core-resources/roles`);
// export const getEntityDataByTemplateId = id => request.get(`/service-mgmt/v1/service-requests/service-templates/${id}/root-entities`);
// export const getPreprocessDataByTaskId = id => request.get(`/service-mgmt/v1/tasks/${id}/preprocess`);

export const addRequestTemplateGroup = data => request.post(`/taskman/v1/request/template/group/save`, data);
export const deleteRequestTemplateGroup = id => request.delete(`/taskman/v1/request/template/group/delete/${id}`);
export const editRequestTemplateGroup = data => request.post(`/taskman/v1/request/template/group/edit`, data);
export const searchRequestTemplateGroup = data => request.post(`/taskman/v1/request/template/group/search/${data.page}/${data.pageSize}`, data.data);
export const getRoleList = () => request.get(`/taskman/v1/core-resources/roles`);
export const getProcessDefinitionKeysList = () => request.get(`/taskman/v1/core-resources/platform/process-definition-keys`);
export const getAllTemplateGroup = () => request.get(`/taskman/v1/request/template/group/available`);
export const saveRequestTemplate = data => request.post(`/taskman/v1/request/template/save`, data);
export const searchRequestTemplate = data => request.post(`/taskman/v1/request/template/search/${data.page}/${data.pageSize}`, data.data);
export const deleteRequestTemplate = id => request.delete(`/taskman/v1/request/template/delete/${id}`);
export const getTaskNodesEntitys = id => request.get(`/taskman/v1/core-resources/platform/process-definitions-nodes/${id}`);
export const saveFormTemplate = data => request.post(`/taskman/v1/form/template/save`, data);
export const saveTaskTemplate = data => request.post(`/taskman/v1/task/template/save`, data);
export const getAllDataModels = () => request.get(`/taskman/v1/core-resources/platform/models`);
export const getTargetOptions = (pkgName, entityName) =>request.get(`/taskman/v1/core-resources/platform/packages/${pkgName}/entities/${entityName}/retrieve`)
export const getFormTemplateDetail = (tempType, tempId) => request.get(`/taskman/v1/form/template/detail/${tempType}/${tempId}`);
export const releaseRequestTemplate = data => request.post(`/taskman/v1/request/template/release`, data);
export const searchRequest = data => request.post(`/taskman/v1/request/search/${data.page}/${data.pageSize}`, data.data);
export const saveRequestInfo = data => request.post(`/taskman/v1/request/save`, data);
export const requestTemplateAvailable = () => request.get(`/taskman/v1/request/template/available`);
export const getEntityDataByTemplateId = key => request.get(`/taskman/v1/core-resources/platform/${key}/root-entity`);
export const workflowProcessPrevieEntities = (procId,entityDataId) => request.get(`/taskman/v1/core-resources/platform/process/definitions/${procId}/preview/entities/${entityDataId}`);
export const searchTask = data => request.post(`/taskman/v1/task/task/search/${data.page}/${data.pageSize}`, data.data);
export const getTaskInfoDetails = data => request.post(`/taskman/v1/task/task/details`, data);
export const taskInfoReceive = data => request.post(`/taskman/v1/task/task/receive`, data);
export const taskInfoProcessing = data => request.post(`/taskman/v1/task/task/processing`, data);
export const getTaskInfoInstance = data => request.post(`/taskman/v1/task/task/instance`, data);




