package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.SaveRequestTemplateReq;
import com.webank.taskman.dto.resp.RequestTemplateResp;

import java.util.List;


public interface RequestTemplateService extends IService<RequestTemplate> {
    RequestTemplateResp saveRequestTemplate(SaveRequestTemplateReq saveRequestTemplateReq);

    void deleteRequestTemplateService(String id) throws Exception;

    QueryResponse<RequestTemplateResp> selectAllequestTemplateService(Integer current, Integer limit, SaveRequestTemplateReq req) throws Exception;

    RequestTemplateResp detailRequestTemplate(String id) throws Exception;

   }
