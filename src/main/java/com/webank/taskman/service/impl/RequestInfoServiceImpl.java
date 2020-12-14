package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.commons.TaskmanException;
import com.webank.taskman.constant.RoleTypeEnum;
import com.webank.taskman.constant.StatusCodeEnum;
import com.webank.taskman.converter.FormInfoConverter;
import com.webank.taskman.converter.FormItemInfoConverter;
import com.webank.taskman.converter.RequestInfoConverter;
import com.webank.taskman.domain.FormInfo;
import com.webank.taskman.domain.FormItemInfo;
import com.webank.taskman.domain.FormTemplate;
import com.webank.taskman.domain.RequestInfo;
import com.webank.taskman.dto.DoneServiceRequestRequest;
import com.webank.taskman.dto.PageInfo;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.SaveRequestInfoReq;
import com.webank.taskman.dto.resp.FormInfoResq;
import com.webank.taskman.dto.resp.RequestInfoResq;
import com.webank.taskman.mapper.FormItemInfoMapper;
import com.webank.taskman.mapper.RequestInfoMapper;
import com.webank.taskman.service.FormInfoService;
import com.webank.taskman.service.FormTemplateService;
import com.webank.taskman.service.RequestInfoService;
import com.webank.taskman.service.RoleRelationService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;


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
    public QueryResponse<RequestInfoResq> selectRequestInfoService(Integer current, Integer limit, SaveRequestInfoReq req) throws Exception {

        IPage<RequestInfo> iPage = requestInfoMapper.selectRequestInfo(new Page<>(current, limit),req);
        List<RequestInfoResq> respList = requestInfoConverter.toDto(iPage.getRecords());

        for (RequestInfoResq requestInfoResq : respList) {
            FormInfo formInfo=formInfoService.getOne(new QueryWrapper<FormInfo>().eq("record_id",requestInfoResq.getId()));
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
    public void doneServiceRequest(DoneServiceRequestRequest completedRequest) {
        RequestInfo serviceRequestResult = this.baseMapper.selectById(completedRequest.getServiceRequestId());
        if (null == serviceRequestResult) {
            String msg = String.format("Service Request [%s] not found", completedRequest.getServiceRequestId());
            throw new TaskmanException(StatusCodeEnum.NOT_FOUND_RECORD.getCode(), msg, completedRequest.getServiceRequestId());
        }
//        serviceRequestResult.setResult(completedRequest.getResult());
        serviceRequestResult.setStatus(STATUS_DONE);
        saveOrUpdate(serviceRequestResult);
    }


    @Override
    @Transactional
    public SaveRequestInfoReq saveRequestInfo(SaveRequestInfoReq req) {
        RequestInfo requestInfo = requestInfoConverter.reqToDomain(req);
        requestInfo.setCurrenUserName(requestInfo,requestInfo.getId());
        saveOrUpdate(requestInfo);
        List<FormItemInfo> formItemInfos = formItemInfoConverter.toEntity(req.getFormItems());
        if(null != formItemInfos && formItemInfos.size() > 0 ){
            FormTemplate formTemplate = formTemplateService.getOne(
                    new QueryWrapper<FormTemplate>().eq("temp_id",requestInfo.getRequestTempId()).eq("temp_type",0) );
            FormInfo form = new FormInfo();
            form.setRecordId(requestInfo.getId());
            form.setFormTemplateId(formTemplate.getId());
            form.setName(requestInfo.getName()+"_form");
            form.setCurrenUserName(form,form.getId());
            formInfoService.save(form);
            formItemInfos.stream().forEach(item -> {
                item.setFormId(form.getId());
                formItemInfoService.save(item);
            });
        }

        return new SaveRequestInfoReq().setId(requestInfo.getId());
    }
}
