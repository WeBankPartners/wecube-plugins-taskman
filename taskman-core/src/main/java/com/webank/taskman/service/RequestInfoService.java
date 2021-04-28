package com.webank.taskman.service;


import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.LocalPageableQueryResult;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.domain.RequestInfo;
import com.webank.taskman.dto.CreateTaskDto;
import com.webank.taskman.dto.req.RequestInfoQueryReqDto;
import com.webank.taskman.dto.resp.RequestInfoQueryResultDto;
import com.webank.taskman.support.core.dto.DynamicWorkflowInstInfoDto;

public interface RequestInfoService extends IService<RequestInfo> {

    LocalPageableQueryResult<RequestInfoQueryResultDto> searchRequestInfos
            (Integer current, Integer limit, RequestInfoQueryReqDto req) throws TaskmanRuntimeException;

    RequestInfoQueryResultDto createNewRequestInfo(CreateTaskDto req);

    DynamicWorkflowInstInfoDto createNewRemoteWorkflowInstance(CreateTaskDto req);


    RequestInfoQueryResultDto fetchRequestInfoDetail(String id);
}
