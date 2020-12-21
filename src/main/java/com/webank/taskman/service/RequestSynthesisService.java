package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.req.SynthesisRequestInfoReq;
import com.webank.taskman.dto.resp.SynthesisRequestInfoFormRequest;
import com.webank.taskman.dto.resp.SynthesisRequestInfoResp;

public interface RequestSynthesisService extends IService<RequestTemplate> {


    QueryResponse<SynthesisRequestInfoResp> selectSynthesisRequestInfoService
            (Integer current, Integer limit, SynthesisRequestInfoReq req) throws TaskmanRuntimeException;

    SynthesisRequestInfoFormRequest selectSynthesisRequestInfoFormService(String  id) throws TaskmanRuntimeException;

}
