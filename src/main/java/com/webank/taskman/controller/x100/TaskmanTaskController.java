package com.webank.taskman.controller.x100;


import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.dto.*;
import com.webank.taskman.dto.req.SelectTaskInfoReq;
import com.webank.taskman.dto.req.SaveTaskTemplateReq;
import com.webank.taskman.dto.req.SynthesisTaskInfoReq;
import com.webank.taskman.dto.resp.*;
import com.webank.taskman.service.TaskInfoService;
import com.webank.taskman.service.TaskTemplateService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.ApiParam;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.validation.BindingResult;
import org.springframework.validation.ObjectError;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;
import java.util.List;


@RestController
@RequestMapping("/v1/task")
@Api(tags = {"4、 Task inteface API"})
public class TaskmanTaskController {



    private static final Logger log = LoggerFactory.getLogger(TaskmanTaskController.class);


    @Autowired
    private TaskTemplateService taskTemplateService;

    @Autowired
    private TaskInfoService taskInfoService;

    @ApiOperationSupport(order = 18)
    @PostMapping("/template/save")
    @ApiOperation(value = "Task-Template-save", notes = "Need to pass in object: ")
    public JsonResponse createTaskTemplate(@Valid @RequestBody SaveTaskTemplateReq taskTemplateReq) throws Exception {

        TaskTemplateResp taskTemplateResp = taskTemplateService.saveTaskTemplateByReq(taskTemplateReq);
        return JsonResponse.okayWithData(taskTemplateResp);
    }

    @ApiOperationSupport(order = 19)
    @PostMapping("/template/search/{page}/{pageSize}")
    @ApiOperation(value = "Task-Template-selectAll")
    public JsonResponse<QueryResponse<TaskTemplateByRoleResp>> selectTaskSynthesis(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize") @PathVariable("pageSize") Integer pageSize)
            throws Exception{
        QueryResponse<TaskTemplateByRoleResp> queryResponse = taskTemplateService.selectTaskTemplateByRole(page,pageSize);
        return JsonResponse.okayWithData(queryResponse);
    }

    @ApiOperationSupport(order = 20)
    @GetMapping("/template/detail/{id}")
    @ApiOperation(value = "Task-Template-detail", notes = "需要传入id")
    public JsonResponse detail(@PathVariable("id") String id) throws Exception {
        TaskTemplateResp taskTemplateResp = taskTemplateService.selectTaskTemplateOne(id);
        return JsonResponse.okayWithData(taskTemplateResp);
    }

    @ApiOperationSupport(order = 21)
    @PostMapping("/task/search/{page}/{pageSize}")
    @ApiOperation(value = "Task-Info-search")
    public JsonResponse<QueryResponse<SynthesisTaskInfoResp>> selectSynthesisTaskInfo(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize")  @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) SynthesisTaskInfoReq req)
            throws Exception {
        QueryResponse<SynthesisTaskInfoResp> queryResponse = taskInfoService.selectSynthesisTaskInfoService(page, pageSize,req);
        return JsonResponse.okayWithData(queryResponse);
    }

    @ApiOperationSupport(order =22)
    @PostMapping("/task/details")
    @ApiOperation(value = "Task-Info-details")
    public JsonResponse<SynthesisTaskInfoFormTask> selectSynthesisTaskInfoForm(String id)
            throws Exception {
        SynthesisTaskInfoFormTask synthesisTaskInfoFormTask = taskInfoService.selectSynthesisTaskInfoFormService(id);
        return JsonResponse.okayWithData(synthesisTaskInfoFormTask);
    }


}
