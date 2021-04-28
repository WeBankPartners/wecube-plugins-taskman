package com.webank.taskman.service;

import java.util.List;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.LocalPageableQueryResult;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.RequestTemplateDto;
import com.webank.taskman.dto.req.RequestTemplateQueryDto;
import com.webank.taskman.dto.resp.RequestTemplateQueryResultDto;

public interface RequestTemplateService extends IService<RequestTemplate> {

    RequestTemplateDto saveRequestTemplate(RequestTemplateDto requestTemplateDto);

    void deleteRequestTemplate(String id);

    LocalPageableQueryResult<RequestTemplateDto> searchRequestTemplates(Integer current, Integer limit,
            RequestTemplateQueryDto req);

    RequestTemplateQueryResultDto fetchRequestTemplateDetail(String id);
    
    RequestTemplateDto releaseRequestTemplate(RequestTemplateDto requestTemplateDto);
    
    List<RequestTemplateDto> fetchAvailableRequestTemplates();

}
