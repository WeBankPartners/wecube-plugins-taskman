package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.dto.TaskInfoDTO;
import com.webank.taskman.dto.req.QueryTaskInfoReq;
import com.webank.taskman.dto.req.SaveTaskInfoReq;
import com.webank.taskman.dto.req.SynthesisTaskInfoReq;
import com.webank.taskman.dto.resp.SaveTaskInfoResp;
import com.webank.taskman.dto.resp.SynthesisTaskInfoFormTask;
import com.webank.taskman.dto.resp.SynthesisTaskInfoResp;


public interface TaskInfoService extends IService<TaskInfo> {

    QueryResponse<TaskInfoDTO> selectTaskInfo(Integer page, Integer pageSize, QueryTaskInfoReq req);

    SaveTaskInfoResp saveTaskInfo(SaveTaskInfoReq saveTaskInfoReq);

    QueryResponse<SynthesisTaskInfoResp> selectSynthesisTaskInfoService(Integer page, Integer pageSize, SynthesisTaskInfoReq req);

    SynthesisTaskInfoFormTask selectSynthesisTaskInfoFormService(String id) throws Exception;
}
