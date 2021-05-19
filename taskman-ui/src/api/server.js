import {req} from "./base";
const request = window.request ? window.request : req

export const getCmdbDate = id => request.post(`/wecmdb/ui/v2/ci-types/${id}/ci-data/query`, {dialect:{showCiHistory: false},filters:[],paging:false});
export const addRequestTemplateGroup = data => request.post(`/taskman/v1/request/template/group/save`, data);
export const deleteRequestTemplateGroup = id => request.delete(`/taskman/v1/request/template/group/delete/${id}`);
export const editRequestTemplateGroup = data => request.post(`/taskman/v1/request/template/group/edit`, data);
export const searchRequestTemplateGroup = data => request.post(`/taskman/v1/request/template/group/search/${data.page}/${data.pageSize}`, data.data);
export const getRoleList = () => request.get(`/taskman/v1/core-resources/roles`);
export const getProcessDefinitionKeysList = () => request.get(`/taskman/v1/core-resources/platform/process-definitions`);
export const getAllTemplateGroup = () => request.get(`/taskman/v1/request/template/group/available`);
export const saveRequestTemplate = data => request.post(`/taskman/v1/request/template/save`, data);
export const searchRequestTemplate = data => request.post(`/taskman/v1/request/template/search/${data.page}/${data.pageSize}`, data.data);
export const deleteRequestTemplate = id => request.delete(`/taskman/v1/request/template/delete/${id}`);
export const getTaskNodesEntitys = id => request.get(`/taskman/v1/core-resources/platform/process-definitions/${id}/nodes`);
export const saveFormTemplate = data => request.post(`/taskman/v1/form/template/save`, data);
export const saveTaskTemplate = data => request.post(`/taskman/v1/task/template/save`, data);
export const getAllDataModels = () => request.get(`/taskman/v1/core-resources/platform/models`);
export const getTargetOptions = (pkgName, entityName) =>request.get(`/taskman/v1/core-resources/platform/packages/${pkgName}/entities/${entityName}/retrieve`)
export const getFormTemplateDetail = (tempType, tempId) => request.get(`/taskman/v1/form/template/detail/${tempType}/${tempId}`);
export const releaseRequestTemplate = data => request.post(`/taskman/v1/request/template/release`, data);
export const searchRequest = data => request.post(`/taskman/v1/request/search/${data.page}/${data.pageSize}`, data.data);
export const saveRequestInfo = data => request.post(`/taskman/v1/request/save`, data);
export const requestTemplateAvailable = () => request.get(`/taskman/v1/request/template/available`);
export const getEntityDataByTemplateId = key => request.get(`/taskman/v1/core-resources/platform/process/${key}/root-entities`);
export const workflowProcessPrevieEntities = (procId,entityDataId) => request.get(`/taskman/v1/core-resources/platform/process/definitions/${procId}/preview/entities/${entityDataId}`);
export const searchTask = data => request.post(`/taskman/v1/task/search/${data.page}/${data.pageSize}`, data.data);
export const getTaskInfoDetails = data => request.post(`/taskman/v1/task/details`, data);
export const taskInfoReceive = data => request.post(`/taskman/v1/task/receive`, data);
export const taskInfoProcessing = data => request.post(`/taskman/v1/task/processing`, data);
export const getTaskInfoInstance = data => request.post(`/taskman/v1/task/instance`, data);
export const getRequestInfoDetails = id => request.get(`/taskman/v1/request/details/${id}`);



