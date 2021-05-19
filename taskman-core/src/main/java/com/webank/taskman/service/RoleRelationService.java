package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.constant.RoleType;
import com.webank.taskman.domain.RoleRelation;
import com.webank.taskman.dto.RoleDto;

import java.util.List;

public interface RoleRelationService extends IService<RoleRelation> {


    int deleteByTemplate(String tempId);

    void saveRoleRelation(String recordId,RoleType roleType,List<RoleDto> roles);

    void saveRoleRelationByTemplate(String recordId,List<RoleDto> useRoles,List<RoleDto> manageRoles);
}
