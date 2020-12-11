package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanException;
import com.webank.taskman.converter.FormItemTemplateConverter;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.PageInfo;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.SaveFormItemTemplateReq;
import com.webank.taskman.dto.req.SelectFormItemTemplateReq;
import com.webank.taskman.dto.resp.FormItemTemplateResq;
import com.webank.taskman.mapper.FormItemTemplateMapper;
import com.webank.taskman.service.FormItemTemplateService;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.BeanUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.List;


@Service
public class FormItemTemplateServiceImpl extends ServiceImpl<FormItemTemplateMapper, FormItemTemplate> implements FormItemTemplateService {

    @Autowired
    FormItemTemplateConverter formItemTemplateConverter;

    @Override
    public FormItemTemplate saveFormItemTemplateByReq(SaveFormItemTemplateReq req) throws TaskmanException {
        FormItemTemplate formItemTemplate = formItemTemplateConverter.toEntityBySaveReq(req);
        formItemTemplate.setUpdatedBy(AuthenticationContextHolder.getCurrentUsername());
        QueryWrapper<FormItemTemplate> queryWrapper = new QueryWrapper<FormItemTemplate>()
                .eq("form_template_id",req.getFormTemplateId())
                .eq("name",req.getName());
        FormItemTemplate queryBean = getBaseMapper().selectOne(queryWrapper);
        if(null != queryBean){
            formItemTemplate.setId(queryBean.getId());
            updateById(formItemTemplate);
        }else{
            formItemTemplate.setCreatedBy(AuthenticationContextHolder.getCurrentUsername());
            save(formItemTemplate);
        }
        return null;
    }

    @Override
    public void deleteRequestTemplateByID(String id) {
        this.getBaseMapper().deleteRequestTemplateByIDMapper(id);
    }


    @Override
    public QueryResponse<FormItemTemplateResq> selectAllFormItemTemplateService(Integer current, Integer limit, SelectFormItemTemplateReq req) {
        Page<FormItemTemplate> page = new Page<>(current, limit);
        QueryWrapper<FormItemTemplate> wrapper = new QueryWrapper<>();
        if (!StringUtils.isEmpty(req.getId())) {
            wrapper.eq("id", req.getId());
        }
        if (!StringUtils.isEmpty(req.getName())) {
            wrapper.like("name", req.getName());
        }
        if (!StringUtils.isEmpty(req.getFormTemplateId())) {
            wrapper.eq("form_template_id", req.getFormTemplateId());
        }
        if (!StringUtils.isEmpty(req.getTitle())) {
            wrapper.like("title", req.getTitle());
        }
        if (!StringUtils.isEmpty(req.getElementType())) {
            wrapper.eq("element_type", req.getElementType());
        }
        if (!StringUtils.isEmpty(req.getDataCiId())) {
            wrapper.eq("data_ci_id", req.getDataCiId());
        }
        if (req.getIsPublic()!=null){
            wrapper.eq("is_public",req.getIsPublic());
        }
        IPage<FormItemTemplate> iPage = this.getBaseMapper().selectPage(page, wrapper);
        List<FormItemTemplate> records = iPage.getRecords();

        List<FormItemTemplateResq> formItemTemplateResqs = formItemTemplateConverter.toDto(records);
        QueryResponse<FormItemTemplateResq> queryResponse = new QueryResponse<>();
        PageInfo pageInfo = new PageInfo();
        pageInfo.setStartIndex(iPage.getCurrent());
        pageInfo.setPageSize(iPage.getSize());
        pageInfo.setTotalRows(iPage.getTotal());
        queryResponse.setPageInfo(pageInfo);
        queryResponse.setContents(formItemTemplateResqs);
        return queryResponse;
    }
}
