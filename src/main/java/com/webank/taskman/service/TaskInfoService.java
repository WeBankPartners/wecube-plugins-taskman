package com.webank.taskman.service;


import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.dto.CoreCancelTaskDTO;
import com.webank.taskman.dto.CoreCreateTaskDTO;
import com.webank.taskman.dto.TaskInfoDTO;
import com.webank.taskman.dto.req.ProcessingTasksReq;
import com.webank.taskman.dto.req.QueryTaskInfoReq;
import com.webank.taskman.dto.resp.RequestInfoInstanceResq;
import com.webank.taskman.dto.resp.TaskInfoResp;
import com.webank.taskman.support.core.CommonResponseDto;

public interface TaskInfoService extends IService<TaskInfo> {

    QueryResponse<TaskInfoDTO> selectTaskInfo(Integer page, Integer pageSize, QueryTaskInfoReq req);

    RequestInfoInstanceResq selectTaskInfoInstanceService(String taskId, String requestId);

    TaskInfoDTO taskInfoReceive(String id);

    CommonResponseDto createTask(CoreCreateTaskDTO req);

    String ProcessingTasksService(ProcessingTasksReq ptr) throws TaskmanRuntimeException;

    CommonResponseDto cancelTask(CoreCancelTaskDTO req);

    TaskInfoResp taskInfoDetail(String id);
}
