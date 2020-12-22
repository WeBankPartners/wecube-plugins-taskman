package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.domain.RequestInfo;
import com.webank.taskman.dto.req.SaveRequestInfoReq;
import com.webank.taskman.dto.req.SynthesisRequestInfoReq;
import com.webank.taskman.dto.resp.RequestInfoResq;
import com.webank.taskman.dto.resp.SynthesisRequestInfoFormRequest;
import com.webank.taskman.dto.resp.SynthesisRequestInfoResp;


public interface RequestInfoService extends IService<RequestInfo> {

    QueryResponse<RequestInfoResq> selectRequestInfoService
            (Integer current, Integer limit, SaveRequestInfoReq req) throws TaskmanRuntimeException;

    SaveRequestInfoReq saveRequestInfo(SaveRequestInfoReq req);

    SynthesisRequestInfoFormRequest selectSynthesisRequestInfoFormService(String id) throws TaskmanRuntimeException;

    QueryResponse<SynthesisRequestInfoResp> selectSynthesisRequestInfoService
            (Integer current, Integer limit, SynthesisRequestInfoReq req) throws TaskmanRuntimeException;
}
