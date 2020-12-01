package com.webank.taskman.service;

import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.RequestTemplateDTO;
import com.webank.taskman.dto.RequestTemplateReq;
import com.webank.taskman.dto.RequestTemplateVO;
import com.webank.taskman.dto.req.AddRequestTemplateReq;

import java.util.List;


public interface RequestTemplateService extends IService<RequestTemplate> {

    public void createRequestTemplateService(RequestTemplateVO requestTemplateVO) throws Exception;

    public void deleteRequestTemplateService(String id) throws Exception;

    public void  updateRequestTemplateService(RequestTemplateVO requestTemplateVO) throws Exception;

    QueryResponse<RequestTemplateDTO> selectAllequestTemplateService(Integer current, Integer limit, RequestTemplateReq req) throws Exception;

    List<RequestTemplateDTO> selectRequestTemplateService(RequestTemplateVO requestTemplateVO) throws Exception;
}
