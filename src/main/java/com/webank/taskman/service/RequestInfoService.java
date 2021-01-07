package com.webank.taskman.service;


import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.domain.RequestInfo;
import com.webank.taskman.dto.CreateTaskDto;
import com.webank.taskman.dto.req.QueryRequestInfoReq;
import com.webank.taskman.dto.req.SaveRequestInfoReq;
import com.webank.taskman.dto.resp.RequestInfoResq;
import com.webank.taskman.support.core.dto.DynamicWorkflowInstCreationInfoDto;
import com.webank.taskman.support.core.dto.DynamicWorkflowInstInfoDto;

public interface RequestInfoService extends IService<RequestInfo> {

    QueryResponse<RequestInfoResq> selectRequestInfoPage
            (Integer current, Integer limit, QueryRequestInfoReq req) throws TaskmanRuntimeException;

    RequestInfoResq saveRequestInfoByDto(CreateTaskDto req);

    DynamicWorkflowInstInfoDto createNewWorkflowInstance(CreateTaskDto req);


    RequestInfoResq selectDetail(String id);
}
