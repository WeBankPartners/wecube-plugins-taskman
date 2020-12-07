package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.RoleRelation;
import org.apache.ibatis.annotations.Param;

public interface RoleRelationService extends IService<RoleRelation> {


    int deleteByTemplate( String tempName, String tempId);
}
