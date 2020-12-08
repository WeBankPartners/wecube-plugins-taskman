package com.webank.taskman.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.webank.taskman.domain.RoleRelation;
import org.apache.ibatis.annotations.Param;

import java.util.List;


public interface RoleRelationMapper extends BaseMapper<RoleRelation> {


    int deleteByTemplate(@Param("tempTable") String tempName, @Param("tempId")  String tempId);

    List<RoleRelation> selectRole(String recordTable, String recordId);
}
