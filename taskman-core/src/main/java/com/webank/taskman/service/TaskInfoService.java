package com.webank.taskman.service;


import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.JsonResponse;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.dto.TaskInfoDto;
import com.webank.taskman.dto.platform.CoreCancelTaskDto;
import com.webank.taskman.dto.req.ProcessingTasksReqDto;
import com.webank.taskman.dto.req.TaskInfoQueryReqDto;
import com.webank.taskman.dto.resp.RequestInfoInstanceResqDto;
import com.webank.taskman.dto.resp.TaskInfoRespDto;
import com.webank.taskman.support.core.CommonResponseDto;

public interface TaskInfoService extends IService<TaskInfo> {

    QueryResponse<TaskInfoDto> selectTaskInfo(Integer page, Integer pageSize, TaskInfoQueryReqDto req);

    RequestInfoInstanceResqDto selectTaskInfoInstanceService(String requestId, String taskId);

    TaskInfoDto taskInfoReceive(String id);

//    CommonResponseDto createTask(PlatformTaskCreationReqDto req);

    JsonResponse taskInfoProcessing(ProcessingTasksReqDto ptr) throws TaskmanRuntimeException;

    CommonResponseDto cancelTask(CoreCancelTaskDto req);

    TaskInfoRespDto taskInfoDetail(String id);
}