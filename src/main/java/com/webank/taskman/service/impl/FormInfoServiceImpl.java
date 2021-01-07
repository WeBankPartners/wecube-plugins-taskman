package com.webank.taskman.service.impl;


import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.StatusEnum;
import com.webank.taskman.domain.FormInfo;
import com.webank.taskman.domain.FormItemInfo;
import com.webank.taskman.domain.FormTemplate;
import com.webank.taskman.mapper.FormInfoMapper;
import com.webank.taskman.service.FormInfoService;
import com.webank.taskman.service.FormItemInfoService;
import com.webank.taskman.service.FormTemplateService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class FormInfoServiceImpl extends ServiceImpl<FormInfoMapper, FormInfo> implements FormInfoService {


    private static final Logger log = LoggerFactory.getLogger(FormInfoServiceImpl.class);

    @Autowired
    FormTemplateService formTemplateService;

    @Autowired
    FormItemInfoService formItemInfoService;



    @Override
    public FormInfo saveFormInfoByExists(String requestTempId, String recordId) {
        FormTemplate formTemplate = formTemplateService.getOne(new FormTemplate(null, requestTempId, StatusEnum.DEFAULT.ordinal() + "").getLambdaQueryWrapper());
        if (null == formTemplate) {
            throw new TaskmanRuntimeException("The FormTemplate do not exist");
        }
        remove(new QueryWrapper<FormInfo>().setEntity(new FormInfo().setRecordId(recordId)));
        FormInfo form = new FormInfo();
        form.setRecordId(recordId);
        form.setFormTemplateId(formTemplate.getId());
        form.setCurrenUserName(form, form.getId());
        save(form);
        return form;
    }

    @Override
    public void saveFormInfoAndItems(List<FormItemInfo> formItems, String requestTempId, String recordId) {
        if(null == formItems || formItems.size() == 0){
            log.info("The formItems is null.");
            return;
        }
        FormInfo form = saveFormInfoByExists(requestTempId, recordId);
        String formId = form.getId();
        formItemInfoService.saveItemInfoByList(formItems, recordId, formId);
    }


}
