package com.webank.taskman.service.impl;

import java.util.Date;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.StringUtils;

import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.StatusEnum;
import com.webank.taskman.converter.RequestTemplateGroupConverter;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.RequestTemplateGroupDto;
import com.webank.taskman.dto.req.RequestTemplateGroupSaveReqDto;
import com.webank.taskman.mapper.RequestTemplateGroupMapper;
import com.webank.taskman.service.RequestTemplateGroupService;

@Service
public class RequestTemplateGroupServiceImpl extends ServiceImpl<RequestTemplateGroupMapper, RequestTemplateGroup>
        implements RequestTemplateGroupService {

    @Autowired
    private RequestTemplateGroupMapper templateGroupMapper;

    @Autowired
    private RequestTemplateGroupConverter requestTemplateGroupConverter;

    @Override
    @Transactional
    public RequestTemplateGroupDto saveTemplateGroupByReq(RequestTemplateGroupSaveReqDto req) {
        RequestTemplateGroup requestTemplateGroup = requestTemplateGroupConverter.saveReqToDomain(req);
        requestTemplateGroup.setCreatedBy(AuthenticationContextHolder.getCurrentUsername());
        requestTemplateGroup.setUpdatedBy(AuthenticationContextHolder.getCurrentUsername());
        requestTemplateGroup.setStatus(RequestTemplateGroup.STATUS_AVAILABLE);
        if (!StringUtils.isEmpty(requestTemplateGroup.getId())) {
            RequestTemplateGroup query = this.getById(requestTemplateGroup.getId());
            if (query == null) {
                throw new TaskmanRuntimeException("NOT_FOUND_RECORD");
            }
            requestTemplateGroup.setUpdatedTime(new Date());
            update(requestTemplateGroup.getUpdateWrapper());
        } else {
            save(requestTemplateGroup);
        }
        return requestTemplateGroupConverter.toDto(requestTemplateGroup);
    }

    @Override
    public QueryResponse<RequestTemplateGroupDto> selectRequestTemplateGroupPage(Integer current, Integer limit,
            RequestTemplateGroupDto req) {
        IPage<RequestTemplateGroup> iPage = templateGroupMapper.selectPage(new Page<>(current, limit),
                requestTemplateGroupConverter.toEntity(req).getLambdaQueryWrapper());
        return new QueryResponse(iPage, requestTemplateGroupConverter.toDto(iPage.getRecords()));
    }

    @Override
    public void deleteTemplateGroupByIDService(String id) {
        RequestTemplateGroup requestTemplateGroup = new RequestTemplateGroup();
        requestTemplateGroup.setId(id);
        requestTemplateGroup.setDelFlag(StatusEnum.ENABLE.ordinal());
        requestTemplateGroup.setUpdatedTime(new Date());
        templateGroupMapper.update(null, requestTemplateGroup.getUpdateWrapper());
    }
}
