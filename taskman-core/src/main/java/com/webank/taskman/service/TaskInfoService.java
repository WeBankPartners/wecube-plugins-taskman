package com.webank.taskman.service;


import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.LocalPageableQueryResult;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.dto.TaskInfoDto;
import com.webank.taskman.dto.platform.CoreCancelTaskDto;
import com.webank.taskman.dto.req.ProceedUserTaskReqDto;
import com.webank.taskman.dto.req.TaskInfoQueryReqDto;
import com.webank.taskman.dto.resp.RequestInfoInstanceResqDto;
import com.webank.taskman.dto.resp.TaskInfoRespDto;
import com.webank.taskman.support.platform.dto.CommonPlatformResponseDto;

public interface TaskInfoService extends IService<TaskInfo> {

    LocalPageableQueryResult<TaskInfoDto> selectTaskInfo(Integer page, Integer pageSize, TaskInfoQueryReqDto req);

    RequestInfoInstanceResqDto selectTaskInfoInstanceService(String requestId, String taskId);

    TaskInfoDto taskInfoReceive(String id);

    void proceedUserTask(ProceedUserTaskReqDto ptr);

    CommonPlatformResponseDto cancelTask(CoreCancelTaskDto req);

    TaskInfoRespDto taskInfoDetail(String id);
}
