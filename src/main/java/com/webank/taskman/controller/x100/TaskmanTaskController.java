package com.webank.taskman.controller.x100;


import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.base.JsonResponse;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.dto.TaskInfoDTO;
import com.webank.taskman.dto.req.ProcessingTasksReq;
import com.webank.taskman.dto.req.QueryTaskInfoReq;
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
    @ApiOperation(value = "Task-Template-save", notes = "Need to pass in object: ")
    public JsonResponse createTaskTemplate(@Valid @RequestBody SaveTaskTemplateReq taskTemplateReq) throws Exception {

        TaskTemplateResp taskTemplateResp = taskTemplateService.saveTaskTemplateByReq(taskTemplateReq);
        return JsonResponse.okayWithData(taskTemplateResp);
    }

    @ApiOperationSupport(order = 2)
    @PostMapping("/template/search/{page}/{pageSize}")
    @ApiOperation(value = "Task-Template-selectAll")
    public JsonResponse<QueryResponse<TaskTemplateByRoleResp>> selectTaskSynthesis(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize") @PathVariable("pageSize") Integer pageSize){
        QueryResponse<TaskTemplateByRoleResp> queryResponse = taskTemplateService.selectTaskTemplateByRole(page,pageSize);
        return JsonResponse.okayWithData(queryResponse);
    }

    @ApiOperationSupport(order = 3)
    @GetMapping("/template/detail/{id}")
    @ApiOperation(value = "Task-Template-detail", notes = "需要传入id")
    public JsonResponse detail(@PathVariable("id") String id) throws Exception {
        TaskTemplateResp taskTemplateResp = taskTemplateService.selectTaskTemplateOne(id);
        return JsonResponse.okayWithData(taskTemplateResp);
    }

    @ApiOperationSupport(order = 4)
    @PostMapping("/search/{page}/{pageSize}")
    @ApiOperation(value = "task-Info-search")
    public JsonResponse<QueryResponse<TaskInfoDTO>> selectTaskInfo(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize")  @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) QueryTaskInfoReq req) {
        QueryResponse<TaskInfoDTO> queryResponse = taskInfoService.selectTaskInfo(page, pageSize,req);
        return JsonResponse.okayWithData(queryResponse);
    }


    @ApiOperationSupport(order =5)
    @PostMapping("/details")
    @ApiOperation(value = "Task-Info-details")
    public JsonResponse<SynthesisTaskInfoFormTask> selectSynthesisTaskInfoForm(String id)
            throws Exception {
        SynthesisTaskInfoFormTask synthesisTaskInfoFormTask = taskInfoService.selectSynthesisTaskInfoFormService(id);
        return JsonResponse.okayWithData(synthesisTaskInfoFormTask);
    }


    @ApiOperationSupport(order =6)
    @PostMapping("/receive")
    @ApiOperation(value = "Task-Info-receive")
    public JsonResponse<TaskInfoGetResp> getTheTaskInfo(String id)
            throws Exception {
        TaskInfoGetResp taskInfoGetResp = taskInfoService.getTheTaskInfoService(id);
        if (taskInfoGetResp.getId()==null){
            return JsonResponse.customError("The task is not in an unclaimed state");
        }
        return JsonResponse.okayWithData(taskInfoGetResp);
    }


    @ApiOperationSupport(order =7)
    @PostMapping("/instance")
    @ApiOperation(value = "Task-Info-instance")
    public JsonResponse<RequestInfoInstanceResq> selectTaskInfoinstance(@RequestParam("taskId") String taskId,
                                                                        @RequestParam("requestId") String requestId)
            throws Exception {
        RequestInfoInstanceResq requestInfoInstanceResq = taskInfoService.selectTaskInfoInstanceService(taskId,requestId);
        return JsonResponse.okayWithData(requestInfoInstanceResq);
    }

    @ApiOperationSupport(order =8)
    @PostMapping("/processing")
    @ApiOperation(value = "Task-Info-processing")
    public JsonResponse<String> ProcessingTasksController(@Valid @RequestBody ProcessingTasksReq ptr)
            throws Exception {
        String msg=taskInfoService.ProcessingTasksService(ptr);
        return JsonResponse.okayWithData(msg);
    }

}

