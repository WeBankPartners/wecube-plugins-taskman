package com.webank.taskman.service;


import java.util.List;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.PageableQueryResult;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.RequestTemplateGroupDto;

public interface RequestTemplateGroupService extends IService<RequestTemplateGroup> {

    List<RequestTemplateGroupDto> fetchAvailableGroupTemplates();

    RequestTemplateGroupDto saveOrUpdateTemplateGroup(RequestTemplateGroupDto requestTemplateGroupDto);

    PageableQueryResult<RequestTemplateGroupDto> searchRequestTemplateGroups(Integer current, Integer limit, RequestTemplateGroupDto req);

    void deleteRequestTemplateGroup(String id);
}
