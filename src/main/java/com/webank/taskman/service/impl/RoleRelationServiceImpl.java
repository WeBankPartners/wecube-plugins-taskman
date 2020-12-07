package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.domain.RoleRelation;
import com.webank.taskman.mapper.RoleRelationMapper;
import com.webank.taskman.service.RoleRelationService;
import org.springframework.stereotype.Service;

@Service
public class RoleRelationServiceImpl  extends ServiceImpl<RoleRelationMapper, RoleRelation> implements RoleRelationService {


    @Override
    public int deleteByTemplate(String tempName, String tempId) {
        return this.getBaseMapper().deleteByTemplate(tempName,tempId);
    }
}
