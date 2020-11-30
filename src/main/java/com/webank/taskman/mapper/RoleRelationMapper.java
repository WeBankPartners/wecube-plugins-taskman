package com.webank.taskman.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.webank.taskman.domain.RoleInfo;
import com.webank.taskman.domain.RoleRelation;
import org.apache.ibatis.annotations.Param;

import java.util.List;
import java.util.Map;


public interface RoleRelationMapper extends BaseMapper<RoleRelation> {

    List<RoleInfo> selectRoleInfoByParams(@Param("params") Map<String,Object> params);
}
