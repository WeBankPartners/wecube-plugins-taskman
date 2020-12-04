package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.TaskTemplate;
import com.webank.taskman.dto.req.SaveTaskTemplateReq;


public interface TaskTemplateService extends IService<TaskTemplate> {

    TaskTemplate addOrUpdateTaskTemplate(SaveTaskTemplateReq taskTemplateReq);
}
