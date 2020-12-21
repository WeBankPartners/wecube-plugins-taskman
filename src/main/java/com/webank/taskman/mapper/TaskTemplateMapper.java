package com.webank.taskman.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.webank.taskman.domain.TaskTemplate;
import com.webank.taskman.dto.req.TaskTemplateReq;
import org.apache.ibatis.annotations.Param;


public interface TaskTemplateMapper extends BaseMapper<TaskTemplate> {


    IPage<TaskTemplate> selectSynthesisRequestTemple(Page page, @Param("param") TaskTemplateReq req);
}
