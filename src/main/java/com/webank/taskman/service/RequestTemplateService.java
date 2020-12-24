package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.req.QueryRequestTemplateReq;
import com.webank.taskman.dto.req.SaveRequestTemplateReq;
import com.webank.taskman.dto.resp.DetailRequestTemplateResq;
import com.webank.taskman.dto.RequestTemplateDTO;

import java.util.List;


public interface RequestTemplateService extends IService<RequestTemplate> {

    RequestTemplateDTO saveRequestTemplate(SaveRequestTemplateReq saveRequestTemplateReq);

    void deleteRequestTemplateService(String id) throws TaskmanRuntimeException;

    QueryResponse<RequestTemplateDTO> selectRequestTemplatePage
            (Integer current, Integer limit, QueryRequestTemplateReq req) throws TaskmanRuntimeException;

    DetailRequestTemplateResq detailRequestTemplate(String id) throws TaskmanRuntimeException;

    List<RequestTemplateDTO> requestTemplateAvailable(QueryRequestTemplateReq req);

    List<RequestTemplate>  selectListByParam(QueryRequestTemplateReq req);

   }
