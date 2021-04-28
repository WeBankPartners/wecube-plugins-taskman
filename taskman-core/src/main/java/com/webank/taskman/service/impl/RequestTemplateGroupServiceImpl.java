package com.webank.taskman.service.impl;

import java.util.Date;
import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.StringUtils;

import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.base.LocalPageInfo;
import com.webank.taskman.base.LocalPageableQueryResult;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.GenernalStatus;
import com.webank.taskman.converter.RequestTemplateGroupConverter;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.RequestTemplateGroupDto;
import com.webank.taskman.mapper.RequestTemplateGroupMapper;
import com.webank.taskman.service.RequestTemplateGroupService;

@Service
public class RequestTemplateGroupServiceImpl extends ServiceImpl<RequestTemplateGroupMapper, RequestTemplateGroup>
        implements RequestTemplateGroupService {

    @Autowired
    private RequestTemplateGroupMapper templateGroupMapper;

    @Autowired
    private RequestTemplateGroupConverter requestTemplateGroupConverter;
    
    public List<RequestTemplateGroupDto> fetchAvailableGroupTemplates(){
        LambdaQueryWrapper<RequestTemplateGroup> lambdaQueryWrapper = new RequestTemplateGroup()
                .setStatus(RequestTemplateGroup.STATUS_AVAILABLE).getLambdaQueryWrapper();
        List<RequestTemplateGroupDto> requestTemplateGroupDtos = requestTemplateGroupConverter
                .convertToDtos(this.list(lambdaQueryWrapper));
        
        return requestTemplateGroupDtos;
    }

    /**
     * 
     */
    @Override
    @Transactional
    public RequestTemplateGroupDto saveOrUpdateTemplateGroup(RequestTemplateGroupDto requestTemplateGroupDto) {
        RequestTemplateGroup requestTemplateGroup = requestTemplateGroupConverter
                .convertToEntity(requestTemplateGroupDto);
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
        return requestTemplateGroupConverter.convertToDto(requestTemplateGroup);
    }

    /**
     * 
     */
    @Override
    public LocalPageableQueryResult<RequestTemplateGroupDto> searchRequestTemplateGroups(Integer current, Integer limit,
            RequestTemplateGroupDto requestTemplateGroupSearchCriteriaDto) {
        IPage<RequestTemplateGroup> iPage = templateGroupMapper.selectPage(new Page<>(current, limit),
                requestTemplateGroupConverter.convertToEntity(requestTemplateGroupSearchCriteriaDto)
                        .getLambdaQueryWrapper());

        // TODO to test start index here
        LocalPageInfo pageInfo = new LocalPageInfo(iPage.getTotal(), iPage.getCurrent() * iPage.getSize(),
                iPage.getSize());
        return new LocalPageableQueryResult<RequestTemplateGroupDto>(pageInfo,
                requestTemplateGroupConverter.convertToDtos(iPage.getRecords()));
    }

    @Override
    public void deleteRequestTemplateGroup(String id) {
        RequestTemplateGroup requestTemplateGroup = new RequestTemplateGroup();
        requestTemplateGroup.setId(id);
        requestTemplateGroup.setDelFlag(GenernalStatus.ENABLE.ordinal());
        requestTemplateGroup.setUpdatedTime(new Date());
        templateGroupMapper.update(null, requestTemplateGroup.getUpdateWrapper());
    }
}
