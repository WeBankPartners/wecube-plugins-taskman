package com.webank.taskman.controller.x100;

import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.QueryRequestTemplateReq;
import com.webank.taskman.dto.req.SynthesisRequestInfoReq;
import com.webank.taskman.dto.req.SynthesisTaskInfoReq;
import com.webank.taskman.dto.resp.*;
import com.webank.taskman.service.RequestSynthesisService;
import com.webank.taskman.service.TaskInfoService;
import com.webank.taskman.service.TaskTemplateService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.ApiParam;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@Api(tags = {"6、 RequestORTask Synthesis API"})
@RestController
@RequestMapping("/v1/Synthesis")
public class RequestSynthesisController {


    private static final Logger log = LoggerFactory.getLogger(RequestSynthesisController.class);

    @Autowired
    RequestSynthesisService requestSynthesisService;

    @Autowired
    TaskTemplateService taskTemplateService;

    @Autowired
    TaskInfoService taskInfoService;

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
    public JsonResponse<QueryResponse<SynthesisRequestInfoResp>>selectSynthesisRequestInfo(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize") @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) SynthesisRequestInfoReq req)
            throws Exception {
        QueryResponse<SynthesisRequestInfoResp> list = requestSynthesisService.selectSynthesisRequestInfoService(page, pageSize,req);
        return JsonResponse.okayWithData(list);
    }


    @ApiOperationSupport(order = 3)
    @PostMapping("/Request/selectSynthesisRequestInfoForm")
    @ApiOperation(value = "Synthesis-Request-Info-Form-search")
    public JsonResponse<SynthesisRequestInfoFormRequest> selectSynthesisRequestInfoForm(String id)
            throws Exception {
        SynthesisRequestInfoFormRequest synthesisRequestInfoFormRequest = requestSynthesisService.selectSynthesisRequestInfoFormService(id);
        return JsonResponse.okayWithData(synthesisRequestInfoFormRequest);
    }

    @ApiOperationSupport(order = 4)
    @PostMapping("/task/selectTaskSynthesis/{page}/{pageSize}")
    @ApiOperation(value = "task-Synthesis-search")
    public JsonResponse<QueryResponse<TaskTemplateByRoleResp>> selectTaskSynthesis(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize") @PathVariable("pageSize") Integer pageSize)
            throws Exception{
        QueryResponse<TaskTemplateByRoleResp> queryResponse = taskTemplateService.selectTaskTemplateByRole(page,pageSize);
        return JsonResponse.okayWithData(queryResponse);
    }

    @ApiOperationSupport(order = 5)
    @PostMapping("/task/selectSynthesisTaskInfo/{page}/{pageSize}")
    @ApiOperation(value = "Synthesis-Task-Info-search")
    public JsonResponse<QueryResponse<SynthesisTaskInfoResp>> selectSynthesisTaskInfo(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize")  @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) SynthesisTaskInfoReq req)
            throws Exception {
        QueryResponse<SynthesisTaskInfoResp> queryResponse = taskInfoService.selectSynthesisTaskInfoService(page, pageSize,req);
        return JsonResponse.okayWithData(queryResponse);
    }

    @ApiOperationSupport(order = 6)
    @PostMapping("/task/selectSynthesisTaskInfoForm")
    @ApiOperation(value = "Synthesis-Task-Info-Form-search")
    public JsonResponse<SynthesisTaskInfoFormTask> selectSynthesisTaskInfoForm(String id)
            throws Exception {
        SynthesisTaskInfoFormTask synthesisTaskInfoFormTask = taskInfoService.selectSynthesisTaskInfoFormService(id);
        return JsonResponse.okayWithData(synthesisTaskInfoFormTask);
    }






}
