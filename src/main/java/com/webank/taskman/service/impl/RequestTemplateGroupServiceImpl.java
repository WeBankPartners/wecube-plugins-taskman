package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.update.UpdateWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.TaskmanException;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.StatusCodeEnum;
import com.webank.taskman.converter.RequestTemplateGroupConverter;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.RequestTemplateGroupDTO;
import com.webank.taskman.dto.req.SaveRequestTemplateGropReq;
import com.webank.taskman.mapper.RequestTemplateGroupMapper;
import com.webank.taskman.service.RequestTemplateGroupService;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.Date;


@Service
public class RequestTemplateGroupServiceImpl extends ServiceImpl<RequestTemplateGroupMapper, RequestTemplateGroup> implements RequestTemplateGroupService {

    @Autowired
    RequestTemplateGroupMapper templateGroupMapper;

    @Autowired
    RequestTemplateGroupConverter requestTemplateGroupConverter;


    @Override
    @Transactional
    public RequestTemplateGroupDTO saveTemplateGroupByReq(SaveRequestTemplateGropReq req)  throws TaskmanException {
        RequestTemplateGroup requestTemplateGroup = requestTemplateGroupConverter.saveReqToDomain(req);
        requestTemplateGroup.setCurrenUserName(requestTemplateGroup,requestTemplateGroup.getId());
        if(!StringUtils.isEmpty(requestTemplateGroup.getId())){
            RequestTemplateGroup query = this.getById(requestTemplateGroup.getId());
            if(null == query){
                throw new TaskmanRuntimeException(StatusCodeEnum.NOT_FOUND_RECORD);
            }
            requestTemplateGroup.setUpdatedTime(new Date());
            updateById(requestTemplateGroup);
        }else {
            save(requestTemplateGroup);
        }
        return new RequestTemplateGroupDTO().setId(requestTemplateGroup.getId());
    }

    @Override
    public QueryResponse<RequestTemplateGroupDTO> selectByParam(Integer current, Integer limit, RequestTemplateGroupDTO req) {
        IPage<RequestTemplateGroup> iPage = templateGroupMapper.selectPage(new Page<>(current, limit),
                requestTemplateGroupConverter.toEntity(req).getLambdaQueryWrapper());
        return new QueryResponse(iPage,requestTemplateGroupConverter.toDto(iPage.getRecords()));
    }

    @Override
    public void deleteTemplateGroupByIDService(String id)  {

        UpdateWrapper<RequestTemplateGroup> wrapper = new UpdateWrapper<>();
        wrapper.lambda().eq(RequestTemplateGroup::getId, id)
                .set(RequestTemplateGroup::getDelFlag, 1)
                .set(RequestTemplateGroup::getUpdatedTime,new Date());
        templateGroupMapper.update(null,wrapper);
    }
}
