package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.SynthesisRequestInfoReq;
import com.webank.taskman.dto.resp.SynthesisRequestInfoFormRequest;
import com.webank.taskman.dto.resp.SynthesisRequestInfoResp;
import com.webank.taskman.dto.resp.SynthesisRequestTempleResp;

import java.util.List;
import java.util.Map;

public interface RequestSynthesisService extends IService<RequestTemplate> {
    QueryResponse<SynthesisRequestTempleResp> selectSynthesisRequestTempleService(Integer current, Integer limit) throws Exception;

    QueryResponse<Map<String,Object>> selectSynthesisRequestInfoService(Integer current, Integer limit, SynthesisRequestInfoReq req) throws Exception;

    SynthesisRequestInfoFormRequest selectSynthesisRequestInfoFormService(String  id) throws Exception;

    QueryResponse<Map<String,Object>> selectSynthesisRequestInfoCurrencyResp(Integer current, Integer limit, SynthesisRequestInfoReq req) throws Exception;

}
