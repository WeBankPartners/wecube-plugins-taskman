package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.conditions.update.UpdateWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.converter.RequestTemplateConverter;
import com.webank.taskman.domain.RequestTemplaeRole;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.*;
import com.webank.taskman.dto.req.SaveRequestTemplateReq;
import com.webank.taskman.dto.resp.RequestTemplateResp;
import com.webank.taskman.mapper.RequestTemplaeRoleMapper;
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

   @Autowired
    RequestTemplaeRoleMapper requestTemplaeRoleMapper;


    @Override
    public RequestTemplateResp saveRequestTemplate(SaveRequestTemplateReq requestTemplateReq)  {
        RequestTemplate requestTemplate=requestTemplateConverter.reqToDomain(requestTemplateReq);
        requestTemplate.setCreatedBy("zhangsan");
        requestTemplate.setUpdatedBy("lisi");
        String[] roleIds=requestTemplateReq.getRoleids().split(",");
        RequestTemplaeRole requestTemplaeRole=new RequestTemplaeRole();

        RequestTemplateResp requestTemplateResp=new RequestTemplateResp();
        if (StringUtils.isEmpty(requestTemplateReq.getId())){
            requestTemplateMapper.insert(requestTemplate);
            requestTemplateResp.setId(requestTemplate.getId());
            for (String roleId : roleIds) {
                requestTemplaeRole.setRoleId(roleId);
                requestTemplaeRole.setRequestTemplateId(requestTemplate.getId());
                requestTemplaeRoleMapper.insert(requestTemplaeRole);
                requestTemplaeRole.setId("");
            }
        }else{
            requestTemplateMapper.updateById(requestTemplate);
            requestTemplateResp.setId(requestTemplateReq.getId());
            if (requestTemplateMapper.selectById(requestTemplateReq.getId())==null){
                for (String roleId : roleIds) {
                    requestTemplaeRole.setRoleId(roleId);
                    requestTemplaeRole.setRequestTemplateId(requestTemplateReq.getId());
                    requestTemplaeRoleMapper.insert(requestTemplaeRole);
                    requestTemplaeRole.setId("");
                }
            }
        }

       return requestTemplateResp;

    }

    @Override
    public void deleteRequestTemplateService(String id) throws Exception {
        if (id == "") {
            throw new Exception("Request template parameter cannot be ID");
        }
        UpdateWrapper<RequestTemplate> wrapper = new UpdateWrapper<>();
        wrapper.eq("id",id).set("del_flag",1);
        requestTemplateMapper.update(null,wrapper);
    }

    @Override
    public QueryResponse<RequestTemplateResp> selectAllequestTemplateService(Integer current, Integer limit, SaveRequestTemplateReq req) throws Exception {
        Page<RequestTemplate> page = new Page<>(current, limit);
        QueryWrapper<RequestTemplate> wrapper = new QueryWrapper<>();
        wrapper.eq(!StringUtils.isEmpty(req.getId()),"id",req.getId());
        wrapper.eq(!StringUtils.isEmpty(req.getRequestTempGroup()),"request_temp_group",req.getRequestTempGroup());
        wrapper.eq(!StringUtils.isEmpty(req.getProcDefKey()),"proc_def_key",req.getProcDefKey());
        wrapper.eq(!StringUtils.isEmpty(req.getProcDefId()),"proc_def_id",req.getProcDefId());
        wrapper.eq(!StringUtils.isEmpty(req.getProcDefName()),"proc_def_name",req.getProcDefName());
        wrapper.eq(!StringUtils.isEmpty(req.getName()),"name",req.getName());
        IPage<RequestTemplate> iPage = requestTemplateMapper.selectPage(page, wrapper);
        List<RequestTemplate> records = iPage.getRecords();
        List<RequestTemplateResp> requestTemplateDTOS = requestTemplateConverter.toDto(records);

        QueryResponse<RequestTemplateResp> queryResponse = new QueryResponse<>();
        PageInfo pageInfo = new PageInfo();
        pageInfo.setStartIndex(iPage.getCurrent());
        pageInfo.setPageSize(iPage.getSize());
        pageInfo.setTotalRows(iPage.getTotal());
        queryResponse.setPageInfo(pageInfo);
        queryResponse.setContents(requestTemplateDTOS);
        return queryResponse;
    }

    @Override
    public RequestTemplateResp detailRequestTemplate(String id) throws Exception {
        RequestTemplate requestTemplate=requestTemplateMapper.selectById(id);
        RequestTemplateResp formTemplateResp=requestTemplateConverter.toDto(requestTemplate);
        return formTemplateResp;
    }


}
