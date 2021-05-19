package com.webank.taskman.service;


import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.FormItemTemplate;

public interface FormItemTemplateService extends IService<FormItemTemplate> {


    void deleteRequestTemplateByID(String id);

    int deleteByDomain(FormItemTemplate formItemTemplate);

//    TaskServiceMetaRespDto getTaskCreateServiceMeta(String procDefId, String nodeDefId);
}
