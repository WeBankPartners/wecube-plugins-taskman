package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.dto.req.ProcessingTasksReq;
import com.webank.taskman.dto.req.SaveTaskInfoReq;
import com.webank.taskman.dto.CoreCreateTaskDTO;
import com.webank.taskman.dto.TaskInfoDTO;
import com.webank.taskman.dto.req.QueryTaskInfoReq;
import com.webank.taskman.dto.req.SaveTaskInfoReq;
import com.webank.taskman.dto.req.SynthesisTaskInfoReq;
import com.webank.taskman.dto.resp.*;
import com.webank.taskman.support.core.CommonResponseDto;

import java.util.List;


public interface TaskInfoService extends IService<TaskInfo> {

    QueryResponse<TaskInfoDTO> selectTaskInfo(Integer page, Integer pageSize, QueryTaskInfoReq req);


    SynthesisTaskInfoFormTask selectSynthesisTaskInfoFormService(String id) throws Exception;

    RequestInfoInstanceResq selectTaskInfoInstanceService(String taskId, String requestId);

    TaskInfoGetResp getTheTaskInfoService(String id);

    CommonResponseDto createTask(CoreCreateTaskDTO req);

    List<FormItemInfoResp> returnDetail(String id);

    String ProcessingTasksService(ProcessingTasksReq ptr) throws Exception;
}
