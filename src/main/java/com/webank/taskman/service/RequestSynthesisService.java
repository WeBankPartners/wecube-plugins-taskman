package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.resp.SynthesisRequestInfoFormRequest;
import com.webank.taskman.dto.resp.SynthesisRequestInfoResp;
import com.webank.taskman.dto.resp.SynthesisRequestTempleResp;

public interface RequestSynthesisService extends IService<RequestTemplate> {
    QueryResponse<SynthesisRequestTempleResp> selectSynthesisRequestTempleService(Integer current, Integer limit) throws Exception;

    QueryResponse<SynthesisRequestInfoResp> selectSynthesisRequestInfoService(Integer current, Integer limit) throws Exception;

   SynthesisRequestInfoFormRequest selectSynthesisRequestInfoFormService(String  id) throws Exception;

}
