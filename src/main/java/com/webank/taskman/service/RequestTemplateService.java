package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.req.QueryRequestTemplateReq;
import com.webank.taskman.dto.req.SaveRequestTemplateReq;
import com.webank.taskman.dto.resp.DetailRequestTemplateResq;
import com.webank.taskman.dto.resp.RequestTemplateResp;

import java.util.List;


public interface RequestTemplateService extends IService<RequestTemplate> {
    RequestTemplateResp saveRequestTemplate(SaveRequestTemplateReq saveRequestTemplateReq);

    void deleteRequestTemplateService(String id) throws TaskmanRuntimeException;

    QueryResponse<RequestTemplateResp> selectRequestTemplatePage
            (Integer current, Integer limit, QueryRequestTemplateReq req) throws TaskmanRuntimeException;

    DetailRequestTemplateResq detailRequestTemplate(String id) throws TaskmanRuntimeException;

    List<RequestTemplateResp> selectAvailableRequest(QueryRequestTemplateReq req);

   }
