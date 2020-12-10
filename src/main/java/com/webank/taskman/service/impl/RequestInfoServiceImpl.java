package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.commons.TaskmanException;
import com.webank.taskman.constant.StatusCodeEnum;
import com.webank.taskman.converter.FormInfoConverter;
import com.webank.taskman.converter.RequestInfoConverter;
import com.webank.taskman.domain.FormInfo;
import com.webank.taskman.domain.RequestInfo;
import com.webank.taskman.dto.DoneServiceRequestRequest;
import com.webank.taskman.dto.PageInfo;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.SaveRequestInfoReq;
import com.webank.taskman.dto.resp.FormInfoResq;
import com.webank.taskman.dto.resp.RequestInfoResq;
import com.webank.taskman.mapper.FormInfoMapper;
import com.webank.taskman.mapper.FormItemInfoMapper;
import com.webank.taskman.mapper.RequestInfoMapper;
import com.webank.taskman.service.RequestInfoService;
import org.checkerframework.checker.units.qual.A;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;


@Service
public class RequestInfoServiceImpl extends ServiceImpl<RequestInfoMapper, RequestInfo> implements RequestInfoService {

    @Autowired
    RequestInfoMapper requestInfoMapper;

    @Autowired
    RequestInfoConverter requestInfoConverter;

    @Autowired
    FormItemInfoMapper formItemInfoMapper;

    @Autowired
    FormInfoMapper formInfoMapper;

    @Autowired
    FormInfoConverter formInfoConverter;

    private final static String STATUS_DONE = "Done";


    @Override
    public QueryResponse<RequestInfoResq> selectRequestInfoService(Integer current, Integer limit, SaveRequestInfoReq req) throws Exception {
        IPage<RequestInfo> iPage = requestInfoMapper.selectRequestInfo(new Page<>(current, limit),req);
        List<RequestInfoResq> respList = requestInfoConverter.toDto(iPage.getRecords());

        for (RequestInfoResq requestInfoResq : respList) {
            FormInfo formInfo=formInfoMapper.selectOne(new QueryWrapper<FormInfo>().eq("record_id",requestInfoResq.getId()));
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
}
