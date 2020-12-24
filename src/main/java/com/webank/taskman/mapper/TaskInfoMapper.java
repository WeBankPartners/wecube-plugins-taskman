package com.webank.taskman.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.dto.req.QueryTaskInfoReq;
import com.webank.taskman.dto.req.SynthesisTaskInfoReq;
import org.apache.ibatis.annotations.Param;

import java.util.List;


public interface TaskInfoMapper extends BaseMapper<TaskInfo> {

    List<TaskInfo> selectTaskInfo(QueryTaskInfoReq req);

    IPage<TaskInfo> selectSynthesisRequestInfo(Page page, @Param("param") SynthesisTaskInfoReq req);
}
