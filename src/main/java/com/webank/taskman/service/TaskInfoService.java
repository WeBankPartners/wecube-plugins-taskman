package com.webank.taskman.service;


import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.JsonResponse;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.dto.CoreCancelTaskDto;
import com.webank.taskman.dto.CoreCreateTaskDto;
import com.webank.taskman.dto.TaskInfoDto;
import com.webank.taskman.dto.req.ProcessingTasksReq;
import com.webank.taskman.dto.req.QueryTaskInfoReq;
import com.webank.taskman.dto.resp.RequestInfoInstanceResq;
import com.webank.taskman.dto.resp.TaskInfoResp;
import com.webank.taskman.support.core.CommonResponseDto;

public interface TaskInfoService extends IService<TaskInfo> {

    QueryResponse<TaskInfoDto> selectTaskInfo(Integer page, Integer pageSize, QueryTaskInfoReq req);

    RequestInfoInstanceResq selectTaskInfoInstanceService(String requestId, String taskId);

    TaskInfoDto taskInfoReceive(String id);

    CommonResponseDto createTask(CoreCreateTaskDto req);

    JsonResponse taskInfoProcessing(ProcessingTasksReq ptr) throws TaskmanRuntimeException;

    CommonResponseDto cancelTask(CoreCancelTaskDto req);

    TaskInfoResp taskInfoDetail(String id);
}
