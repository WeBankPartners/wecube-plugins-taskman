package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.resp.TaskServiceMetaResp;


public interface FormItemTemplateService extends IService<FormItemTemplate> {


    void deleteRequestTemplateByID(String id);

    int deleteByDomain(FormItemTemplate formItemTemplate);

    TaskServiceMetaResp getTaskCreateServiceMeta(String procDefId, String nodeDefId);
}
