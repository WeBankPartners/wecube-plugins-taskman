package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.TaskTemplate;
import com.webank.taskman.dto.req.SaveTaskTemplateReq;
import com.webank.taskman.dto.resp.TaskTemplateResp;


public interface TaskTemplateService extends IService<TaskTemplate> {

    TaskTemplateResp saveTaskTemplateByReq(SaveTaskTemplateReq taskTemplateReq);

}
