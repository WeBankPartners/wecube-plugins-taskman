package com.webank.taskman.support.core;

import com.webank.taskman.commons.AppProperties.ServiceTaskmanProperties;
import com.webank.taskman.support.core.dto.CoreResponse.*;
import com.webank.taskman.support.core.dto.*;
import com.webank.taskman.utils.SpringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.stream.Collectors;

import static com.webank.taskman.support.core.CoreServiceTestData.*;

@Service
public class CoreServiceStub {


    private static final Logger log = LoggerFactory.getLogger(CoreServiceStub.class);

    public static final String GET_ALL_ROLES = "/auth/v1/roles";
    public static final String GET_ROLES_BY_USER_NAME = "/auth/v1/users/{user-name}/roles";

    public static final String CREATE_NEW_WORKFLOW_INSTANCE = "/platform/v1/release/process/instances";

    public static final String GET_MODELS_ALL_URL= "/platform/v1/models";
    public static final String GET_MODELS_BY_PACKAGE_URL= "/platform/v1/packages/{package-name}/models";
    public static final String FETCH_LATEST_RELEASED_WORKFLOW_DEFS = "/platform/v1/release/process/definitions";
    public static final String FETCH_WORKFLOW_TASKNODE_INFOS = "/platform/v1/release/process/definitions/{proc-def-id}/tasknodes";
    public static final String GET_ROOT_ENTITIES_BY_PROC_URL= "/platform/v1/process/definitions/{proc-def-id}/root-entities";
    public static final String GET_ATTRIBUTES_BY_PACKAGE_ENTITY_URL= "/platform/v1/models/package/{plugin-package-name}/entity/{entity-name}/attributes";
    public static final String GET_ENTITY_RETRIEVE_URL = "/platform/v1/packages/{package-name}/entities/{entity-name}/retrieve";
    public static final String QUERY_ENTITY_RETRIEVE_URL = "/platform/v1/packages/{package-name}/entities/{entity-name}/query";
    public static final String GET_PROCESS_DATA_PREVIEW_URL = "/platform/process/definitions/{proc-def-id}/preview/entities/{entity-data-id}";

    @Autowired
    private CoreRestTemplate template;

    @Autowired
    private ServiceTaskmanProperties smProperties;

    private String asCoreUrl(String path, Object... pathVariables) {
        log.info("URL before formatting:{}",path);
        if (null != pathVariables  && pathVariables.length > 0) {
            String pattern = "\\{(.*?)}";
            Matcher m = Pattern.compile(pattern).matcher(path);
            if(m.find()){
                for(Object param:pathVariables){
                    path = path.replaceFirst(pattern,param+"");
                }
            }else{
                path = String.format(path, pathVariables);
            }
        }
        log.info("URL after formatting:{}",path);
        return null!=smProperties? smProperties.getWecubeCoreAddress() + path :path;
    }

    //1
    public List<RolesDataResponse> getAllRoles() {
        if("dev".equals(SpringUtils.getActiveProfile())){
            return CoreServiceTestData.addRoles();
        }
        return template.get(asCoreUrl(GET_ALL_ROLES), GetAllRolesResponse.class);
    }
    // 1
    public List<RolesDataResponse> getRolesByUserName(String userName) {
        if("dev".equals(SpringUtils.getActiveProfile())){
            return CoreServiceTestData.addRoleTestData();
        }
        return template.get(asCoreUrl(GET_ROLES_BY_USER_NAME, userName), GetAllRolesResponse.class);
    }


    // 1
    public List<WorkflowDefInfoDto> fetchLatestReleasedWorkflowDefs()  {
        if("dev".equals(SpringUtils.getActiveProfile())){
            return CoreServiceTestData.addPdfTestData();
        }
        return template.get(asCoreUrl(FETCH_LATEST_RELEASED_WORKFLOW_DEFS), CommonResponseDto.class);
    }

