package com.webank.taskman.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.dto.req.SelectTaskInfoReq;
import org.apache.ibatis.annotations.Param;


public interface TaskInfoMapper extends BaseMapper<TaskInfo> {
    IPage<TaskInfo> selectTaskInfo(Page page, @Param("Info") SelectTaskInfoReq selectTaskInfoReq);
}
