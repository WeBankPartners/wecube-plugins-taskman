package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.domain.RequestTemplateRole;
import com.webank.taskman.mapper.RequestTemplateRoleMapper;
import com.webank.taskman.service.RequestTemplateRoleService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class RequestTemplateRoleImpl extends ServiceImpl<RequestTemplateRoleMapper, RequestTemplateRole> implements RequestTemplateRoleService {

    @Autowired
    RequestTemplateRoleMapper requestTemplateRoleMapper;

    @Override
    public int deleteByRequestTemplate(String requestTemplateId) {
        return this.getBaseMapper().deleteByRequestTemplate(requestTemplateId);
    }
}
