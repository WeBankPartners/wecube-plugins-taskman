package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.constant.RoleTypeEnum;
import com.webank.taskman.domain.RoleRelation;
import com.webank.taskman.dto.RoleDTO;
import org.apache.ibatis.annotations.Param;

import java.util.List;

public interface RoleRelationService extends IService<RoleRelation> {


    int deleteByTemplate(String tempId);

    void saveRoleRelation(String recordId,RoleTypeEnum roleType,List<RoleDTO> roles);

    void saveRoleRelationByTemplate(String recordId,List<RoleDTO> useRoles,List<RoleDTO> manageRoles);
}
