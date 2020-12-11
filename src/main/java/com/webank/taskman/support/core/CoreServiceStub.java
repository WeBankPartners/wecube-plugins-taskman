package com.webank.taskman.support.core;

import com.webank.taskman.commons.AppProperties.ServiceTaskmanProperties;
import com.webank.taskman.constant.TaskNodeTypeEnum;
import com.webank.taskman.support.core.dto.*;
import com.webank.taskman.support.core.dto.CoreResponse.*;
import com.webank.taskman.utils.SpringUtils;
import org.hibernate.validator.internal.util.classhierarchy.Filters;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.text.MessageFormat;
import java.util.*;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.stream.Collectors;

import static com.webank.taskman.support.core.CoreServiceTestData.*;

@Service
public class CoreServiceStub {
    private static final int NOT_INCLUDE_DRAFT = 0;

    private static final String GET_ALL_ROLES = "/auth/v1/roles";
    public static final String GET_ROLES_BY_USER_NAME = "/auth/v1/users/%s/roles";

    private static final String FETCH_LATEST_RELEASED_WORKFLOW_DEFS = "/platform/v1/release/process/definitions";
    private static final String FETCH_WORKFLOW_TASKNODE_INFOS = "/platform/release/process/definitions/{%s}/tasknodes";

    private static final String CREATE_NEW_WORKFLOW_INSTANCE = "/platform/release/process/instances";

    public  static final String GET_MODELS_ALL_URL= "/platform/v1/models";

    public  static final String GET_MODELS_BY_PACKAGE_URL= "/platform/v1/packages/{package-name}/models";

    public  static final String GET_ROOT_ENTITIES_BY_PROC_URL= "/platform/v1/process/definitions/{proc-def-id}/root-entities";


    // entity.attributes
    public  static final String GET_ATTRIBUTES_BY_PACKAGE_ENTITY_URL=
            "/platform/v1/models/package/{plugin-package-name}/entity/{entity-name}/attributes";
    // entity to retrieve
    public static final String GET_ENTITY_RETRIEVE_URL =
            "/platform/v1/packages/{package-name}/entities/{entity-name}/retrieve";




    @Autowired
    private CoreRestTemplate template;

    @Autowired
    private ServiceTaskmanProperties smProperties;

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

    private String asCoreUrl(String path, Object... pathVariables) {
        if (pathVariables != null && pathVariables.length > 0) {
            String pattern = "\\{[^}]*\\}"; // 匹配{} 的表达式
            Matcher m = Pattern.compile(pattern).matcher(path);
            if(m.find()){
                for(Object param:pathVariables){
                    path = path.replaceFirst(pattern,param+"");
                }
            }else{
                path = String.format(path, pathVariables);
            }
        }
        return smProperties.getWecubeCoreAddress() + path;
    }

    // 1
    public List<WorkflowDefInfoDto> fetchLatestReleasedWorkflowDefs()  {
        if("dev".equals(SpringUtils.getActiveProfile())){
            return CoreServiceTestData.addPdfTestData();
        }
        return template.get(asCoreUrl(FETCH_LATEST_RELEASED_WORKFLOW_DEFS, NOT_INCLUDE_DRAFT), CommonResponseDto.class);
    }

    //2
    public List<WorkflowNodeDefInfoDto> fetchWorkflowTasknodeInfos(String procDefId)  {
        List<LinkedHashMap> list = new ArrayList<>();
        if("dev".equals(SpringUtils.getActiveProfile())){
            if("rYsEQg2D2Bu".equals(procDefId)) {
                list =  addTestNodeList();
            }
            List filterList = new LinkedList<>();
            for(LinkedHashMap node:list){
                if(TaskNodeTypeEnum.SUTN.getType().equals(node.get("taskCategory"))){
                    filterList.add(node);
                }
            }
            return filterList;
        }
        return template.get(asCoreUrl(FETCH_WORKFLOW_TASKNODE_INFOS, NOT_INCLUDE_DRAFT), CommonResponseDto.class);
    }

    // 3
    public Set<PluginPackageDataModelDto> allDataModels() {
        if("dev".equals(SpringUtils.getActiveProfile())){
            return  new HashSet(addAllDataModels());
        }
        return template.get(asCoreUrl(GET_MODELS_ALL_URL),GetModelsAllResponse.class);
    }

    // 4
    public Set<PluginPackageDataModelDto> getModelsByPackage(String packageName) {
        if("dev".equals(SpringUtils.getActiveProfile())){
            Set<LinkedHashMap> list =  addAllDataModels().stream().filter(model->model.get("packageName").equals(packageName)).collect(Collectors.toSet());
            return new HashSet(list);
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
            List<LinkedHashMap> models =new ArrayList(getModelsByPackage(packageName));
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



}
