package com.webank.taskman.support.platform;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.webank.taskman.commons.AppProperties.ServiceTaskmanProperties;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.support.platform.dto.AuthRoleResponseDto;
import com.webank.taskman.support.platform.dto.DataModelEntityDto;
import com.webank.taskman.support.platform.dto.GenericPlatformResponseDto;
import com.webank.taskman.support.platform.dto.PluginPackageAttributeDto;
import com.webank.taskman.support.platform.dto.PluginPackageDataModelDto;
import com.webank.taskman.support.platform.dto.PluginPackageDataModelListResponseDto;
import com.webank.taskman.support.platform.dto.PluginPackageDataModelResponseDto;
import com.webank.taskman.support.platform.dto.ProcessDataPreviewDto;
import com.webank.taskman.support.platform.dto.SimpleLocalRoleDto;
import com.webank.taskman.support.platform.dto.TaskNodeDefObjectBindInfoDto;
import com.webank.taskman.support.platform.dto.WorkflowDefInfoDto;
import com.webank.taskman.support.platform.dto.WorkflowDefInfoResponseDto;
import com.webank.taskman.support.platform.dto.WorkflowNodeDefInfoDto;
import com.webank.taskman.support.platform.dto.WorkflowNodeDefInfoResponseDto;

@Service
public class PlatformCoreServiceRestClient {

    private static final Logger log = LoggerFactory.getLogger(PlatformCoreServiceRestClient.class);

    public static final String GET_ALL_ROLES = "/auth/v1/roles";
    public static final String GET_ROLES_BY_USER_NAME = "/auth/v1/users/{user-name}/roles";

    public static final String CREATE_NEW_WORKFLOW_INSTANCE = "/platform/v1/public/process/instances";
    public static final String FETCH_LATEST_RELEASED_WORKFLOW_DEFS = "/platform/v1/public/process/definitions";
    public static final String FETCH_WORKFLOW_TASKNODE_INFOS = "/platform/v1/public/process/definitions/{proc-def-id}/tasknodes";

    public static final String GET_MODELS_ALL_URL = "/platform/v1/models";
    public static final String GET_MODELS_BY_PACKAGE_URL = "/platform/v1/packages/{package-name}/models";

    public static final String GET_ROOT_ENTITIES_BY_PROC_URL = "/platform/v1/process/definitions/{proc-def-id}/root-entities";
    public static final String GET_ENTITY_BY_PACKAGE_NAME_AND_ENTITY_NAME_URL = "/platform/v1/models/package/{plugin-package-name}/entity/{entity-name}";
    public static final String GET_ATTRIBUTES_BY_PACKAGE_ENTITY_URL = "/platform/v1/models/package/{plugin-package-name}/entity/{entity-name}/attributes";
    public static final String GET_ENTITY_RETRIEVE_URL = "/platform/v1/packages/{package-name}/entities/{entity-name}/retrieve";
    public static final String QUERY_ENTITY_RETRIEVE_URL = "/platform/v1/packages/{package-name}/entities/{entity-name}/query";
    public static final String GET_PROCESS_DATA_PREVIEW_URL = "/platform/v1/public/process/definitions/{proc-def-id}/preview/entities/{entity-data-id}";
    public static final String GET_PROCESS_INSTANCES_TASKNODE_BINDINGS_URL = "/platform/v1/process/instances/tasknodes/session/{process-session-id}/tasknode-bindings";

    private UserJwtSsoTokenRestTemplate restTemplate = new UserJwtSsoTokenRestTemplate();

    @Autowired
    private ServiceTaskmanProperties smProperties;

    private String asCoreUrl(String path, Object... pathVariables) {
        log.info("URL before formatting:{}", path);
        if (null != pathVariables && pathVariables.length > 0) {
            String pattern = "\\{(.*?)}";
            Matcher m = Pattern.compile(pattern).matcher(path);
            if (m.find()) {
                for (Object param : pathVariables) {
                    path = path.replaceFirst(pattern, param + "");
                }
            } else {
                path = String.format(path, pathVariables);
            }
        }
        log.info("URL after formatting:{}", path);
        return smProperties != null ? smProperties.getWecubeCoreAddress() + path : path;
    }

    public List<SimpleLocalRoleDto> getAllPlatformAuthRoles() {
        String url = asCoreUrl(GET_ALL_ROLES);

        AuthRoleResponseDto responseDto = restTemplate.getForObject(url, AuthRoleResponseDto.class);
        List<SimpleLocalRoleDto> roleDtos = responseDto.getData();
        return roleDtos;
    }

    public List<SimpleLocalRoleDto> getAllAuthRolesOfCurrentUser() {
        String currentUserName = AuthenticationContextHolder.getCurrentUsername();
        if (StringUtils.isBlank(currentUserName)) {
            return Collections.emptyList();
        }

        String url = asCoreUrl(GET_ROLES_BY_USER_NAME, currentUserName);

        AuthRoleResponseDto responseDto = restTemplate.getForObject(url, AuthRoleResponseDto.class);
        List<SimpleLocalRoleDto> roleDtos = responseDto.getData();
        return roleDtos;
    }

