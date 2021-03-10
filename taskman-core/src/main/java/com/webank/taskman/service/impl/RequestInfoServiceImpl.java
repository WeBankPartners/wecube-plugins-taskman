package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.base.QueryResponse.PageInfo;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.StatusEnum;
import com.webank.taskman.converter.EntityAttrValueConverter;
import com.webank.taskman.converter.FormItemInfoConverter;
import com.webank.taskman.converter.RequestInfoConverter;
import com.webank.taskman.domain.FormItemInfo;
import com.webank.taskman.domain.FormTemplate;
import com.webank.taskman.domain.RequestInfo;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.CreateTaskDto;
import com.webank.taskman.dto.req.RequestInfoQueryReqDto;
import com.webank.taskman.dto.resp.RequestInfoResqDto;
import com.webank.taskman.mapper.RequestInfoMapper;
import com.webank.taskman.service.*;
import com.webank.taskman.support.core.PlatformCoreServiceRestClient;
import com.webank.taskman.support.core.dto.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.*;
import java.util.stream.Collectors;
import java.util.stream.Stream;

@Service
public class RequestInfoServiceImpl extends ServiceImpl<RequestInfoMapper, RequestInfo> implements RequestInfoService {

    @Autowired
    private FormInfoService formInfoService;
    @Autowired
    private PlatformCoreServiceRestClient platformCoreServiceRestClient;
    @Autowired
    private FormItemInfoService formItemInfoService;
    @Autowired
    private RequestTemplateService requestTemplateService;

    @Autowired
    private RequestInfoMapper requestInfoMapper;
    @Autowired
    private RequestInfoConverter requestInfoConverter;
    @Autowired
    private FormItemInfoConverter formItemInfoConverter;

    @Autowired
    private FormTemplateService formTemplateService;
    @Autowired
    private EntityAttrValueConverter entityAttrValueConverter;

    @Override
    public QueryResponse<RequestInfoResqDto> selectRequestInfoPage(Integer current, Integer limit,
            RequestInfoQueryReqDto req) {
        req.queryCurrentUserRoles();
        IPage<RequestInfoResqDto> iPage = requestInfoMapper.selectRequestInfo(new Page<>(current, limit), req);
        QueryResponse<RequestInfoResqDto> queryResponse = new QueryResponse<>();
        queryResponse.setPageInfo(new PageInfo(iPage.getTotal(), iPage.getCurrent(), iPage.getSize()));
        queryResponse.setContents(iPage.getRecords());
        return queryResponse;
    }

    @Override
    public RequestInfoResqDto selectDetail(String id) {
        RequestInfo requestInfo = requestInfoMapper.selectOne(new RequestInfo().setId(id).getLambdaQueryWrapper());
        RequestInfoResqDto requestInfoResq = requestInfoConverter.toResp(requestInfo);
        requestInfoResq.setFormItemInfos(formItemInfoService.returnDetail(id));
        return requestInfoResq;
    }

    /**
     * Submit new request
     */
    @Override
    @Transactional
    public RequestInfoResqDto createNewRequestInfo(CreateTaskDto reqDto) {
        RequestInfo requestInfoEntity = buildRequestInfoEntity(reqDto);
        
        saveOrUpdate(requestInfoEntity);
        
        reqDto.setId(requestInfoEntity.getId());
        saveRequestFormInfo(reqDto);
        
        //remotely invoke platform service to create new process instance.
        DynamicWorkflowInstInfoDto dynamicWorkflowInstInfoDto = createNewRemoteWorkflowInstance(reqDto);
        
        if (dynamicWorkflowInstInfoDto == null) {
            log.error("Remotely create new process instance failed due to null response.");
            throw new TaskmanRuntimeException("Remotely create new process instance failed!");
        }
        
        if (StatusEnum.InProgress.name().equals(dynamicWorkflowInstInfoDto.getStatus())) {
            requestInfoEntity.setProcInstId(dynamicWorkflowInstInfoDto.getProcInstKey());
            requestInfoEntity.setStatus(dynamicWorkflowInstInfoDto.getStatus());
            requestInfoEntity.setUpdatedTime(new Date());
            updateById(requestInfoEntity);
        }
        return requestInfoConverter.toResp(requestInfoEntity);
    }

    public void saveRequestFormInfo(CreateTaskDto req) {
        List<FormItemInfo> items = new ArrayList<>();
        req.getEntities().stream().forEach(e -> {
            items.addAll(formItemInfoConverter.toEntityByAttrValue(e.getAttrValues()));
        });
        formInfoService.saveFormInfoAndItems(items, req.getRequestTempId(), req.getId());
    }

    @Override
    public DynamicWorkflowInstInfoDto createNewRemoteWorkflowInstance(CreateTaskDto req) {
        RequestTemplate requestTemplate = requestTemplateService.getById(req.getRequestTempId());
        String rootEntity = req.getRootEntity();
        
        ProcessDataPreviewDto processDataPreviewDto = platformCoreServiceRestClient
                .platformProcessDataPreview(requestTemplate.getProcDefId(), rootEntity);
        
        DynamicWorkflowInstCreationInfoDto creationInfoDto = new DynamicWorkflowInstCreationInfoDto();
        creationInfoDto.setProcDefId(requestTemplate.getProcDefId());
        creationInfoDto.setProcDefKey(requestTemplate.getProcDefKey());

        DynamicEntityValueDto rootEntityValue = createDynamicEntityValues(processDataPreviewDto, rootEntity);
        creationInfoDto.setRootEntityValue(rootEntityValue);
        List<DynamicTaskNodeBindInfoDto> taskNodeBindInfos = createTaskNodeBindInfos(
                processDataPreviewDto.getProcessSessionId(), req);
        creationInfoDto.setTaskNodeBindInfos(taskNodeBindInfos);
        
        DynamicWorkflowInstInfoDto dto = platformCoreServiceRestClient.createNewWorkflowInstance(creationInfoDto);
        return dto;
    }
    
