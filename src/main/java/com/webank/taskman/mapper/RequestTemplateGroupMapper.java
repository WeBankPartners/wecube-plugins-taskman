package com.webank.taskman.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.webank.taskman.domain.RequestTemplateGroup;
import org.apache.ibatis.annotations.Param;


public interface RequestTemplateGroupMapper extends BaseMapper<RequestTemplateGroup> {

    public void deleteTemplateGroupByIDMapper(@Param("id") String id);


}
