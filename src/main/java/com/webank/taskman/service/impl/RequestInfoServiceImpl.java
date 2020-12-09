package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.converter.FormInfoConverter;
import com.webank.taskman.converter.RequestInfoConverter;
import com.webank.taskman.domain.FormInfo;
import com.webank.taskman.domain.RequestInfo;
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

    @Override
    public QueryResponse<RequestInfoResq> selectRequestInfoService(Integer current, Integer limit, SaveRequestInfoReq req) throws Exception {
        IPage<RequestInfo> iPage = requestInfoMapper.selectRequestInfo(new Page<>(current, limit),req);
        List<RequestInfoResq> respList = requestInfoConverter.toDto(iPage.getRecords());

        for (RequestInfoResq requestInfoResq : respList) {
            FormInfo formInfo=formInfoMapper.selectOne(new QueryWrapper<FormInfo>().eq("record_id",requestInfoResq.getRequestTempId()));
            FormInfoResq formInfoResq=formInfoConverter.toDto(formInfo);
            formInfoResq.setFormItemInfo(formItemInfoMapper.selectFormItemInfo(requestInfoResq.getRequestTempId()));
            requestInfoResq.setFormInfoResq(formInfoResq);
        }

        QueryResponse<RequestInfoResq> queryResponse = new QueryResponse<>();
        queryResponse.setPageInfo(new PageInfo(iPage.getTotal(),iPage.getCurrent(),iPage.getSize()));
        queryResponse.setContents(respList);
        return queryResponse;
    }
}
