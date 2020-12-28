package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.resp.TaskServiceMetaResp;
import com.webank.taskman.mapper.FormItemTemplateMapper;
import com.webank.taskman.service.FormItemTemplateService;
import org.springframework.stereotype.Service;



@Service
public class FormItemTemplateServiceImpl extends ServiceImpl<FormItemTemplateMapper, FormItemTemplate> implements FormItemTemplateService {


    @Override
    public void deleteRequestTemplateByID(String id) {
        this.getBaseMapper().deleteRequestTemplateByIDMapper(id);
    }

    @Override
    public int deleteByDomain(FormItemTemplate formItemTemplate) {
        return this.getBaseMapper().deleteByDomain(formItemTemplate);
    }

    @Override
    public TaskServiceMetaResp getTaskCreateServiceMeta(String procInstKey, String nodeDefId) {
        TaskServiceMetaResp resp = new TaskServiceMetaResp();
        resp.setFormItems(this.getBaseMapper().getCreateTaskServiceMeta(procInstKey,nodeDefId));
        return resp;
    }

}
