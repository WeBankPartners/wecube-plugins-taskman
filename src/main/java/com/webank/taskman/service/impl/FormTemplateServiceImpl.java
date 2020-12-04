package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.conditions.update.UpdateWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.converter.FormTemplateConverter;
import com.webank.taskman.dto.PageInfo;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.SaveFormTemplateReq;
import com.webank.taskman.dto.resp.FormTemplateResp;
import com.webank.taskman.mapper.FormTemplateMapper;
import com.webank.taskman.service.FormTemplateService;
import com.webank.taskman.domain.FormTemplate;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;


@Service
public class FormTemplateServiceImpl extends ServiceImpl<FormTemplateMapper, FormTemplate> implements FormTemplateService {

    @Autowired
    FormTemplateMapper formTemplateMapper;
    @Autowired
    FormTemplateConverter formTemplateConverter;

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
    public void deleteFormTemplate(String id) throws Exception {
        if (id==""){
            throw new Exception("Form template parameter cannot be ID");
        }
        UpdateWrapper<FormTemplate> wrapper=new UpdateWrapper<>();
        wrapper.eq("id",id).set("del_flag",1);
        formTemplateMapper.update(null,wrapper);
    }

    @Override
    public FormTemplateResp detailFormTemplate(SaveFormTemplateReq req) throws Exception {
        QueryWrapper<FormTemplate> wrapper=new QueryWrapper<>();
        wrapper.eq(!StringUtils.isEmpty(req.getId()),"id",req.getId());
        wrapper.eq(!StringUtils.isEmpty(req.getTempId()),"temp_id",req.getTempId());
        wrapper.eq(!StringUtils.isEmpty(req.getTempType()),"temp_type",req.getTempType());
        FormTemplate formTemplate=formTemplateMapper.selectOne(wrapper);
        FormTemplateResp formTemplateResp=formTemplateConverter.toDto(formTemplate);
        return formTemplateResp;
    }

    @Override
    public FormTemplateResp saveFormTemplate(SaveFormTemplateReq formTemplateReq) {
        FormTemplate formTemplate=formTemplateConverter.reqToDomain(formTemplateReq);
        FormTemplateResp formTemplateResp=new FormTemplateResp();
        formTemplate.setCreatedBy("zhangsan");
        formTemplate.setUpdatedBy("lisi");
        if (org.springframework.util.StringUtils.isEmpty(formTemplateReq.getId())){
            formTemplateMapper.insert(formTemplate);
            formTemplateResp.setId(formTemplate.getId());
        }else{
            formTemplateMapper.updateById(formTemplate);
            formTemplateResp.setId(formTemplate.getId());
        }
        return formTemplateResp;
    }
}
