package com.webank.taskman.service.impl;

import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.TemplateType;
import com.webank.taskman.domain.FormInfo;
import com.webank.taskman.domain.FormItemInfo;
import com.webank.taskman.domain.FormTemplate;
import com.webank.taskman.mapper.FormInfoMapper;
import com.webank.taskman.service.FormInfoService;
import com.webank.taskman.service.FormItemInfoService;
import com.webank.taskman.service.FormTemplateService;

@Service
public class FormInfoServiceImpl extends ServiceImpl<FormInfoMapper, FormInfo> implements FormInfoService {

    private static final Logger log = LoggerFactory.getLogger(FormInfoServiceImpl.class);

    @Autowired
    private FormTemplateService formTemplateService;

    @Autowired
    private FormItemInfoService formItemInfoService;

    @Override
    public FormInfo saveFormInfoByExists(String requestTempId, String recordId) {
        LambdaQueryWrapper<FormTemplate> formTemplateQueryWrapper = new FormTemplate(null, requestTempId,
                TemplateType.REQUEST.getType()).getLambdaQueryWrapper();

        FormTemplate formTemplate = formTemplateService.getOne(formTemplateQueryWrapper);
        if (formTemplate == null) {
            throw new TaskmanRuntimeException("The form template do not exist");
        }
        
        remove(new QueryWrapper<FormInfo>().setEntity(new FormInfo().setRecordId(recordId)));
        
        FormInfo form = new FormInfo();
        form.setRecordId(recordId);
        form.setFormTemplateId(formTemplate.getId());
        form.setCreatedBy(AuthenticationContextHolder.getCurrentUsername());
        form.setUpdatedBy(AuthenticationContextHolder.getCurrentUsername());

        save(form);

        return form;
    }

    @Override
    public void saveFormInfoAndFormItems(List<FormItemInfo> formItems, String requestTempId, String recordId) {
        if (formItems == null || formItems.isEmpty()) {
            log.info("The formItems is null.");
            return;
        }
        FormInfo form = saveFormInfoByExists(requestTempId, recordId);
        String formId = form.getId();
        formItemInfoService.saveItemInfoByList(formItems, recordId, formId);
    }

}
