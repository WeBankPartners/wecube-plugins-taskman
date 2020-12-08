package com.webank.taskman.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.webank.taskman.domain.RoleRelation;
import org.apache.ibatis.annotations.Param;



public interface RoleRelationMapper extends BaseMapper<RoleRelation> {


    int deleteByTemplate(@Param("tempTable") String tempName, @Param("tempId")  String tempId);
}
