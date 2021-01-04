package com.webank.taskman.support.core;

import com.google.gson.reflect.TypeToken;
import com.webank.taskman.commons.AppProperties.ServiceTaskmanProperties;
import com.webank.taskman.converter.RoleRelationConverter;
import com.webank.taskman.support.core.dto.CoreResponse.*;
import com.webank.taskman.support.core.dto.*;
import com.webank.taskman.utils.GsonUtil;
import com.webank.taskman.utils.SpringUtils;
import org.checkerframework.checker.units.qual.A;
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

    @Autowired
    RoleRelationConverter roleRelationConverter;

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
            return  CoreServiceTestData.addRoles();
        }
        String json = template.get(asCoreUrl(GET_ALL_ROLES));
        List<RolesDataResponse> list =  roleRelationConverter.roleDTOToRolesDataResponseList(
                GsonUtil.toObject(json,new TypeToken<List<RolesDataDTO>>(){}));
        return list;
    }

    public List<RolesDataResponse> authRoleCurrentUser(String userName) {
        if(isDev()){
            return CoreServiceTestData.addRoleTestData();
        }
        String json = template.get(asCoreUrl(GET_ROLES_BY_USER_NAME, userName));
        List<RolesDataResponse> list = roleRelationConverter.roleDTOToRolesDataResponseList(
                GsonUtil.toObject(json,new TypeToken<List<RolesDataDTO>>(){}));
        return list;
    }

    public List<WorkflowDefInfoDto> platformProcessAll()
    {
        if(isDev()){
            return CoreServiceTestData.addPdfTestData();
        }
        String json = template.get(asCoreUrl(FETCH_LATEST_RELEASED_WORKFLOW_DEFS));
        List<WorkflowDefInfoDto> list = GsonUtil.toObject(json,new TypeToken<List<WorkflowDefInfoDto>>(){});
        return list;
    }

    public List<WorkflowNodeDefInfoDto> platformProcessNodes(String procDefId)
    {
        if(isDev()){
            if("rYsEQg2D2Bu".equals(procDefId)) {
                return  addTestNodeList();
            }
        }
        String json = template.get(asCoreUrl(FETCH_WORKFLOW_TASKNODE_INFOS, procDefId));
        List<WorkflowNodeDefInfoDto> list = GsonUtil.toObject(json,new TypeToken<List<WorkflowNodeDefInfoDto>>(){});
        return list;
    }

    public List<PluginPackageDataModelDto> platformProcessModels(String packageName)
    {
        List<PluginPackageDataModelDto> list = new ArrayList<>();
        String url = StringUtils.isEmpty(packageName) ? GET_MODELS_ALL_URL:GET_MODELS_BY_PACKAGE_URL;
        if(isDev()){
            list =  addAllDataModels();
            PluginPackageDataModelDto dto = null;
            if(!StringUtils.isEmpty(packageName)){
                list.clear();
                Optional<PluginPackageDataModelDto> optional = list.stream().filter(model-> model.getPackageName().equals(packageName)).findFirst();
                dto = optional.isPresent() ? optional.get():dto;
                list.add(dto);
            }
            return  list;
        }
        String json = template.get(asCoreUrl(url,packageName));
        if(!StringUtils.isEmpty(packageName)) {
            PluginPackageDataModelDto dataModelDto = GsonUtil.toObject(json, new TypeToken<PluginPackageDataModelDto>(){});
            list.add(dataModelDto);
        }else{
            list = GsonUtil.toObject(json, new TypeToken<List<PluginPackageDataModelDto>>(){});
        }
        return list;
    }

    public List<Map<String,Object>> platformProcessRootEntity(String procDefKey)
    {
        if(isDev()){
            return addRootEntityTestData();
        }
        String json = template.get(asCoreUrl(GET_ROOT_ENTITIES_BY_PROC_URL,procDefKey));
        List<Map<String,Object>> list = GsonUtil.toObject(json,new TypeToken<List<Map<String,Object>>>(){});
        return list;
    }

    public List<PluginPackageAttributeDto> platformProcessEntityAttributes(String packageName, String entity)
    {
        if(isDev()){
            List<PluginPackageAttributeDto> attributes = new ArrayList<>();
            List<PluginPackageDataModelDto> models = platformProcessModels(packageName);
            if(null != models && models.stream().findFirst().isPresent()){
                Set<PluginPackageEntityDto> entities =  models.stream().findFirst().get().getPluginPackageEntities();
                Optional<PluginPackageEntityDto> attributeDto = entities.stream().filter(e->e.getName().equals(entity)).findFirst();
                attributes = attributeDto.isPresent() ? attributeDto.get().getAttributes() : attributes;
            }
            return attributes;
        }
        String json = template.get(asCoreUrl(GET_ATTRIBUTES_BY_PACKAGE_ENTITY_URL, packageName,entity));
        List<PluginPackageAttributeDto> list = GsonUtil.toObject(json,new TypeToken<List<PluginPackageAttributeDto>>(){});
        return list;
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

    public ProcessDataPreviewDto platformProcessDataPreview(String procDefId, String guid)
    {
        if(isDev()){
            if("sjqH9YVJ2DP".equals(procDefId) && "0045_0000000100".equals(guid)){
                ProcessDataPreviewDto processDataPreviewDto = reviewEntities();
                return processDataPreviewDto;
            }
            return  null;
        }
        String json = template.get(asCoreUrl(GET_PROCESS_DATA_PREVIEW_URL, procDefId,guid));
        ProcessDataPreviewDto result = GsonUtil.toObject(json,new TypeToken<ProcessDataPreviewDto>(){});
        return result;
    }

    public List<TaskNodeDefObjectBindInfoDto> platformProcessTasknodeBindings(String processSessionId)
    {
        if(isDev()){
            return addProcessTasknodes();
        }
        String json = template.get(asCoreUrl(GET_PROCESS_INSTANCES_TASKNODE_BINDINGS_URL, processSessionId));
        List<TaskNodeDefObjectBindInfoDto> list = GsonUtil.toObject(json,new TypeToken<List<TaskNodeDefObjectBindInfoDto>>(){});
        return list;
    }


    public DynamicWorkflowInstInfoDto createNewWorkflowInstance(DynamicWorkflowInstCreationInfoDto creationInfoDto)
    {
        log.info("try to create new workflow instance with data: {}", creationInfoDto);
        LinkedHashMap result = template.postForResponse(asCoreUrl(CREATE_NEW_WORKFLOW_INSTANCE), creationInfoDto,LinkedHashMapResponse.class);
        return GsonUtil.toObject(GsonUtil.GsonString(result),new TypeToken<DynamicWorkflowInstInfoDto>(){});
    }

    private List<DynamicTaskNodeBindInfoDto> createTaskNodeBindInfos(String processSessionId)
    {
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

    private List<DynamicTaskNodeBindInfoDto> createTaskNodeBindInfos(List<TaskNodeDefObjectBindInfoDto> taskNodeDefObjectBindInfoDtos)
    {
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



}