    private RequestInfo buildRequestInfoEntity(CreateTaskDto reqDto){
        RequestInfo requestInfo = new RequestInfo();
        requestInfo.setEmergency(reqDto.getEmergency());
        requestInfo.setDescription(reqDto.getDescription());
        requestInfo.setName(reqDto.getName());
        requestInfo.setRootEntity(reqDto.getRootEntity());
        requestInfo.setId(reqDto.getId());
        requestInfo.setRequestTempId(reqDto.getRequestTempId());
        
        
        requestInfo.setCreatedBy(AuthenticationContextHolder.getCurrentUsername());
        requestInfo.setUpdatedBy(AuthenticationContextHolder.getCurrentUsername());
        requestInfo.setReporter(AuthenticationContextHolder.getCurrentUsername());
        requestInfo.setReportTime(new Date());
        requestInfo.setReportRole(AuthenticationContextHolder.getCurrentUserRolesToString());
        
        return requestInfo;
    }

    private DynamicEntityValueDto createDynamicEntityValues(ProcessDataPreviewDto processDataPreviewDto, String guid) {

        List<ProcessDataPreviewDto.GraphNodeDto> entityTreeNodes = processDataPreviewDto.getEntityTreeNodes();
        if (null == entityTreeNodes) {
            throw new TaskmanRuntimeException(String.format("getProcessDataPreview is error:%s", entityTreeNodes));
        }
        DynamicEntityValueDto rootEntityValue = new DynamicEntityValueDto();
        rootEntityValue.setDataId(guid);
        rootEntityValue.setOid(guid);
        rootEntityValue.setEntityDefId(guid);
        entityTreeNodes.stream().forEach(entityTreeNode -> {
            List<String> previousOids = Stream.of(rootEntityValue.getPreviousOids(), entityTreeNode.getPreviousIds())
                    .flatMap(Collection::stream).distinct().collect(Collectors.toList());
            previousOids.sort((e, s) -> e.compareTo(s));
            List<String> succeedingOids = Stream
                    .of(rootEntityValue.getSucceedingOids(), entityTreeNode.getSucceedingIds())
                    .flatMap(Collection::stream).distinct().collect(Collectors.toList());
            succeedingOids.sort((e, s) -> e.compareTo(s));
            rootEntityValue.setEntityName(entityTreeNode.getEntityName());
            rootEntityValue.setPackageName(entityTreeNode.getPackageName());
            rootEntityValue.setPreviousOids(previousOids);
            rootEntityValue.setSucceedingOids(succeedingOids);
        });
        return rootEntityValue;
    }

    private List<DynamicTaskNodeBindInfoDto> createTaskNodeBindInfos(String processSessionId, CreateTaskDto req) {
        List<DynamicTaskNodeBindInfoDto> dtoList = new ArrayList<>();
        Map<String, DynamicTaskNodeBindInfoDto> maps = new HashMap<>();
        //TODO
        FormTemplate formTemplate = formTemplateService
                .getOne(new FormTemplate().setTempId(req.getRequestTempId()).getLambdaQueryWrapper());

        List<TaskNodeDefObjectBindInfoDto> taskNodebindInfos = platformCoreServiceRestClient
                .platformProcessTasknodeBindings(processSessionId);
        taskNodebindInfos.stream().forEach(taskNode -> {
            String nodeDefId = taskNode.getNodeDefId();
            String[] entityTypeIds = taskNode.getEntityTypeId().split(":");
            String guid = taskNode.getEntityDataId();
            DynamicEntityValueDto entityValueDto = new DynamicEntityValueDto(guid, entityTypeIds[0], entityTypeIds[1],
                    guid, guid);
            addBondEntityByNodeDto(maps, nodeDefId, entityValueDto);
        });
        req.getEntities().stream().forEach(e -> {
            String nodeDefId = e.getNodeDefId();
            String dataId = e.getDataId();
            DynamicEntityValueDto entityValueDto = new DynamicEntityValueDto(dataId, e.getPackageName(),
                    e.getEntityName(), dataId, dataId);
            entityValueDto.setAttrValues(entityAttrValueConverter.toEntity(e.getAttrValues()));
            addBondEntityByNodeDto(maps, nodeDefId, entityValueDto);
        });
        dtoList = maps.entrySet().stream().map(e -> e.getValue()).collect(Collectors.toList());
        maps.clear();
        return dtoList;
    }

    private void addBondEntityByNodeDto(Map<String, DynamicTaskNodeBindInfoDto> maps, String nodeDefId,
            DynamicEntityValueDto entityValueDto) {
        DynamicTaskNodeBindInfoDto nodeDto = maps.get(nodeDefId);
        if (nodeDto == null) {
            nodeDto = new DynamicTaskNodeBindInfoDto(nodeDefId, nodeDefId);
        }
        List<DynamicEntityValueDto> boundEntityValues = nodeDto.getBoundEntityValues();
        boundEntityValues.add(entityValueDto);
        nodeDto.setBoundEntityValues(boundEntityValues.stream()
                .collect(Collectors.collectingAndThen(
                        Collectors.toCollection(
                                () -> new TreeSet<>(Comparator.comparing(DynamicEntityValueDto::getDataId))),
                        ArrayList<DynamicEntityValueDto>::new)));
        maps.put(nodeDefId, nodeDto);
    }
}
