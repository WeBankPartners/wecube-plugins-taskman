package com.webank.taskman.service;


import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.RequestTemplateGroupDto;
import com.webank.taskman.dto.req.SaveRequestTemplateGroupReq;

public interface RequestTemplateGroupService extends IService<RequestTemplateGroup> {


    RequestTemplateGroupDto saveTemplateGroupByReq(SaveRequestTemplateGroupReq gropReq);

    QueryResponse<RequestTemplateGroupDto> selectRequestTemplateGroupPage(Integer current, Integer limit, RequestTemplateGroupDto req);

    void deleteTemplateGroupByIDService(String id);
}