    //2
    public List<WorkflowNodeDefInfoDto> fetchWorkflowTasknodeInfos(String procDefId)  {
        List list = new ArrayList<>();
        if("dev".equals(SpringUtils.getActiveProfile())){
            if("rYsEQg2D2Bu".equals(procDefId)) {
                list =  addTestNodeList();
                return list;
            }
        }
        return template.get(asCoreUrl(FETCH_WORKFLOW_TASKNODE_INFOS, procDefId), CommonResponseDto.class);
    }

    // 3
    public List<PluginPackageDataModelDto> allDataModels() {
        if("dev".equals(SpringUtils.getActiveProfile())){
            return  new ArrayList(addAllDataModels());
        }
        return template.get(asCoreUrl(GET_MODELS_ALL_URL),GetModelsAllResponse.class);
    }

    // 4
    public List<PluginPackageDataModelDto> getModelsByPackage(String packageName) {
        if("dev".equals(SpringUtils.getActiveProfile())){
            List list =  addAllDataModels().stream().filter(model->model.get("packageName").equals(packageName)).collect(Collectors.toList());
            return list;
        }
        return template.get(asCoreUrl(GET_MODELS_BY_PACKAGE_URL, packageName), GetModelsAllResponse.class);
    }

    public List<Map<String, Object>> getProcessDefinitionRootEntitiesByProcDefKey(String procDefKey){
        if("dev".equals(SpringUtils.getActiveProfile())){
            return addRootEntityTestData();
        }
        return template.get(asCoreUrl(GET_ROOT_ENTITIES_BY_PROC_URL,procDefKey),ListMapDataResponse.class);
    }

    // 5
    public List<PluginPackageAttributeDto> getAttributesByPackageEntity(String packageName,String entity) {
        if("dev".equals(SpringUtils.getActiveProfile())){
            List<LinkedHashMap> models = new ArrayList(getModelsByPackage(packageName));
            if(models.size() > 0){
                List<LinkedHashMap> entitys = new ArrayList((ArrayList)models.get(0).get("pluginPackageEntities"));
                if(entitys.size()>0){
                    List<LinkedHashMap> entityList =  entitys.stream().filter(e-> e.get("name").equals(entity)).collect(Collectors.toList());
                    List attributes = (ArrayList)entityList.get(0).get("attributes");
                    return attributes;
                }
            }
            return new ArrayList<>();
        }
        return template.get(asCoreUrl(GET_ATTRIBUTES_BY_PACKAGE_ENTITY_URL, packageName,entity), GetModelsAllResponse.class);
    }

    // 6
    public List<Object> retrieveEntity(String packageName, String entity, String filters) {
        if("dev".equals(SpringUtils.getActiveProfile())){
            if("resource_set".equals(entity)){
                return addRetrieveEntityData();
            }
            return new ArrayList<>();
        }

        return template.get(asCoreUrl(GET_ENTITY_RETRIEVE_URL, packageName,entity), ListDataResponse.class,filters);
    }

    //7
    public DynamicWorkflowInstInfoDto createNewWorkflowInstance(DynamicWorkflowInstCreationInfoDto creationInfoDto) {
        return template.postForResponse(CREATE_NEW_WORKFLOW_INSTANCE, creationInfoDto,DefaultCoreResponse.class);
    }

    public List<Object> rootEntities(String procDefId){
        if("dev".equals(SpringUtils.getActiveProfile())){
            if("sjqH9YVJ2DP".equals(procDefId)){
                return addRetrieveEntityData();
            }
            return new ArrayList<>();
        }
        return template.get(asCoreUrl(GET_ROOT_ENTITIES_BY_PROC_URL, procDefId),
                ListDataResponse.class);
    }

    public ProcessDataPreviewDto getProcessDataPreview(String procDefId, String guid){
        if("dev".equals(SpringUtils.getActiveProfile())){
            if("sjqH9YVJ2DP".equals(procDefId) && "0045_0000000100".equals(guid)){
                return reviewEntities();
            }
            return  null;
        }
        return template.get(asCoreUrl(GET_PROCESS_DATA_PREVIEW_URL, procDefId,guid),
                ReviewEntitiesDTOResponse.class);
    }

}
