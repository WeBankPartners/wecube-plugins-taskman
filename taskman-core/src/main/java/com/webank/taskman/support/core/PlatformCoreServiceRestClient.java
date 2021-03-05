package com.webank.taskman.support.core;

import java.util.ArrayList;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.google.gson.reflect.TypeToken;
import com.webank.taskman.commons.AppProperties.ServiceTaskmanProperties;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.support.core.dto.CallbackRequestDto;
import com.webank.taskman.support.core.dto.CoreResponse.DefaultCoreResponse;
import com.webank.taskman.support.core.dto.CoreResponse.LinkedHashMapResponse;
import com.webank.taskman.support.core.dto.CoreResponse.ListDataResponse;
import com.webank.taskman.support.core.dto.DataModelEntityDto;
import com.webank.taskman.support.core.dto.DynamicWorkflowInstCreationInfoDto;
import com.webank.taskman.support.core.dto.DynamicWorkflowInstInfoDto;
import com.webank.taskman.support.core.dto.PluginPackageAttributeDto;
import com.webank.taskman.support.core.dto.PluginPackageDataModelDto;
import com.webank.taskman.support.core.dto.ProcessDataPreviewDto;
import com.webank.taskman.support.core.dto.RolesDataResponse;
import com.webank.taskman.support.core.dto.TaskNodeDefObjectBindInfoDto;
import com.webank.taskman.support.core.dto.WorkflowDefInfoDto;
import com.webank.taskman.support.core.dto.WorkflowNodeDefInfoDto;
import com.webank.taskman.utils.GsonUtil;

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

    @Autowired
    private CoreRestTemplate template;

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
        return null != smProperties ? smProperties.getWecubeCoreAddress() + path : path;
    }

    public List<RolesDataResponse> getAllPlatformAuthRoles() {
        String json = template.get(asCoreUrl(GET_ALL_ROLES));
        List<RolesDataResponse> list = GsonUtil.toObject(json, new TypeToken<List<RolesDataResponse>>() {
        });
        return list;
    }

    public List<RolesDataResponse> getAllAuthRolesOfCurrentUser() {
        String currentUserName = AuthenticationContextHolder.getCurrentUsername();
        if(StringUtils.isBlank(currentUserName)){
            return Collections.emptyList();
        }
        
        String json = template.get(asCoreUrl(GET_ROLES_BY_USER_NAME, currentUserName));
        List<RolesDataResponse> list = GsonUtil.toObject(json, new TypeToken<List<RolesDataResponse>>() {
        });
        return list;
    }

    public List<WorkflowDefInfoDto> getAllLatestPlatformProcesses() {
        String json = template.get(asCoreUrl(FETCH_LATEST_RELEASED_WORKFLOW_DEFS));
        List<WorkflowDefInfoDto> list = GsonUtil.toObject(json, new TypeToken<List<WorkflowDefInfoDto>>() {
        });
        return list;
    }

    public List<WorkflowNodeDefInfoDto> getPlatformProcessDefinitionNodes(String procDefId) {
        String json = template.get(asCoreUrl(FETCH_WORKFLOW_TASKNODE_INFOS, procDefId));
        List<WorkflowNodeDefInfoDto> list = GsonUtil.toObject(json, new TypeToken<List<WorkflowNodeDefInfoDto>>() {
        });
        return list;
    }

    public List<PluginPackageDataModelDto> getAllPlatformProcessModels(String packageName) {
        List<PluginPackageDataModelDto> list = new ArrayList<>();
        String url = StringUtils.isEmpty(packageName) ? GET_MODELS_ALL_URL : GET_MODELS_BY_PACKAGE_URL;
        String json = template.get(asCoreUrl(url, packageName));
        if (!StringUtils.isEmpty(packageName)) {
            PluginPackageDataModelDto dataModelDto = GsonUtil.toObject(json,
                    new TypeToken<PluginPackageDataModelDto>() {
                    });
            list.add(dataModelDto);
        } else {
            list = GsonUtil.toObject(json, new TypeToken<List<PluginPackageDataModelDto>>() {
            });
        }
        return list;
    }

    public List<Map<String, Object>> getPlatformProcessRootEntities(String procDefKey) {
        String json = template.get(asCoreUrl(GET_ROOT_ENTITIES_BY_PROC_URL, procDefKey));
        List<Map<String, Object>> list = GsonUtil.toObject(json, new TypeToken<List<Map<String, Object>>>() {
        });
        return list;
    }

    public DataModelEntityDto getEntityByPackageNameAndName(String packageName, String entity) {
        String json = template.get(asCoreUrl(GET_ENTITY_BY_PACKAGE_NAME_AND_ENTITY_NAME_URL, packageName, entity));
        DataModelEntityDto dto = GsonUtil.toObject(json, new TypeToken<DataModelEntityDto>() {
        });
        return dto;
    }

    public List<PluginPackageAttributeDto> platformProcessEntityAttributes(String packageName, String entity) {
        String json = template.get(asCoreUrl(GET_ATTRIBUTES_BY_PACKAGE_ENTITY_URL, packageName, entity));
        List<PluginPackageAttributeDto> list = GsonUtil.toObject(json,
                new TypeToken<List<PluginPackageAttributeDto>>() {
                });
        return list;
    }

    public List<Object> platformProcessEntityRetrieve(String packageName, String entity, String filters) {
        List list = template.get(asCoreUrl(GET_ENTITY_RETRIEVE_URL, packageName, entity), ListDataResponse.class,
                filters);
        return list;
    }

    public ProcessDataPreviewDto platformProcessDataPreview(String procDefId, String guid) {
        String json = template.get(asCoreUrl(GET_PROCESS_DATA_PREVIEW_URL, procDefId, guid));
        ProcessDataPreviewDto result = GsonUtil.toObject(json, new TypeToken<ProcessDataPreviewDto>() {
        });
        return result;
    }

    public List<TaskNodeDefObjectBindInfoDto> platformProcessTasknodeBindings(String processSessionId) {
        String json = template.get(asCoreUrl(GET_PROCESS_INSTANCES_TASKNODE_BINDINGS_URL, processSessionId));
        List<TaskNodeDefObjectBindInfoDto> list = GsonUtil.toObject(json,
                new TypeToken<List<TaskNodeDefObjectBindInfoDto>>() {
                });
        return list;
    }

    public DynamicWorkflowInstInfoDto createNewWorkflowInstance(DynamicWorkflowInstCreationInfoDto creationInfoDto) {
        LinkedHashMap result = template.postForResponse(asCoreUrl(CREATE_NEW_WORKFLOW_INSTANCE), creationInfoDto,
                LinkedHashMapResponse.class);
        return GsonUtil.toObject(GsonUtil.GsonString(result), new TypeToken<DynamicWorkflowInstInfoDto>() {
        });
    }

    public Object callback(String callbackUrl, CallbackRequestDto callbackRequest) {
        return template.postForResponse(asCoreUrl(callbackUrl), callbackRequest, DefaultCoreResponse.class);
    }

}
