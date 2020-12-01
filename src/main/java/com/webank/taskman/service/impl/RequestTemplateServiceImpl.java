package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.conditions.update.UpdateWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.converter.RequestTemplateConverter;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.*;
import com.webank.taskman.dto.req.AddRequestTemplateReq;
import com.webank.taskman.mapper.RequestTemplateMapper;
import com.webank.taskman.service.RequestTemplateService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.List;


@Service
public class RequestTemplateServiceImpl extends ServiceImpl<RequestTemplateMapper, RequestTemplate> implements RequestTemplateService {

   @Autowired
   RequestTemplateMapper requestTemplateMapper;

   @Autowired
    RequestTemplateConverter requestTemplateConverter;

    @Override
    public void createRequestTemplateService(RequestTemplateVO requestTemplateVO) throws Exception {
        if (requestTemplateVO == null) {
            throw new Exception("Template group objects cannot be empty");
        }
        RequestTemplate requestTemplate=requestTemplateConverter.voToDomain(requestTemplateVO);
        requestTemplateMapper.insert(requestTemplate);
    }

    @Override
    public void deleteRequestTemplateService(String id) throws Exception {
        if (id == "") {
            throw new Exception("Template group objects cannot be empty");
        }
        UpdateWrapper<RequestTemplate> wrapper = new UpdateWrapper<>();
        wrapper.eq("id",id).set("del_flag",1);
        requestTemplateMapper.update(null,wrapper);
    }

    @Override
    public void updateRequestTemplateService(RequestTemplateVO requestTemplateVO) throws Exception {
        if (requestTemplateVO==null){
            throw new Exception("Template group objects cannot be empty");
        }
        RequestTemplate requestTemplate=requestTemplateConverter.voToDomain(requestTemplateVO);
        SimpleDateFormat df = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
        requestTemplate.setUpdatedTime(df.parse(df.format(new Date())));
        requestTemplateMapper.updateById(requestTemplate);
    }

    @Override
    public QueryResponse<RequestTemplateDTO> selectAllequestTemplateService(Integer current, Integer limit, RequestTemplateReq req) throws Exception {
        Page<RequestTemplate> page = new Page<>(current, limit);
        QueryWrapper<RequestTemplate> wrapper = new QueryWrapper<>();
        wrapper.eq(!StringUtils.isEmpty(req.getId()),"id",req.getId());
        wrapper.eq(!StringUtils.isEmpty(req.getFormTempId()),"formTempId",req.getFormTempId());
        wrapper.eq(!StringUtils.isEmpty(req.getGroupId()),"groupId",req.getGroupId());
        wrapper.like(!StringUtils.isEmpty(req.getName()),"name",req.getName());

        IPage<RequestTemplate> iPage = requestTemplateMapper.selectPage(page, wrapper);
        List<RequestTemplate> records = iPage.getRecords();
        List<RequestTemplateDTO> requestTemplateDTOS = requestTemplateConverter.toDto(records);

        QueryResponse<RequestTemplateDTO> queryResponse = new QueryResponse<>();
        PageInfo pageInfo = new PageInfo();
        pageInfo.setStartIndex(iPage.getCurrent());
        pageInfo.setPageSize(iPage.getSize());
        pageInfo.setTotalRows(iPage.getTotal());
        queryResponse.setPageInfo(pageInfo);
        queryResponse.setContents(requestTemplateDTOS);
        return queryResponse;
    }


    @Override
    public List<RequestTemplateDTO> selectRequestTemplateService(RequestTemplateVO requestTemplateVO) throws Exception {
        return null;
    }

}
