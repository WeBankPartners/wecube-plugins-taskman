package com.webank.taskman.mapper;

import org.apache.ibatis.annotations.Param;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.webank.taskman.domain.RoleRelation;


public interface RoleRelationMapper extends BaseMapper<RoleRelation> {


    int deleteByTemplate(@Param("tempId")  String tempId);

}
