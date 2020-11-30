package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.*;


public interface RequestTemplateGroupService extends IService<RequestTemplateGroup> {


    void addTemplateGroup(RequestTemplateGroup req) throws Exception;

    void updateTemplateGroupService(TemplateGroupVO templateGroupVO) throws Exception;

    QueryResponse<TemplateGroupDTO> selectAllTemplateGroupService(Integer current, Integer limit, TemplateGroupReq req) throws Exception;

    void deleteTemplateGroupByIDService(String id) throws Exception;
}
