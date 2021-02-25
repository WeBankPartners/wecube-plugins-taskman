package com.webank.taskman.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.webank.taskman.domain.RoleRelation;
import org.apache.ibatis.annotations.Param;

import java.util.List;


public interface RoleRelationMapper extends BaseMapper<RoleRelation> {


    int deleteByTemplate(@Param("tempId")  String tempId);

}
