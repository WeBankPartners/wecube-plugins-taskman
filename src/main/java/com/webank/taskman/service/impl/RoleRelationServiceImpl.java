package com.webank.taskman.service.impl;


import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.constant.RoleTypeEnum;
import com.webank.taskman.domain.RoleRelation;
import com.webank.taskman.dto.RoleDto;
import com.webank.taskman.mapper.RoleRelationMapper;
import com.webank.taskman.service.RoleRelationService;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.StringUtils;

import java.util.List;

@Service
public class RoleRelationServiceImpl  extends ServiceImpl<RoleRelationMapper, RoleRelation> implements RoleRelationService {


    @Override
    public int deleteByTemplate( String tempId) {
        return this.getBaseMapper().deleteByTemplate(tempId);
    }

    @Override
    public void saveRoleRelation(String recordId, RoleTypeEnum roleType, List<RoleDto> roles) {
        roles.stream().forEach(role-> {
            if(!StringUtils.isEmpty(role.getRoleName()) && !StringUtils.isEmpty(role.getDisplayName())  ){
                this.getBaseMapper().insert( new RoleRelation(recordId, roleType.getType(),role.getRoleName(),role.getDisplayName()));
            }
        });
    }

    @Override
    @Transactional
    public void saveRoleRelationByTemplate(String requestTemplateId, List<RoleDto> useRoles, List<RoleDto> manageRoles) {
        deleteByTemplate(requestTemplateId);
        saveRoleRelation(requestTemplateId, RoleTypeEnum.USE_ROLE,useRoles);
        saveRoleRelation(requestTemplateId,RoleTypeEnum.MANAGE_ROLE,manageRoles);
    }

}
