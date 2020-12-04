package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.converter.RequestTemplateGroupConverter;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.*;
import com.webank.taskman.dto.req.SaveAndUpdateTemplateGropReq;
import com.webank.taskman.mapper.RequestTemplateGroupMapper;
import com.webank.taskman.service.RequestTemplateGroupService;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.List;


@Service
public class RequestTemplateGroupServiceImpl extends ServiceImpl<RequestTemplateGroupMapper, RequestTemplateGroup> implements RequestTemplateGroupService {

    @Autowired
    RequestTemplateGroupMapper templateGroupMapper;

    @Autowired
    RequestTemplateGroupConverter requestTemplateGroupConverter;


    @Override
    public RequestTemplateGroup addOrUpdateTemplateGroup(SaveAndUpdateTemplateGropReq gropReq) throws Exception {
        RequestTemplateGroup req = requestTemplateGroupConverter.addOrUpdateDomain(gropReq);
        if (StringUtils.isEmpty(req.getId())) {
            req.setCreatedBy("11");
            req.setUpdatedBy("22");
            templateGroupMapper.insert(req);
            RequestTemplateGroup requestTemplateGroup = templateGroupMapper.selectById(req);
            return requestTemplateGroup;
        }
        if (!StringUtils.isEmpty(req.getId())) {
            req.setUpdatedTime(new Date());
            templateGroupMapper.updateById(req);
            RequestTemplateGroup requestTemplateGroup = templateGroupMapper.selectById(req);
            return requestTemplateGroup;
        }
        return null;
    }

    @Override
    public void updateTemplateGroupService(TemplateGroupVO templateGroupVO) throws Exception {
        if (templateGroupVO == null) {
            throw new Exception("Template group objects cannot be empty");
        }
        RequestTemplateGroup templateGroup = requestTemplateGroupConverter.voToDomain(templateGroupVO);
        //templateGroup.setUpdatedBy();
        templateGroupMapper.updateById(templateGroup);
    }

    @Override
    public QueryResponse<TemplateGroupDTO> selectAllTemplateGroupService(Integer current, Integer limit, TemplateGroupReq req) throws Exception {
        Page<RequestTemplateGroup> page = new Page<>(current, limit);
        QueryWrapper<RequestTemplateGroup> wrapper = new QueryWrapper<>();
        if (!StringUtils.isEmpty(req.getId())) {
            wrapper.eq("id", req.getId());
        }
        if (!StringUtils.isEmpty(req.getName())) {
            wrapper.like("name", req.getName());
        }
        if (!StringUtils.isEmpty(req.getManageRole())) {
            wrapper.eq("manage_role_id", req.getManageRole());
        }

        IPage<RequestTemplateGroup> iPage = templateGroupMapper.selectPage(page, wrapper);
        List<RequestTemplateGroup> records = iPage.getRecords();
        List<TemplateGroupDTO> templateGroupDTOS = requestTemplateGroupConverter.toDto(records);

        QueryResponse<TemplateGroupDTO> queryResponse = new QueryResponse<>();
        PageInfo pageInfo = new PageInfo();
        pageInfo.setStartIndex(iPage.getCurrent());
        pageInfo.setPageSize(iPage.getSize());
        pageInfo.setTotalRows(iPage.getTotal());
        queryResponse.setPageInfo(pageInfo);
        queryResponse.setContents(templateGroupDTOS);
        return queryResponse;
    }

    @Override
    public void deleteTemplateGroupByIDService(String id) throws Exception {

        templateGroupMapper.deleteTemplateGroupByIDMapper(id);
    }
}
