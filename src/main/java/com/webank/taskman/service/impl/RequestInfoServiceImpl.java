package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.base.PageInfo;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.StatusEnum;
import com.webank.taskman.converter.*;
import com.webank.taskman.domain.*;
import com.webank.taskman.dto.req.QueryRequestInfoReq;
import com.webank.taskman.dto.req.SaveFormItemInfoReq;
import com.webank.taskman.dto.req.SaveRequestInfoReq;
import com.webank.taskman.dto.resp.*;
import com.webank.taskman.mapper.RequestInfoMapper;
import com.webank.taskman.service.*;
import com.webank.taskman.support.core.CoreServiceStub;
import com.webank.taskman.support.core.dto.*;
import com.webank.taskman.support.core.dto.ProcessDataPreviewDto.GraphNodeDto;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.*;
import java.util.stream.Collectors;
import java.util.stream.Stream;


@Service
public class RequestInfoServiceImpl extends ServiceImpl<RequestInfoMapper, RequestInfo> implements RequestInfoService {


    @Autowired
    RequestInfoMapper requestInfoMapper;

    @Autowired
    RequestInfoConverter requestInfoConverter;

    @Autowired
    FormInfoService formInfoService;

    @Autowired
    FormItemInfoConverter formItemInfoConverter;

    @Autowired
    FormItemInfoService formItemInfoService;



    @Override
    public QueryResponse<RequestInfoResq> selectRequestInfoService(Integer current, Integer limit, QueryRequestInfoReq req) {
        req.setEqUseRole("rt");
        IPage<RequestInfoResq> iPage = requestInfoMapper.selectRequestInfo(new Page<>(current, limit), req);
        QueryResponse<RequestInfoResq> queryResponse = new QueryResponse<>();
        queryResponse.setPageInfo(new PageInfo(iPage.getTotal(), iPage.getCurrent(), iPage.getSize()));
        queryResponse.setContents(iPage.getRecords());
        return queryResponse;
    }


    @Override
    @Transactional
    public RequestInfoResq saveRequestInfo(SaveRequestInfoReq req) {
        RequestInfo requestInfo = requestInfoConverter.reqToDomain(req);
        requestInfo.setCurrenUserName(requestInfo, requestInfo.getId());
        requestInfo.setReporter(AuthenticationContextHolder.getCurrentUsername());
        requestInfo.setReportTime(new Date());
        requestInfo.setReportRole(AuthenticationContextHolder.getCurrentUserRolesToString());
        saveOrUpdate(requestInfo);
        formInfoService.saveFormInfoAndItems(formItemInfoConverter.toEntityByReqs(req.getFormItems()), requestInfo.getRequestTempId(), requestInfo.getId());

        DynamicWorkflowInstInfoDto response = createNewWorkflowInstance(requestInfo);
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



    @Autowired
    CoreServiceStub coreServiceStub;

    @Autowired
    RequestTemplateService requestTemplateService;

    @Override
    public DynamicWorkflowInstInfoDto createNewWorkflowInstance(RequestInfo requestInfo) {
        RequestTemplate requestTemplate = requestTemplateService.getById(requestInfo.getRequestTempId());
        String guid = requestInfo.getRootEntity();
        ProcessDataPreviewDto processDataPreviewDto = coreServiceStub.platformProcessDataPreview(requestTemplate.getProcDefId(), guid);
        DynamicWorkflowInstCreationInfoDto creationInfoDto = new DynamicWorkflowInstCreationInfoDto();
        creationInfoDto.setProcDefId(requestTemplate.getProcDefId());
        creationInfoDto.setProcDefKey(requestTemplate.getProcDefKey());

        DynamicEntityValueDto rootEntityValue = createDynamicEntityValues(processDataPreviewDto, guid);
        creationInfoDto.setRootEntityValue(rootEntityValue);
        creationInfoDto.setTaskNodeBindInfos(createTaskNodeBindInfos(processDataPreviewDto.getProcessSessionId()));
        DynamicWorkflowInstInfoDto dto = coreServiceStub.createNewWorkflowInstance(creationInfoDto);
        return dto;
    }

    @Override
    public DynamicWorkflowInstCreationInfoDto createDynamicWorkflowInstCreationInfoDto(String procDefId, String guid) {
        RequestTemplate requestTemplate = requestTemplateService.getOne(new RequestTemplate().setProcDefId(procDefId).getLambdaQueryWrapper());
        ProcessDataPreviewDto processDataPreviewDto = coreServiceStub.platformProcessDataPreview(procDefId, guid);
        DynamicWorkflowInstCreationInfoDto creationInfoDto = new DynamicWorkflowInstCreationInfoDto();
        creationInfoDto.setProcDefId(requestTemplate.getProcDefId());
        creationInfoDto.setProcDefKey(requestTemplate.getProcDefKey());
        DynamicEntityValueDto rootEntityValue = createDynamicEntityValues(processDataPreviewDto, guid);
        creationInfoDto.setRootEntityValue(rootEntityValue);
        creationInfoDto.setTaskNodeBindInfos(createTaskNodeBindInfos(processDataPreviewDto.getProcessSessionId()));
        return creationInfoDto;
    }

    @Override
    public RequestInfoResq selectDetail(String id) {
        RequestInfo requestInfo = requestInfoMapper.selectOne(new RequestInfo().setId(id).getLambdaQueryWrapper());
        RequestInfoResq requestInfoResq = requestInfoConverter.toResp(requestInfo);
        requestInfoResq.setFormItemInfos(formItemInfoService.returnDetail(id));
        return requestInfoResq;
    }

    private DynamicEntityValueDto createDynamicEntityValues(ProcessDataPreviewDto processDataPreviewDto, String guid) {

        List<GraphNodeDto> entityTreeNodes = processDataPreviewDto.getEntityTreeNodes();
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

    private List<DynamicTaskNodeBindInfoDto> createTaskNodeBindInfos(String processSessionId) {
        List<DynamicTaskNodeBindInfoDto> dtoList = new ArrayList<>();
        List<TaskNodeDefObjectBindInfoDto> taskNodeDefObjectBindInfoDtos =
                coreServiceStub.platformProcessTasknodeBindings(processSessionId);
        Map<String, DynamicTaskNodeBindInfoDto> maps = new HashMap<>();
        taskNodeDefObjectBindInfoDtos.stream().forEach(taskNode -> {
            String nodeDefId = taskNode.getNodeDefId();
            DynamicTaskNodeBindInfoDto nodeDto = maps.get(nodeDefId);
            if (null == nodeDto) {
                nodeDto = new DynamicTaskNodeBindInfoDto(nodeDefId, nodeDefId);
            }
            String[] entityTypeIds = taskNode.getEntityTypeId().split(":");
            String guid = taskNode.getEntityDataId();
            List<DynamicEntityValueDto> boundEntityValues = nodeDto.getBoundEntityValues();
            boundEntityValues.add(new DynamicEntityValueDto(guid, entityTypeIds[0], entityTypeIds[1], guid, guid));
            nodeDto.setBoundEntityValues(boundEntityValues.stream().collect(Collectors.collectingAndThen(
                    Collectors.toCollection(() -> new TreeSet<>(Comparator.comparing(DynamicEntityValueDto::getDataId))), ArrayList::new)));
            maps.put(nodeDefId, nodeDto);
        });
        dtoList = maps.entrySet().stream().map(e -> e.getValue()).collect(Collectors.toList());
        maps.clear();
        return dtoList;
    }
}
