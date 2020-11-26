package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.update.LambdaUpdateWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.converter.TemplateGroupConverter;
import com.webank.taskman.domain.TemplateGroup;
import com.webank.taskman.dto.TemplateGroupDTO;
import com.webank.taskman.dto.TemplateGroupVO;
import com.webank.taskman.mapper.TemplateGroupMapper;
import com.webank.taskman.service.TemplateGroupService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;


@Service
public class TemplateGroupServiceImpl extends ServiceImpl<TemplateGroupMapper, TemplateGroup> implements TemplateGroupService {

    @Autowired
    TemplateGroupMapper templateGroupMapper;

    @Autowired
    TemplateGroupConverter templateGroupConverter;

    @Override
    public void createTemplateGroupService(TemplateGroupVO templateGroupVO) throws Exception {
        if (templateGroupVO == null) {
            throw new Exception("Template group objects cannot be empty");
        }
        TemplateGroup templateGroup = templateGroupConverter.voToDomain(templateGroupVO);
        templateGroupMapper.insert(templateGroup);
    }

    @Override
    public void updateTemplateGroupService(TemplateGroupVO templateGroupVO) throws Exception {
        if (templateGroupVO == null) {
            throw new Exception("Template group objects cannot be empty");
        }
        TemplateGroup templateGroup = templateGroupConverter.voToDomain(templateGroupVO);
        templateGroupMapper.updateById(templateGroup);
    }

    @Override
    public List<TemplateGroupDTO> selectAllTemplateGroupService() throws Exception {
        List<TemplateGroup> templateGroups = templateGroupMapper.selectList(null);
        return templateGroupConverter.toDto(templateGroups);
    }

    @Override
    public void deleteTemplateGroupByIDService(String id) throws Exception {
        LambdaUpdateWrapper<TemplateGroup> lambdaUpdateWrapper = new LambdaUpdateWrapper<>();
        lambdaUpdateWrapper.eq(TemplateGroup::getId, id).set(TemplateGroup::getDelFlag, 1);
        int update = templateGroupMapper.update(null, lambdaUpdateWrapper);
        if (update != 1) {
            throw new Exception("Template group deletion failed");
        }
    }
}
