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
import com.webank.taskman.dto.req.QueryRequestInfoReq;
import com.webank.taskman.dto.resp.RequestInfoResq;
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
    FormInfoService formInfoService;
    @Autowired
    PlatformCoreServiceRestClient coreServiceStub;
    @Autowired
    FormItemInfoService formItemInfoService;
    @Autowired
    RequestTemplateService requestTemplateService;

    @Autowired
    RequestInfoMapper requestInfoMapper;
    @Autowired
    RequestInfoConverter requestInfoConverter;
    @Autowired
    FormItemInfoConverter formItemInfoConverter;


    @Override
    public QueryResponse<RequestInfoResq> selectRequestInfoPage(Integer current, Integer limit, QueryRequestInfoReq req) {
        req.queryCurrentUserRoles();
        IPage<RequestInfoResq> iPage = requestInfoMapper.selectRequestInfo(new Page<>(current, limit), req);
        QueryResponse<RequestInfoResq> queryResponse = new QueryResponse<>();
        queryResponse.setPageInfo(new PageInfo(iPage.getTotal(), iPage.getCurrent(), iPage.getSize()));
        queryResponse.setContents(iPage.getRecords());
        return queryResponse;
    }

    @Override
    public RequestInfoResq selectDetail(String id) {
        RequestInfo requestInfo = requestInfoMapper.selectOne(new RequestInfo().setId(id).getLambdaQueryWrapper());
        RequestInfoResq requestInfoResq = requestInfoConverter.toResp(requestInfo);
        requestInfoResq.setFormItemInfos(formItemInfoService.returnDetail(id));
        return requestInfoResq;
    }

    @Override
    @Transactional
    public RequestInfoResq saveRequestInfoByDto(CreateTaskDto req) {
        RequestInfo requestInfo = requestInfoConverter.createDtoToDomain(req);
        requestInfo.setCurrenUserName(requestInfo, requestInfo.getId());
        requestInfo.setReporter(AuthenticationContextHolder.getCurrentUsername());
        requestInfo.setReportTime(new Date());
        requestInfo.setReportRole(AuthenticationContextHolder.getCurrentUserRolesToString());
        saveOrUpdate(requestInfo);
        req.setId(requestInfo.getId());
        saveRequestFormInfo(req);
        DynamicWorkflowInstInfoDto response = createNewWorkflowInstance(req);
        if (null == response) {
            throw new TaskmanRuntimeException("Core interface:[createNewWorkflowInstance] call failed!");
        }
        if (StatusEnum.InProgress.name().equals(response.getStatus())) {
            requestInfo.setProcInstId(response.getProcInstKey());
            requestInfo.setStatus(response.getStatus());
            updateById(requestInfo);
        }
        return requestInfoConverter.toResp(requestInfo);
    }

    public void saveRequestFormInfo(CreateTaskDto req) {
        List<FormItemInfo> items = new ArrayList<>();
        req.getEntitys().stream().forEach(e->{
            items.addAll(formItemInfoConverter.toEntityByAttrValue(e.getAttrValues()));
        });
        formInfoService.saveFormInfoAndItems(items, req.getRequestTempId(), req.getId());
    }



    @Override
    public DynamicWorkflowInstInfoDto createNewWorkflowInstance(CreateTaskDto req) {
        RequestTemplate requestTemplate = requestTemplateService.getById(req.getRequestTempId());
        String rootEntity = req.getRootEntity();
        ProcessDataPreviewDto processDataPreviewDto = coreServiceStub.platformProcessDataPreview(requestTemplate.getProcDefId(), rootEntity);
        DynamicWorkflowInstCreationInfoDto creationInfoDto = new DynamicWorkflowInstCreationInfoDto();
        creationInfoDto.setProcDefId(requestTemplate.getProcDefId());
        creationInfoDto.setProcDefKey(requestTemplate.getProcDefKey());

        DynamicEntityValueDto rootEntityValue = createDynamicEntityValues(processDataPreviewDto, rootEntity);
        creationInfoDto.setRootEntityValue(rootEntityValue);
        List<DynamicTaskNodeBindInfoDto>  taskNodeBindInfos = createTaskNodeBindInfos(processDataPreviewDto.getProcessSessionId(),req);
        creationInfoDto.setTaskNodeBindInfos(taskNodeBindInfos);
        DynamicWorkflowInstInfoDto dto = coreServiceStub.createNewWorkflowInstance(creationInfoDto);
        return dto;
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
            List<String> succeedingOids = Stream.of(rootEntityValue.getSucceedingOids(), entityTreeNode.getSucceedingIds())
                    .flatMap(Collection::stream).distinct().collect(Collectors.toList());
            succeedingOids.sort((e, s) -> e.compareTo(s));
            rootEntityValue.setEntityName(entityTreeNode.getEntityName());
            rootEntityValue.setPackageName(entityTreeNode.getPackageName());
            rootEntityValue.setPreviousOids(previousOids);
            rootEntityValue.setSucceedingOids(succeedingOids);
        });
        return rootEntityValue;
    }

    @Autowired
    FormTemplateService formTemplateService;
    @Autowired
    EntityAttrValueConverter entityAttrValueConverter;

    private List<DynamicTaskNodeBindInfoDto> createTaskNodeBindInfos(String processSessionId,CreateTaskDto req) {
        List<DynamicTaskNodeBindInfoDto> dtoList = new ArrayList<>();
        Map<String, DynamicTaskNodeBindInfoDto> maps = new HashMap<>();
        FormTemplate formTemplate = formTemplateService.getOne(new FormTemplate().setTempId(req.getRequestTempId()).getLambdaQueryWrapper());

        List<TaskNodeDefObjectBindInfoDto> taskNodebindInfos = coreServiceStub.platformProcessTasknodeBindings(processSessionId);
        taskNodebindInfos.stream().forEach(taskNode -> {
            String nodeDefId = taskNode.getNodeDefId();
            String[] entityTypeIds = taskNode.getEntityTypeId().split(":");
            String guid = taskNode.getEntityDataId();
            DynamicEntityValueDto entityValueDto = new DynamicEntityValueDto(guid, entityTypeIds[0], entityTypeIds[1], guid, guid);
            addBondEntityByNodeDto(maps, nodeDefId, entityValueDto);
        });
        req.getEntitys().stream().forEach(e->{
            String nodeDefId = e.getNodeDefId();
            String dataId = e.getDataId();
            DynamicEntityValueDto entityValueDto = new DynamicEntityValueDto(dataId, e.getPackageName(), e.getEntityName(), dataId, dataId);
            entityValueDto.setAttrValues(entityAttrValueConverter.toEntity(e.getAttrValues()));
            addBondEntityByNodeDto(maps, nodeDefId, entityValueDto);
        });
        dtoList = maps.entrySet().stream().map(e -> e.getValue()).collect(Collectors.toList());
        maps.clear();
        return dtoList;
    }

    private void addBondEntityByNodeDto(Map<String, DynamicTaskNodeBindInfoDto> maps, String nodeDefId, DynamicEntityValueDto entityValueDto) {
        DynamicTaskNodeBindInfoDto nodeDto = maps.get(nodeDefId);
        if (null == nodeDto) {
            nodeDto = new DynamicTaskNodeBindInfoDto(nodeDefId, nodeDefId);
        }
        List<DynamicEntityValueDto> boundEntityValues = nodeDto.getBoundEntityValues();
        boundEntityValues.add(entityValueDto);
        nodeDto.setBoundEntityValues(boundEntityValues.stream().collect(Collectors.collectingAndThen(
                Collectors.toCollection(() -> new TreeSet<>(Comparator.comparing(DynamicEntityValueDto::getDataId))), ArrayList::new)));
        maps.put(nodeDefId, nodeDto);
    }
}
