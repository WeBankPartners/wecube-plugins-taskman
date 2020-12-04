package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.RequestTemplateRole;


public interface RequestTemplateRoleService extends IService<RequestTemplateRole> {

    int deleteByRequestTemplate(String requestTemplateId );
}
