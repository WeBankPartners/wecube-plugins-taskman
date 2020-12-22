package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.base.PageInfo;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.StatusEnum;
import com.webank.taskman.converter.FormInfoConverter;
import com.webank.taskman.converter.FormItemInfoConverter;
import com.webank.taskman.converter.RequestInfoConverter;
import com.webank.taskman.domain.*;
import com.webank.taskman.dto.req.SaveRequestInfoReq;
import com.webank.taskman.dto.resp.FormInfoResq;
import com.webank.taskman.dto.resp.RequestInfoResq;
import com.webank.taskman.mapper.FormItemInfoMapper;
import com.webank.taskman.mapper.RequestInfoMapper;
import com.webank.taskman.service.*;
import com.webank.taskman.support.core.CoreServiceStub;
import com.webank.taskman.support.core.dto.DynamicEntityValueDto;
import com.webank.taskman.support.core.dto.DynamicWorkflowInstCreationInfoDto;
import com.webank.taskman.support.core.dto.DynamicWorkflowInstInfoDto;
import com.webank.taskman.support.core.dto.ProcessDataPreviewDto;
import com.webank.taskman.support.core.dto.ProcessDataPreviewDto.GraphNodeDto;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.Collection;
import java.util.List;
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
    RoleRelationService roleRelationService;

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
    public SaveRequestInfoReq saveRequestInfo(SaveRequestInfoReq req) {
        RequestInfo requestInfo = requestInfoConverter.reqToDomain(req);
        requestInfo.setCurrenUserName(requestInfo,requestInfo.getId());
        saveOrUpdate(requestInfo);
        List<FormItemInfo> formItemInfos = formItemInfoConverter.toEntity(req.getFormItems());
        String requestTempId = requestInfo.getRequestTempId();
        if(null != formItemInfos && formItemInfos.size() > 0 ){

            FormTemplate formTemplate = formTemplateService.getOne(
                    new FormTemplate(requestTempId,StatusEnum.DEFAULT.ordinal()+"").getLambdaQueryWrapper());
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
        return new SaveRequestInfoReq().setId(requestInfo.getId());
    }

    @Autowired
    CoreServiceStub coreServiceStub;

    @Autowired
    RequestTemplateService requestTemplateService;

    private DynamicWorkflowInstInfoDto createNewWorkflowInstance(RequestInfo requestInfo){
        RequestTemplate requestTemplate = requestTemplateService.getById(requestInfo.getRequestTempId());
        String guid = requestInfo.getRootEntity();

        DynamicWorkflowInstCreationInfoDto creationInfoDto = new DynamicWorkflowInstCreationInfoDto();
        creationInfoDto.setProcDefId(requestTemplate.getProcDefId());
        creationInfoDto.setProcDefKey(requestTemplate.getProcDefKey());

        DynamicEntityValueDto rootEntityValue = getDynamicEntityValueDto(requestTemplate, guid);
        creationInfoDto.setRootEntityValue(rootEntityValue);

        return coreServiceStub.createNewWorkflowInstance(creationInfoDto);
    }

    private DynamicEntityValueDto getDynamicEntityValueDto(RequestTemplate requestTemplate, String guid) {
        ProcessDataPreviewDto processDataPreviewDto = coreServiceStub.getProcessDataPreview(requestTemplate.getProcDefId(),guid);

        List<GraphNodeDto> entityTreeNodes = processDataPreviewDto.getEntityTreeNodes();
        if(null == entityTreeNodes){
            throw new TaskmanRuntimeException(String.format("getProcessDataPreview is error:%s",entityTreeNodes));
        }
        DynamicEntityValueDto rootEntityValue = new DynamicEntityValueDto();
        rootEntityValue.setDataId(guid);
        rootEntityValue.setOid(guid);
        rootEntityValue.setEntityDefId(guid);
        entityTreeNodes.stream().forEach(entityTreeNode->{
            List previousOids =  Stream.of(rootEntityValue.getPreviousOids(), entityTreeNode.getPreviousIds())
                    .flatMap(Collection::stream).distinct().collect(Collectors.toList());
            List succeedingOids = Stream.of(rootEntityValue.getSucceedingOids(), entityTreeNode.getSucceedingIds())
                    .flatMap(Collection::stream).distinct().collect(Collectors.toList());

            rootEntityValue.setEntityName(entityTreeNode.getEntityName());
            rootEntityValue.setPackageName(entityTreeNode.getPackageName());
            rootEntityValue.setPreviousOids(previousOids);
            rootEntityValue.setSucceedingOids(succeedingOids);
        });
        return rootEntityValue;
    }
}
