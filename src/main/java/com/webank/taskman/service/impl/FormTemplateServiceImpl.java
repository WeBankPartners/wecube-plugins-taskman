package com.webank.taskman.service.impl;


import com.baomidou.mybatisplus.core.conditions.update.UpdateWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.StatusEnum;
import com.webank.taskman.converter.FormItemTemplateConverter;
import com.webank.taskman.converter.FormTemplateConverter;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.domain.FormTemplate;
import com.webank.taskman.dto.req.SaveFormTemplateReq;
import com.webank.taskman.dto.resp.FormTemplateResp;
import com.webank.taskman.mapper.FormTemplateMapper;
import com.webank.taskman.service.FormItemTemplateService;
import com.webank.taskman.service.FormTemplateService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.StringUtils;

import java.util.Date;
import java.util.List;

@Service
public class FormTemplateServiceImpl extends ServiceImpl<FormTemplateMapper, FormTemplate> implements FormTemplateService {

    @Autowired
    FormTemplateMapper formTemplateMapper;
    @Autowired
    FormTemplateConverter formTemplateConverter;

    @Autowired
    FormItemTemplateService formItemTemplateService;

    @Autowired
    FormItemTemplateConverter formItemTemplateConverter;

    @Override
    public QueryResponse<FormTemplateResp> selectFormTemplate(Integer page, Integer pageSize, SaveFormTemplateReq req){
        IPage<FormTemplate> iPage = formTemplateMapper.selectPage(new Page<>(page,pageSize),
                formTemplateConverter.reqToDomain(req).getLambdaQueryWrapper());
        List<FormTemplateResp> formTemplateResps = formTemplateConverter.toDto(iPage.getRecords());
        return new QueryResponse(iPage.getSize(), iPage.getCurrent(), iPage.getSize(), formTemplateResps);
    }

    @Override
    public void deleteFormTemplate(String id) throws TaskmanRuntimeException {
        if (StringUtils.isEmpty(id)){
            throw new TaskmanRuntimeException("Form template parameter cannot be ID");
        }
        UpdateWrapper<FormTemplate> wrapper = new UpdateWrapper<>();
        wrapper.lambda().eq(FormTemplate::getId,id)
                .set(FormTemplate::getDelFlag, StatusEnum.ENABLE.ordinal())
                .set(FormTemplate::getUpdatedTime,new Date());;
        formTemplateMapper.update(null,wrapper);
    }

    @Override
    public FormTemplateResp detailFormTemplate(SaveFormTemplateReq req) {

        FormTemplateResp formTemplateResp = formTemplateConverter.toDto(
                formTemplateMapper.selectOne(formTemplateConverter.reqToDomain(req).getLambdaQueryWrapper()));
        if(null !=formTemplateResp){
            formTemplateResp.setItems(formItemTemplateConverter.toRespByEntity(
                formItemTemplateService.list(new FormItemTemplate().setFormTemplateId(formTemplateResp.getId()).getLambdaQueryWrapper())));
        }
        return formTemplateResp;
    }

    @Override
    @Transactional
    public FormTemplateResp saveFormTemplateByReq(SaveFormTemplateReq formTemplateReq) throws TaskmanRuntimeException {
        FormTemplate formTemplate= formTemplateConverter.reqToDomain(formTemplateReq);
        formTemplate.setName(StringUtils.isEmpty(formTemplate.getName())?"":formTemplate.getName());
        formTemplate.setCurrenUserName(formTemplate,formTemplate.getId());
        saveOrUpdate(formTemplate);
        formItemTemplateService.deleteByDomain(new FormItemTemplate().setFormTemplateId(formTemplate.getId()));

        formTemplateReq.getFormItems().stream().forEach(item->{
            FormItemTemplate formItemTemplate = formItemTemplateConverter.toEntityBySaveReq(item);
            formItemTemplate.setFormTemplateId(formTemplate.getId());
            formItemTemplate.setTempId(formTemplate.getTempId());
            formItemTemplateService.save(formItemTemplate);
        });
        return new FormTemplateResp().setId(formTemplate.getId());
    }



    @Override
    public FormTemplateResp queryDetailByTemp(Integer tempType, String tempId) throws TaskmanRuntimeException {
        return formTemplateConverter.toDto(
                getOne(new FormTemplate().setTempId(tempId).setTempType(tempType+"").getLambdaQueryWrapper()));
    }
}
