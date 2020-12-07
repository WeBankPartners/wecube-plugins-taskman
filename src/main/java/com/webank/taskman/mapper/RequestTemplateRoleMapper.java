package com.webank.taskman.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.webank.taskman.domain.RequestTemplateRole;
import org.apache.ibatis.annotations.Param;

import java.util.List;


public interface RequestTemplateRoleMapper extends BaseMapper<RequestTemplateRole> {

    int deleteByRequestTemplate(@Param("requestTemplateId") String requestTemplateId );

}
