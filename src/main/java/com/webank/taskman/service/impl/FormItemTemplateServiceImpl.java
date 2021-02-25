package com.webank.taskman.service.impl;


import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.converter.FormItemInfoConverter;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.resp.FormItemInfoRespDto;
import com.webank.taskman.dto.resp.TaskServiceMetaRespDto;
import com.webank.taskman.mapper.FormItemTemplateMapper;
import com.webank.taskman.service.FormItemTemplateService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;


@Service
public class FormItemTemplateServiceImpl extends ServiceImpl<FormItemTemplateMapper, FormItemTemplate> implements FormItemTemplateService {

    @Autowired
    private FormItemInfoConverter formItemInfoConverter;

    @Override
    public void deleteRequestTemplateByID(String id) {
        this.getBaseMapper().deleteRequestTemplateByIDMapper(id);
    }

    @Override
    public int deleteByDomain(FormItemTemplate formItemTemplate) {
        return this.getBaseMapper().deleteByDomain(formItemTemplate);
    }

    @Override
    public TaskServiceMetaRespDto getTaskCreateServiceMeta(String procInstId, String nodeDefId) {
        TaskServiceMetaRespDto resp = new TaskServiceMetaRespDto();
        List<FormItemInfoRespDto> list = getBaseMapper().getCreateTaskServiceMeta(procInstId,nodeDefId);
        resp.setFormItems(formItemInfoConverter.respToServiceMeta(list));
        return resp;
    }

}
