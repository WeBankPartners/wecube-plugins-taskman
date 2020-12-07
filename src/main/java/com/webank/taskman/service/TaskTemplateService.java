package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.TaskTemplate;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.SaveTaskTemplateReq;
import com.webank.taskman.dto.req.SelectTaskTemplateRep;
import com.webank.taskman.dto.resp.TaskTemplateResp;

import java.util.List;


public interface TaskTemplateService extends IService<TaskTemplate> {

    TaskTemplate saveTaskTemplate(SaveTaskTemplateReq taskTemplateReq);

    void deleteTaskTemplateByIDService(String id);

    QueryResponse<TaskTemplateResp> selectTaskTemplate(Integer currentPage, Integer pageSize, SelectTaskTemplateRep rep);

    List<TaskTemplateResp> selectTaskTemplateAll();

    TaskTemplateResp selectTaskTemplateOne(String id);
}
