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
import com.webank.taskman.dto.req.SaveRequestInfoReq;
import com.webank.taskman.dto.req.SynthesisRequestInfoReq;
import com.webank.taskman.dto.resp.FormInfoResq;
import com.webank.taskman.dto.resp.RequestInfoResq;
import com.webank.taskman.dto.resp.SynthesisRequestInfoFormRequest;
import com.webank.taskman.dto.resp.SynthesisRequestInfoResp;
import com.webank.taskman.mapper.FormInfoMapper;
import com.webank.taskman.mapper.FormItemInfoMapper;
import com.webank.taskman.mapper.RequestInfoMapper;
import com.webank.taskman.service.FormInfoService;
import com.webank.taskman.service.FormTemplateService;
import com.webank.taskman.service.RequestInfoService;
import com.webank.taskman.service.RequestTemplateService;
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
    FormItemInfoConverter formItemInfoConverter;

    @Autowired
    FormItemInfoMapper formItemInfoMapper;

    @Autowired
    FormInfoService formInfoService;

    @Autowired
    FormTemplateService formTemplateService;

    @Autowired
    FormItemInfoServiceImpl formItemInfoService;

    @Autowired
    FormInfoConverter formInfoConverter;

    @Autowired
    FormInfoMapper formInfoMapper;

    @Autowired
    SynthesisRequestInfoRespConverter synthesisRequestInfoRespConverter;

    @Autowired
    SynthesisRequestInfoFormRequestConverter synthesisRequestInfoFormRequestConverter;

    private final static String STATUS_DONE = "Done";

    @Override
    public QueryResponse<RequestInfoResq> selectRequestInfoService(Integer current, Integer limit, SaveRequestInfoReq req) {

        IPage<RequestInfo> iPage = requestInfoMapper.selectRequestInfo(new Page<>(current, limit),req);
        List<RequestInfoResq> respList = requestInfoConverter.toDto(iPage.getRecords());

        for (RequestInfoResq requestInfoResq : respList) {
            FormInfo formInfo=formInfoService.getOne(new FormInfo().setRecordId(requestInfoResq.getId()).getLambdaQueryWrapper());
            FormInfoResq formInfoResq=formInfoConverter.toDto(formInfo);
            formInfoResq.setFormItemInfo(formItemInfoMapper.selectFormItemInfo(requestInfoResq.getId()));
            requestInfoResq.setFormInfoResq(formInfoResq);
        }

        QueryResponse<RequestInfoResq> queryResponse = new QueryResponse<>();
        queryResponse.setPageInfo(new PageInfo(iPage.getTotal(),iPage.getCurrent(),iPage.getSize()));
        queryResponse.setContents(respList);
        return queryResponse;
    }


    @Override
    @Transactional
    public RequestInfoResq saveRequestInfo(SaveRequestInfoReq req) {
        RequestInfo requestInfo = requestInfoConverter.reqToDomain(req);
        requestInfo.setCurrenUserName(requestInfo,requestInfo.getId());
        requestInfo.setReporter(AuthenticationContextHolder.getCurrentUsername());
        requestInfo.setReportTime(new Date());
        requestInfo.setReportRole(AuthenticationContextHolder.getCurrentUserRolesToString());
        saveOrUpdate(requestInfo);
        List<FormItemInfo> formItemInfos = formItemInfoConverter.toEntity(req.getFormItems());
        String requestTempId = requestInfo.getRequestTempId();
        if(null != formItemInfos && formItemInfos.size() > 0 ){

            FormTemplate formTemplate = formTemplateService.getOne(
                    new FormTemplate(null,requestTempId,StatusEnum.DEFAULT.ordinal()+"").getLambdaQueryWrapper());
            if(null == formTemplate){
                throw new TaskmanRuntimeException("The FormTemplate do not exist");
            }
            formInfoService.remove(new QueryWrapper<FormInfo>().setEntity(new FormInfo().setRecordId(requestInfo.getId())));
            FormInfo form = new FormInfo();
            form.setRecordId(requestInfo.getId());
            form.setFormTemplateId(formTemplate.getId());
            form.setCurrenUserName(form,form.getId());
            formInfoService.save(form);
            formItemInfos.stream().forEach(item -> {
                item.setFormId(form.getId());
                formItemInfoService.save(item);
            });
        }
        DynamicWorkflowInstInfoDto response = createNewWorkflowInstance(requestInfo);
        if(null == response){
            throw new TaskmanRuntimeException("Core interface:[createNewWorkflowInstance] call failed!");
        }
        if(StatusEnum.InProgress.name().equals(response.getStatus())){
            requestInfo.setProcInstKey(response.getProcInstKey());
            requestInfo.setStatus(response.getStatus());
            updateById(requestInfo);
        }
        return requestInfoConverter.toDto(requestInfo);
    }

    @Override
    public SynthesisRequestInfoFormRequest selectSynthesisRequestInfoFormService(String id) throws TaskmanRuntimeException {
        FormInfo formInfo=formInfoMapper.selectOne(new FormInfo().setRecordId(id).getLambdaQueryWrapper());
        if (null==formInfo||"".equals(formInfo)){
            throw new TaskmanRuntimeException("The request details do not exist");
        }
        List<FormItemInfo> formItemInfos=formItemInfoMapper.selectList(new FormItemInfo().setFormId(formInfo.getId()).getLambdaQueryWrapper());
        SynthesisRequestInfoFormRequest srt=synthesisRequestInfoFormRequestConverter.toDto(formInfo);
        srt.setFormItemInfo(formItemInfos);

        return srt;
    }

    @Override
    public QueryResponse<SynthesisRequestInfoResp> selectSynthesisRequestInfoService(Integer current, Integer limit, SynthesisRequestInfoReq req) throws TaskmanRuntimeException {
        String currentUserRolesToString = AuthenticationContextHolder.getCurrentUserRolesToString();
        req.setRoleName(currentUserRolesToString);
        IPage<RequestInfo> iPage = requestInfoMapper.selectSynthesisRequestInfo(new Page<>(current, limit),req);
        List<SynthesisRequestInfoResp> srt=synthesisRequestInfoRespConverter.toDto(iPage.getRecords());

        QueryResponse<SynthesisRequestInfoResp> queryResponse = new QueryResponse<>();
        queryResponse.setPageInfo(new PageInfo(iPage.getTotal(),iPage.getCurrent(),iPage.getSize()));
        queryResponse.setContents(srt);

        return queryResponse;
    }

    @Autowired
    CoreServiceStub coreServiceStub;

    @Autowired
    RequestTemplateService requestTemplateService;

    @Override
    public DynamicWorkflowInstInfoDto createNewWorkflowInstance(RequestInfo requestInfo){
        RequestTemplate requestTemplate = requestTemplateService.getById(requestInfo.getRequestTempId());
        String guid = requestInfo.getRootEntity();
        ProcessDataPreviewDto processDataPreviewDto = coreServiceStub.platformProcessDataPreview(requestTemplate.getProcDefId(),guid);
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
    public DynamicWorkflowInstCreationInfoDto createDynamicWorkflowInstCreationInfoDto(String  procDefId, String guid){
        RequestTemplate requestTemplate = requestTemplateService.getOne(new RequestTemplate().setProcDefId(procDefId).getLambdaQueryWrapper());
        ProcessDataPreviewDto processDataPreviewDto = coreServiceStub.platformProcessDataPreview(procDefId,guid);
        DynamicWorkflowInstCreationInfoDto creationInfoDto = new DynamicWorkflowInstCreationInfoDto();
        creationInfoDto.setProcDefId(requestTemplate.getProcDefId());
        creationInfoDto.setProcDefKey(requestTemplate.getProcDefKey());
        DynamicEntityValueDto rootEntityValue = createDynamicEntityValues(processDataPreviewDto, guid);
        creationInfoDto.setRootEntityValue(rootEntityValue);
        creationInfoDto.setTaskNodeBindInfos(createTaskNodeBindInfos(processDataPreviewDto.getProcessSessionId()));
        return creationInfoDto;
    }

    private DynamicEntityValueDto createDynamicEntityValues(ProcessDataPreviewDto processDataPreviewDto, String guid) {

        List<GraphNodeDto> entityTreeNodes = processDataPreviewDto.getEntityTreeNodes();
        if(null == entityTreeNodes){
            throw new TaskmanRuntimeException(String.format("getProcessDataPreview is error:%s",entityTreeNodes));
        }
        DynamicEntityValueDto rootEntityValue = new DynamicEntityValueDto();
        rootEntityValue.setDataId(guid);
        rootEntityValue.setOid(guid);
        rootEntityValue.setEntityDefId(guid);
        entityTreeNodes.stream().forEach(entityTreeNode->{
            List<String> previousOids =  Stream.of(rootEntityValue.getPreviousOids(), entityTreeNode.getPreviousIds())
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
