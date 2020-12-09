package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.conditions.update.UpdateWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.commons.TaskmanException;
import com.webank.taskman.converter.FormItemTemplateConverter;
import com.webank.taskman.converter.FormTemplateConverter;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.domain.FormTemplate;
import com.webank.taskman.dto.PageInfo;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.SaveFormItemTemplateReq;
import com.webank.taskman.dto.req.SaveFormTemplateReq;
import com.webank.taskman.dto.resp.FormTemplateResp;
import com.webank.taskman.mapper.FormTemplateMapper;
import com.webank.taskman.service.FormItemTemplateService;
import com.webank.taskman.service.FormTemplateService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.StringUtils;

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
    public QueryResponse<FormTemplateResp> selectFormTemplate(Integer current, Integer limit, SaveFormTemplateReq req) throws Exception {
        Page<FormTemplate> page=new Page<>(current,limit);
        QueryWrapper<FormTemplate> wrapper=new QueryWrapper<>();
        wrapper.eq(!StringUtils.isEmpty(req.getId()),"id",req.getId());
        wrapper.eq(!StringUtils.isEmpty(req.getTempId()),"temp_id",req.getTempId());
        wrapper.like(!StringUtils.isEmpty(req.getTempType()),"temp_type",req.getTempType());
        wrapper.like(!StringUtils.isEmpty(req.getName()),"name",req.getName());
        wrapper.like(!StringUtils.isEmpty(req.getDescription()),"description",req.getDescription());
        wrapper.like(!StringUtils.isEmpty(req.getStyle()),"style",req.getStyle());

        IPage<FormTemplate> iPage=formTemplateMapper.selectPage(page,wrapper);
        List<FormTemplate> list=iPage.getRecords();
        List<FormTemplateResp> formTemplateResps=formTemplateConverter.toDto(list);

        QueryResponse<FormTemplateResp> queryResponse=new QueryResponse<>();
        PageInfo pageInfo=new PageInfo();
        pageInfo.setStartIndex(iPage.getCurrent());
        pageInfo.setPageSize(iPage.getSize());
        pageInfo.setTotalRows(iPage.getTotal());
        queryResponse.setPageInfo(pageInfo);
        queryResponse.setContents(formTemplateResps);

        return queryResponse;
    }

    @Override
    public void deleteFormTemplate(String id) throws TaskmanException {
        if (StringUtils.isEmpty(id)){
            throw new TaskmanException("Form template parameter cannot be ID");
        }
        UpdateWrapper<FormTemplate> wrapper=new UpdateWrapper<>();
        wrapper.eq("id",id).set("del_flag",1);
        formTemplateMapper.update(null,wrapper);
    }

    @Override
    public FormTemplateResp detailFormTemplate(SaveFormTemplateReq req) {
        QueryWrapper<FormTemplate> wrapper=new QueryWrapper<>();
        wrapper.eq(!StringUtils.isEmpty(req.getId()),"id",req.getId());
        wrapper.eq(!StringUtils.isEmpty(req.getTempId()),"temp_id",req.getTempId());
        wrapper.eq(!StringUtils.isEmpty(req.getTempType()),"temp_type",req.getTempType());

        FormTemplateResp formTemplateResp = formTemplateConverter.toDto(formTemplateMapper.selectOne(wrapper));
        if(null !=formTemplateResp){
            formTemplateResp.setItems(formItemTemplateConverter.toDto(
                    formItemTemplateService.list(new QueryWrapper<FormItemTemplate>().
                            eq("form_template_id",formTemplateResp.getId()))
                    )
            );
        }
        return formTemplateResp;
    }

    @Override
    @Transactional
    public FormTemplateResp saveFormTemplateByReq(SaveFormTemplateReq formTemplateReq) throws TaskmanException {
        FormTemplate formTemplate= formTemplateConverter.reqToDomain(formTemplateReq);
        formTemplate.setCurrenUserName(formTemplate,formTemplate.getId());
        saveOrUpdate(formTemplate);
        List<SaveFormItemTemplateReq> items = formTemplateReq.getFormItems();
        for(SaveFormItemTemplateReq item:items){
            formItemTemplateService.saveFormItemTemplateByReq(item);
        }

        return new FormTemplateResp().setId(formTemplate.getId());
    }



    @Override
    public FormTemplateResp queryDetailByTemp(Integer tempType, String tempId) throws TaskmanException {
        QueryWrapper<FormTemplate> wrapper=new QueryWrapper<>();
        wrapper.eq("temp_id",tempId);
        wrapper.eq("temp_type",tempType);
        return formTemplateConverter.toDto(getOne(wrapper));
    }
}
