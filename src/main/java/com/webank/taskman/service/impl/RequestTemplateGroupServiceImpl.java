package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.converter.TemplateGroupConverter;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.TemplateGroupDTO;
import com.webank.taskman.dto.TemplateGroupVO;
import com.webank.taskman.mapper.RequestTemplateGroupMapper;
import com.webank.taskman.service.RequestTemplateGroupService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;


@Service
public class RequestTemplateGroupServiceImpl extends ServiceImpl<RequestTemplateGroupMapper, RequestTemplateGroup> implements RequestTemplateGroupService {

    @Autowired
    RequestTemplateGroupMapper templateGroupMapper;

    @Autowired
    TemplateGroupConverter templateGroupConverter;


    @Override
    public void createTemplateGroupService(TemplateGroupVO templateGroupVO) throws Exception {
        RequestTemplateGroup templateGroup = templateGroupConverter.voToDomain(templateGroupVO);
        templateGroup.setCreatedBy("123");
        templateGroup.setUpdatedBy("123");
        templateGroupMapper.insert(templateGroup);
    }

    @Override
    public void updateTemplateGroupService(TemplateGroupVO templateGroupVO) throws Exception {
        if (templateGroupVO == null) {
            throw new Exception("Template group objects cannot be empty");
        }
        RequestTemplateGroup templateGroup = templateGroupConverter.voToDomain(templateGroupVO);
        templateGroupMapper.updateById(templateGroup);
    }

    @Override
    public List<TemplateGroupDTO> selectAllTemplateGroupService() throws Exception {
        List<RequestTemplateGroup> templateGroups = templateGroupMapper.selectList(null);
        return templateGroupConverter.toDto(templateGroups);
    }

    @Override
    public void deleteTemplateGroupByIDService(String id) throws Exception {

    }
}
