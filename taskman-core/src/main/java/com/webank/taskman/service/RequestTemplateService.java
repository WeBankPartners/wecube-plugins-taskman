package com.webank.taskman.service;

import java.util.List;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.PageableQueryResult;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.RequestTemplateDto;
import com.webank.taskman.dto.req.RequestTemplateQueryReqDto;
import com.webank.taskman.dto.req.RequestTemplateSaveReqDto;
import com.webank.taskman.dto.resp.RequestTemplateRespDto;

public interface RequestTemplateService extends IService<RequestTemplate> {

    RequestTemplateDto saveRequestTemplate(RequestTemplateSaveReqDto saveRequestTemplateReq);

    void deleteRequestTemplateService(String id);

    PageableQueryResult<RequestTemplateDto> selectRequestTemplatePage(Integer current, Integer limit,
            RequestTemplateQueryReqDto req);

    RequestTemplateRespDto detailRequestTemplate(String id);
    
    RequestTemplateDto releaseRequestTemplate(RequestTemplateSaveReqDto req);
    
    List<RequestTemplateDto> fetchAvailableRequestTemplates();

}
