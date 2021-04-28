package com.webank.taskman.service;


import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.PageableQueryResult;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.RequestTemplateGroupDto;

public interface RequestTemplateGroupService extends IService<RequestTemplateGroup> {


    RequestTemplateGroupDto saveTemplateGroupByReq(RequestTemplateGroupDto gropReq);

    PageableQueryResult<RequestTemplateGroupDto> selectRequestTemplateGroupPage(Integer current, Integer limit, RequestTemplateGroupDto req);

    void deleteTemplateGroupByIDService(String id);
}
