package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.StatusEnum;
import com.webank.taskman.converter.FormItemInfoConverter;
import com.webank.taskman.domain.FormInfo;
import com.webank.taskman.domain.FormItemInfo;
import com.webank.taskman.domain.FormTemplate;
import com.webank.taskman.dto.req.SaveFormInfoReq;
import com.webank.taskman.dto.req.SaveFormItemInfoReq;
import com.webank.taskman.mapper.FormInfoMapper;
import com.webank.taskman.service.FormInfoService;
import com.webank.taskman.service.FormItemInfoService;
import com.webank.taskman.service.FormTemplateService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;


@Service
public class FormInfoServiceImpl extends ServiceImpl<FormInfoMapper, FormInfo> implements FormInfoService {

    @Autowired
    FormTemplateService formTemplateService;

    @Autowired
    FormItemInfoService formItemInfoService;



    @Override
    public FormInfo saveFormInfoByExists(String requestTempId, String requestInfoId) {
        FormTemplate formTemplate = formTemplateService.getOne(new FormTemplate(null, requestTempId, StatusEnum.DEFAULT.ordinal() + "").getLambdaQueryWrapper());
        if (null == formTemplate) {
            throw new TaskmanRuntimeException("The FormTemplate do not exist");
        }
        remove(new QueryWrapper<FormInfo>().setEntity(new FormInfo().setRecordId(requestInfoId)));
        FormInfo form = new FormInfo();
        form.setRecordId(requestInfoId);
        form.setFormTemplateId(formTemplate.getId());
        form.setCurrenUserName(form, form.getId());
        save(form);
        return form;
    }

    @Override
    public void saveFormInfoAndItems(List<FormItemInfo> formItems, String requestTempId, String requestInfoId) {
        FormInfo form = saveFormInfoByExists(requestTempId, requestInfoId);
        String formId = form.getId();
        formItemInfoService.saveItemInfoByList(formItems, requestInfoId, formId);
    }


}
