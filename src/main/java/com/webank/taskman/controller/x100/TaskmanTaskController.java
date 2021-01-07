package com.webank.taskman.controller.x100;


import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.base.JsonResponse;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.dto.TaskInfoDTO;
import com.webank.taskman.dto.req.ProcessingTasksReq;
import com.webank.taskman.dto.req.QueryTaskInfoReq;
import com.webank.taskman.dto.req.QueryTemplateReq;
import com.webank.taskman.dto.req.SaveTaskTemplateReq;
import com.webank.taskman.dto.resp.*;
import com.webank.taskman.service.TaskInfoService;
import com.webank.taskman.service.TaskTemplateService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.ApiParam;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.web.DefaultSecurityFilterChain;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;


@RestController
@RequestMapping("/v1/task")
@Api(tags = {"4、 Task inteface API"})
public class TaskmanTaskController {

    private static final Logger log = LoggerFactory.getLogger(TaskmanTaskController.class);

    @Autowired
    private TaskTemplateService taskTemplateService;

    @Autowired
    private TaskInfoService taskInfoService;


    @ApiOperationSupport(order = 1)
    @PostMapping("/template/save")
    @ApiOperation(value = "task-template-save", notes = "Need to pass in object: ")
    public JsonResponse taskTemplateSave(@Valid @RequestBody SaveTaskTemplateReq taskTemplateReq) throws Exception {

        TaskTemplateResp taskTemplateResp = taskTemplateService.saveTaskTemplateByReq(taskTemplateReq);
        return JsonResponse.okayWithData(taskTemplateResp);
    }

    @ApiOperationSupport(order = 2)
    @PostMapping("/template/search/{page}/{pageSize}")
    @ApiOperation(value = "task-template-search")
    public JsonResponse<QueryResponse<TaskTemplateResp>> taskTemplateSearch(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize") @PathVariable("pageSize") Integer pageSize,@RequestBody  QueryTemplateReq req){
        DefaultSecurityFilterChain D = null;
        QueryResponse<TaskTemplateByRoleResp> queryResponse = taskTemplateService.selectTaskTemplatePage(page,pageSize,req);
        return JsonResponse.okayWithData(queryResponse);
    }

    @ApiOperationSupport(order = 3)
    @GetMapping("/template/detail/{id}")
    @ApiOperation(value = "task-template-detail", notes = "需要传入id")
    public JsonResponse taskTemplateDetail(@PathVariable("id") String id) throws Exception {
        TaskTemplateResp taskTemplateResp = taskTemplateService.taskTemplateDetail(id);
        return JsonResponse.okayWithData(taskTemplateResp);
    }

    @ApiOperationSupport(order = 4)
    @PostMapping("/search/{page}/{pageSize}")
    @ApiOperation(value = "task-info-search")
    public JsonResponse<QueryResponse<TaskInfoDTO>> taskInfoSearch(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize")  @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) QueryTaskInfoReq req) {
        QueryResponse<TaskInfoDTO> queryResponse = taskInfoService.selectTaskInfo(page, pageSize,req);
        return JsonResponse.okayWithData(queryResponse);
    }

    @ApiOperationSupport(order =5)
    @PostMapping("/details")
    @ApiOperation(value = "task-info-detail")
    public JsonResponse<TaskInfoResp> taskInfoDetail(String id)
    {
        return JsonResponse.okayWithData(taskInfoService.taskInfoDetail(id));
    }


    @ApiOperationSupport(order =6)
    @PostMapping("/receive")
    @ApiOperation(value = "task-info-receive")
    public JsonResponse<TaskInfoDTO> taskInfoReceive(String id)
    {
        TaskInfoDTO taskDTO = taskInfoService.taskInfoReceive(id);
        if (null == taskDTO.getId()){
            return JsonResponse.customError("The task is not in an unclaimed state");
        }
        return JsonResponse.okayWithData(taskDTO);
    }


    @ApiOperationSupport(order =7)
    @GetMapping("/instance")
    @ApiOperation(value = "task-info-instance")
    public JsonResponse<RequestInfoInstanceResq> taskInfoInstance(
        @RequestParam("requestId") String requestId,@RequestParam("taskId") String taskId)
    {
        RequestInfoInstanceResq requestInfoInstanceResq = taskInfoService.selectTaskInfoInstanceService(requestId,taskId);
        return JsonResponse.okayWithData(requestInfoInstanceResq);
    }

    @ApiOperationSupport(order =8)
    @PostMapping("/processing")
    @ApiOperation(value = "task-info-processing")
    public JsonResponse taskInfoProcessing(@Valid @RequestBody ProcessingTasksReq req)
    {
        return taskInfoService.taskInfoProcessing(req);
    }

}