    public List<WorkflowDefInfoDto> getAllLatestPlatformProcesses() {
        String url = asCoreUrl(FETCH_LATEST_RELEASED_WORKFLOW_DEFS);
        WorkflowDefInfoResponseDto responseDto = restTemplate.getForObject(url, WorkflowDefInfoResponseDto.class);
        return responseDto.getData();
    }

    public List<WorkflowNodeDefInfoDto> getPlatformProcessDefinitionNodes(String procDefId) {
        String url = asCoreUrl(FETCH_WORKFLOW_TASKNODE_INFOS, procDefId);
        WorkflowNodeDefInfoResponseDto responseDto = restTemplate.getForObject(url,
                WorkflowNodeDefInfoResponseDto.class);
        return responseDto.getData();
    }

    public List<PluginPackageDataModelDto> getAllPlatformProcessModels() {
        String url = GET_MODELS_ALL_URL;

        PluginPackageDataModelListResponseDto respDto = restTemplate.getForObject(url,
                PluginPackageDataModelListResponseDto.class);

        return respDto.getData();
    }

    public List<PluginPackageDataModelDto> getAllPlatformProcessModels(String packageName) {
        List<PluginPackageDataModelDto> list = new ArrayList<>();
        String url = asCoreUrl(GET_MODELS_BY_PACKAGE_URL, packageName);

        PluginPackageDataModelResponseDto responseDto = restTemplate.getForObject(url,
                PluginPackageDataModelResponseDto.class);
        list.add(responseDto.getData());
        return list;
    }

    public List<Map<String, Object>> getPlatformProcessRootEntities(String procDefId) {

        String url = asCoreUrl(GET_ROOT_ENTITIES_BY_PROC_URL, procDefId);
        ProcDefRootEntitiesResponseDto respDto = restTemplate.getForObject(url, ProcDefRootEntitiesResponseDto.class);
        return respDto.getData();
    }

    public DataModelEntityDto getEntityByPackageNameAndName(String packageName, String entity) {

        String url = asCoreUrl(GET_ENTITY_BY_PACKAGE_NAME_AND_ENTITY_NAME_URL, packageName, entity);
        DataModelEntityResponseDto respDto = restTemplate.getForObject(url, DataModelEntityResponseDto.class);
        return respDto.getData();
    }

    public List<PluginPackageAttributeDto> platformProcessEntityAttributes(String packageName, String entity) {

        String url = asCoreUrl(GET_ATTRIBUTES_BY_PACKAGE_ENTITY_URL, packageName, entity);
        PluginPackageAttributeListResponseDto respDto = restTemplate.getForObject(url,
                PluginPackageAttributeListResponseDto.class);
        return respDto.getData();
    }

    public List<Object> platformProcessEntityRetrieve(String packageName, String entity, String filters) {
        String url = asCoreUrl(GET_ENTITY_RETRIEVE_URL, packageName, entity);
        ObjectListResponseDto respDto = restTemplate.getForObject(url, ObjectListResponseDto.class);
        return respDto.getData();
    }

    public ProcessDataPreviewDto platformProcessDataPreview(String procDefId, String guid) {

        String url = asCoreUrl(GET_PROCESS_DATA_PREVIEW_URL, procDefId, guid);
        ProcessDataPreviewResponseDto respDto = restTemplate.getForObject(url, ProcessDataPreviewResponseDto.class);
        return respDto.getData();
    }

    public List<TaskNodeDefObjectBindInfoDto> platformProcessTasknodeBindings(String processSessionId) {
        String url = asCoreUrl(GET_PROCESS_INSTANCES_TASKNODE_BINDINGS_URL, processSessionId);

        TaskNodeDefObjectBindInfoListResponseDto respDto = restTemplate.getForObject(url,
                TaskNodeDefObjectBindInfoListResponseDto.class);
        return respDto.getData();
    }

//    public DynamicWorkflowInstInfoDto createNewWorkflowInstance(DynamicWorkflowInstCreationInfoDto creationInfoDto) {
//        // TODO
//        LinkedHashMap result = template.postForResponse(asCoreUrl(CREATE_NEW_WORKFLOW_INSTANCE), creationInfoDto,
//                LinkedHashMapResponse.class);
//        return GsonUtil.toObject(GsonUtil.GsonString(result), new TypeToken<DynamicWorkflowInstInfoDto>() {
//        });
//    }

    // public Object callback(String callbackUrl, CallbackRequestDto
    // callbackRequest) {
    // return template.postForResponse(asCoreUrl(callbackUrl), callbackRequest,
    // DefaultCoreResponse.class);
    // }

    public static class ProcDefRootEntitiesResponseDto extends GenericPlatformResponseDto<List<Map<String, Object>>> {

    }

    public static class DataModelEntityResponseDto extends GenericPlatformResponseDto<DataModelEntityDto> {

    }

    public static class PluginPackageAttributeListResponseDto
            extends GenericPlatformResponseDto<List<PluginPackageAttributeDto>> {

    }

    public static class ObjectListResponseDto extends GenericPlatformResponseDto<List<Object>> {

    }

    public static class ProcessDataPreviewResponseDto extends GenericPlatformResponseDto<ProcessDataPreviewDto> {
    }

    public static class TaskNodeDefObjectBindInfoListResponseDto
            extends GenericPlatformResponseDto<List<TaskNodeDefObjectBindInfoDto>> {
    }
}
