package com.webank.taskman.controller.x100;

import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.QueryRequestTemplateReq;
import com.webank.taskman.dto.req.SynthesisRequestInfoReq;
import com.webank.taskman.dto.resp.SynthesisRequestInfoFormRequest;
import com.webank.taskman.dto.resp.SynthesisRequestInfoResp;
import com.webank.taskman.dto.resp.SynthesisRequestTempleResp;
import com.webank.taskman.service.RequestSynthesisService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.ApiParam;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@Api(tags = {"6、 Request Synthesis API"})
@RestController
@RequestMapping("/v1/Synthesis")
public class RequestSynthesisController {
    @Autowired
    RequestSynthesisService requestSynthesisService;


    @ApiOperationSupport(order = 1)
    @PostMapping("/Request/selectRequestSynthesis/{page}/{pageSize}")
    @ApiOperation(value = "request-Synthesis-search")
    public JsonResponse<QueryResponse<SynthesisRequestTempleResp>> selectRequestSynthesis(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize") @PathVariable("pageSize") Integer pageSize)
            throws Exception {
        QueryResponse<SynthesisRequestTempleResp> queryResponse = requestSynthesisService.selectSynthesisRequestTempleService(page, pageSize);
        return JsonResponse.okayWithData(queryResponse);
    }

    @ApiOperationSupport(order = 2)
    @PostMapping("/Request/selectSynthesisRequestInfo/{page}/{pageSize}")
    @ApiOperation(value = "Synthesis-Request-Info-search")
    public JsonResponse<QueryResponse<SynthesisRequestInfoResp>> selectSynthesisRequestInfo(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize") @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) SynthesisRequestInfoReq req)
            throws Exception {
        QueryResponse<SynthesisRequestInfoResp> queryResponse = requestSynthesisService.selectSynthesisRequestInfoService(page, pageSize,req);
        return JsonResponse.okayWithData(queryResponse);
    }


    @ApiOperationSupport(order = 3)
    @PostMapping("/Request/selectSynthesisRequestInfoForm")
    @ApiOperation(value = "Synthesis-Request-Info-Form-search")
    public JsonResponse<SynthesisRequestInfoFormRequest> selectSynthesisRequestInfoForm(String id)
            throws Exception {
        SynthesisRequestInfoFormRequest synthesisRequestInfoFormRequest = requestSynthesisService.selectSynthesisRequestInfoFormService(id);
        return JsonResponse.okayWithData(synthesisRequestInfoFormRequest);
    }
}
