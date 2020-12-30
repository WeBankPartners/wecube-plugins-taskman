package com.webank.taskman.service;


import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.TaskmanException;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.RequestTemplateGroupDTO;
import com.webank.taskman.dto.req.SaveRequestTemplateGropReq;

public interface RequestTemplateGroupService extends IService<RequestTemplateGroup> {


    RequestTemplateGroupDTO saveTemplateGroupByReq(SaveRequestTemplateGropReq gropReq) throws TaskmanException;

    QueryResponse<RequestTemplateGroupDTO> selectRequestTemplateGroupPage(Integer current, Integer limit, RequestTemplateGroupDTO req) throws TaskmanRuntimeException;

    void deleteTemplateGroupByIDService(String id) throws TaskmanRuntimeException;
}
