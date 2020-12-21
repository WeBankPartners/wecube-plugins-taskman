package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.domain.RequestInfo;
import com.webank.taskman.dto.req.DoneRequestReq;
import com.webank.taskman.dto.req.SaveRequestInfoReq;
import com.webank.taskman.dto.resp.RequestInfoResq;


public interface RequestInfoService extends IService<RequestInfo> {

    QueryResponse<RequestInfoResq> selectRequestInfoService
            (Integer current, Integer limit, SaveRequestInfoReq req) throws TaskmanRuntimeException;

    void doneServiceRequest(DoneRequestReq request);

    SaveRequestInfoReq saveRequestInfo(SaveRequestInfoReq req);
}
