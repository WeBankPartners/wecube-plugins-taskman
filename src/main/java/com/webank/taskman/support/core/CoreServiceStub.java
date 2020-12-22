package com.webank.taskman.support.core;

import com.webank.taskman.commons.AppProperties.ServiceTaskmanProperties;
import com.webank.taskman.support.core.dto.CoreResponse.*;
import com.webank.taskman.support.core.dto.*;
import com.webank.taskman.utils.SpringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

import java.util.*;
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
    public static final String GET_PROCESS_DATA_PREVIEW_URL = "/platform/v1/process/definitions/{proc-def-id}/preview/entities/{entity-data-id}";
    public static final String GET_PROCESS_INSTANCES_TASKNODE_BINDINGS_URL = "/platform/v1/process/instances/tasknodes/session/{process-session-id}/tasknode-bindings";

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

    private boolean isDev(){
        return "dev".equals(SpringUtils.getActiveProfile());
    }

    public List<RolesDataResponse> authRoleAll() {
        if(isDev()){
            return CoreServiceTestData.addRoles();
        }
        return template.get(asCoreUrl(GET_ALL_ROLES), ListRolesDataResponse.class);
    }

    public List<RolesDataResponse> authRoleCurrentUser(String userName) {
        if(isDev()){
            return CoreServiceTestData.addRoleTestData();
        }
        return template.get(asCoreUrl(GET_ROLES_BY_USER_NAME, userName), ListRolesDataResponse.class);
    }

    public List<WorkflowDefInfoDto> platformProcessAll()  {
        if(isDev()){
            return CoreServiceTestData.addPdfTestData();
        }
        return template.get(asCoreUrl(FETCH_LATEST_RELEASED_WORKFLOW_DEFS), ListWorkflowDefInfoResponse.class);
    }

    public List<WorkflowNodeDefInfoDto> platformProcessNodes(String procDefId)  {
        List list = new ArrayList<>();
        if(isDev()){
            if("rYsEQg2D2Bu".equals(procDefId)) {
                return  addTestNodeList();
            }
        }
        return template.get(asCoreUrl(FETCH_WORKFLOW_TASKNODE_INFOS, procDefId), ListWorkflowNodeDefInfoResponse.class);
    }

    public List<PluginPackageDataModelDto> platformProcessModels(String packageName) {
        String url = StringUtils.isEmpty(packageName) ? GET_MODELS_ALL_URL:GET_MODELS_BY_PACKAGE_URL;
        if(isDev()){
            List<PluginPackageDataModelDto> list =  addAllDataModels();
            return  StringUtils.isEmpty(packageName)? new ArrayList(list):
                    new ArrayList(list.stream().filter(model-> model.getPackageName().equals(packageName)).collect(Collectors.toList())) ;
        }
        return template.get(asCoreUrl(url,packageName), ListPluginPackageDataModelResponse.class);
    }

    public List<Map<String,Object>> platformProcessRootEntity(String procDefKey){
        if(isDev()){
            return addRootEntityTestData();
        }
        return template.get(asCoreUrl(GET_ROOT_ENTITIES_BY_PROC_URL,procDefKey),ListMapDataResponse.class);
    }

    public List<PluginPackageAttributeDto> platformProcessEntityAttributes(String packageName, String entity) {
        if(isDev()){
            List<LinkedHashMap> models = new ArrayList(platformProcessModels(packageName));
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
        return template.get(asCoreUrl(GET_ATTRIBUTES_BY_PACKAGE_ENTITY_URL, packageName,entity), ListPluginPackageAttributeResponse.class);
    }

    public List<Object> platformProcessEntityRetrieve(String packageName, String entity, String filters) {
        if(isDev()){
            if("resource_set".equals(entity)){
                return addRetrieveEntityData();
            }
            return new ArrayList<>();
        }
        return template.get(asCoreUrl(GET_ENTITY_RETRIEVE_URL, packageName,entity), ListDataResponse.class,filters);
    }

    public ProcessDataPreviewDto platformProcessDataPreview(String procDefId, String guid){
        if(isDev()){
            if("sjqH9YVJ2DP".equals(procDefId) && "0045_0000000100".equals(guid)){
                ProcessDataPreviewDto processDataPreviewDto = reviewEntities();
                return processDataPreviewDto;
            }
            return  null;
        }
        return template.get(asCoreUrl(GET_PROCESS_DATA_PREVIEW_URL, procDefId,guid),ProcessDataPreviewResponse.class);
    }

    public List<TaskNodeDefObjectBindInfoDto> platformProcessTasknodeBindings(String processSessionId) {
        if(isDev()){
            return addProcessTasknodes();
        }
        return template.get(asCoreUrl(GET_PROCESS_INSTANCES_TASKNODE_BINDINGS_URL, processSessionId),
                ListTaskNodeDefObjectBindInfoResponse.class);
    }


    public DynamicWorkflowInstInfoDto createNewWorkflowInstance(DynamicWorkflowInstCreationInfoDto creationInfoDto) {
        log.info("try to create new workflow instance with data: {}", creationInfoDto);

        DynamicWorkflowInstInfoDto dto = template.postForResponse(CREATE_NEW_WORKFLOW_INSTANCE, creationInfoDto,DynamicWorkflowInstInfoResponse.class);

        return dto;
    }

    private List<DynamicTaskNodeBindInfoDto> createTaskNodeBindInfos(String processSessionId) {
        List<DynamicTaskNodeBindInfoDto> dtoList = new ArrayList<>();
        List<TaskNodeDefObjectBindInfoDto> taskNodeDefObjectBindInfoDtos = platformProcessTasknodeBindings(processSessionId);
        Map<String,DynamicTaskNodeBindInfoDto> maps = new HashMap<>();
        taskNodeDefObjectBindInfoDtos.stream().forEach( taskNode->{
            String nodeDefId = taskNode.getNodeDefId();
            DynamicTaskNodeBindInfoDto nodeDto = maps.get(nodeDefId);
            if(null == nodeDto){
                nodeDto = new DynamicTaskNodeBindInfoDto(nodeDefId,nodeDefId);
            }
            String[] entityTypeIds = taskNode.getEntityTypeId().split(":");
            String guid = taskNode.getEntityDataId();
            nodeDto.getBoundEntityValues().add(new DynamicEntityValueDto(guid,guid,guid,entityTypeIds[0],entityTypeIds[1]));
            maps.put(nodeDefId,nodeDto);
        });
        dtoList = maps.entrySet().stream().map(e->e.getValue()).collect(Collectors.toList());
        return dtoList;
    }

    private List<DynamicTaskNodeBindInfoDto> createTaskNodeBindInfos(List<TaskNodeDefObjectBindInfoDto> taskNodeDefObjectBindInfoDtos) {
        List<DynamicTaskNodeBindInfoDto> dtoList = new ArrayList<>();
        Map<String,DynamicTaskNodeBindInfoDto> maps = new HashMap<>();
        taskNodeDefObjectBindInfoDtos.stream().forEach( taskNode->{
            String nodeDefId = taskNode.getNodeDefId();
            DynamicTaskNodeBindInfoDto nodeDto = maps.get(nodeDefId);
            if(null == nodeDto){
                nodeDto = new DynamicTaskNodeBindInfoDto(nodeDefId,nodeDefId);
            }
            String[] entityTypeIds = taskNode.getEntityTypeId().split(":");
            String guid = taskNode.getEntityDataId();
            List<DynamicEntityValueDto> boundEntityValues = nodeDto.getBoundEntityValues();
            boundEntityValues.add(new DynamicEntityValueDto(guid,entityTypeIds[0],entityTypeIds[1],guid,guid));
            nodeDto.setBoundEntityValues(boundEntityValues.stream().collect(Collectors.collectingAndThen(
                Collectors.toCollection(() ->new TreeSet<>(Comparator.comparing(DynamicEntityValueDto :: getDataId))),ArrayList::new)));
            maps.put(nodeDefId,nodeDto);
        });
        dtoList = maps.entrySet().stream().map(e->e.getValue()).collect(Collectors.toList());
        maps.clear();
        return dtoList;
    }

    public static void main(String[] args) {
        CoreServiceStub stub = new CoreServiceStub();
        System.out.println(stub.createTaskNodeBindInfos(addProcessTasknodes()));
    }

}
