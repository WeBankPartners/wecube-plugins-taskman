package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.domain.RequestInfo;
import com.webank.taskman.dto.req.SaveRequestInfoReq;
import com.webank.taskman.dto.req.SynthesisRequestInfoReq;
import com.webank.taskman.dto.resp.RequestInfoResq;
import com.webank.taskman.dto.resp.SynthesisRequestInfoForm;
import com.webank.taskman.dto.resp.SynthesisRequestInfoFormRequest;
import com.webank.taskman.dto.resp.SynthesisRequestInfoResp;
import com.webank.taskman.support.core.dto.DynamicWorkflowInstCreationInfoDto;
import com.webank.taskman.support.core.dto.DynamicWorkflowInstInfoDto;


public interface RequestInfoService extends IService<RequestInfo> {

    QueryResponse<RequestInfoResq> selectRequestInfoService
            (Integer current, Integer limit, SaveRequestInfoReq req) throws TaskmanRuntimeException;

    RequestInfoResq saveRequestInfo(SaveRequestInfoReq req);

    SynthesisRequestInfoForm selectSynthesisRequestInfoFormService(String id) throws TaskmanRuntimeException;

    QueryResponse<SynthesisRequestInfoResp> selectSynthesisRequestInfoService
            (Integer current, Integer limit, SynthesisRequestInfoReq req) throws TaskmanRuntimeException;

    DynamicWorkflowInstInfoDto createNewWorkflowInstance(RequestInfo requestInfo);

    DynamicWorkflowInstCreationInfoDto createDynamicWorkflowInstCreationInfoDto(String procDefId, String guid);

}
