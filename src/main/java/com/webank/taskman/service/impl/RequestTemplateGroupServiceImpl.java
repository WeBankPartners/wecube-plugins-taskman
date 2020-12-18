package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.converter.RequestTemplateGroupConverter;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.PageInfo;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.TemplateGroupDTO;
import com.webank.taskman.dto.TemplateGroupReq;
import com.webank.taskman.dto.req.SaveRequestTemplateGropReq;
import com.webank.taskman.mapper.RequestTemplateGroupMapper;
import com.webank.taskman.service.RequestTemplateGroupService;
import com.webank.taskman.service.RoleRelationService;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;


@Service
public class RequestTemplateGroupServiceImpl extends ServiceImpl<RequestTemplateGroupMapper, RequestTemplateGroup> implements RequestTemplateGroupService {

    @Autowired
    RequestTemplateGroupMapper templateGroupMapper;

    @Autowired
    RequestTemplateGroupConverter requestTemplateGroupConverter;

    @Autowired
    RoleRelationService roleRelationService;


    @Override
    @Transactional
    public RequestTemplateGroup saveTemplateGroupByReq(SaveRequestTemplateGropReq req) throws Exception {
        RequestTemplateGroup requestTemplateGroup = requestTemplateGroupConverter.saveReqToDomain(req);
        requestTemplateGroup.setCurrenUserName(requestTemplateGroup,requestTemplateGroup.getId());
        saveOrUpdate(requestTemplateGroup);
//        roleRelationService.deleteByTemplate(requestTemplateGroup.getId());
        return new RequestTemplateGroup().setId(requestTemplateGroup.getId());
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
