package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.TaskTemplate;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.SaveTaskTemplateReq;
import com.webank.taskman.dto.resp.TaskTemplateByRoleResp;
import com.webank.taskman.dto.resp.TaskTemplateResp;

import java.util.List;


public interface TaskTemplateService extends IService<TaskTemplate> {

    TaskTemplateResp saveTaskTemplateByReq(SaveTaskTemplateReq taskTemplateReq);

    void deleteTaskTemplateByIDService(String id);

    List<TaskTemplateResp> selectTaskTemplateAll();

    TaskTemplateResp selectTaskTemplateOne(String id);

    QueryResponse<TaskTemplateByRoleResp> selectTaskTemplateByRole(Integer page, Integer pageSize);
}
