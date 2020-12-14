package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.SaveTaskInfoAndFormInfoReq;
import com.webank.taskman.dto.req.SelectTaskInfoReq;
import com.webank.taskman.dto.req.SynthesisTaskInfoReq;
import com.webank.taskman.dto.resp.SaveTaskInfoResp;
import com.webank.taskman.dto.resp.SynthesisTaskInfoFormTask;
import com.webank.taskman.dto.resp.SynthesisTaskInfoResp;
import com.webank.taskman.dto.resp.TaskInfoResp;


public interface TaskInfoService extends IService<TaskInfo> {

    QueryResponse<TaskInfoResp> selectTaskInfoService(Integer page, Integer pageSize, SelectTaskInfoReq req);

    SaveTaskInfoResp saveTaskInfo(SaveTaskInfoAndFormInfoReq saveTaskInfoAndFormInfoReq);

    QueryResponse<SynthesisTaskInfoResp> selectSynthesisTaskInfoService(Integer page, Integer pageSize, SynthesisTaskInfoReq req);

    SynthesisTaskInfoFormTask selectSynthesisTaskInfoFormService(String id) throws Exception;
}
