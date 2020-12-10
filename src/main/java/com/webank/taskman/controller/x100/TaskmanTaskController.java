package com.webank.taskman.controller.x100;


import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.dto.*;
import com.webank.taskman.dto.req.SaveTaskInfoReq;
import com.webank.taskman.dto.req.SaveTaskTemplateReq;
import com.webank.taskman.dto.resp.TaskInfoResp;
import com.webank.taskman.dto.resp.TaskTemplateResp;
import com.webank.taskman.service.TaskInfoService;
import com.webank.taskman.service.TaskTemplateService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.ApiParam;
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

    @Autowired
    private TaskTemplateService taskTemplateService;

    @Autowired
    private TaskInfoService taskInfoService;

    @ApiOperationSupport(order = 18)
    @PostMapping("/template/save")
    @ApiOperation(value = "Task-Template-save", notes = "Need to pass in object: ")
    public JsonResponse createTaskTemplate(@Valid @RequestBody SaveTaskTemplateReq taskTemplateReq, BindingResult bindingResult) throws Exception {
        if (bindingResult.hasErrors()) {
            for (ObjectError error : bindingResult.getAllErrors()) {
                return JsonResponse.okayWithData(error.getDefaultMessage());
            }
        }
        TaskTemplateResp taskTemplateResp = taskTemplateService.saveTaskTemplateByReq(taskTemplateReq);
        return JsonResponse.okayWithData(taskTemplateResp);
    }

    @ApiOperationSupport(order = 19)
    @GetMapping("/template/search")
    @ApiOperation(value = "Task-Template-selectAll")
    public JsonResponse selectTaskTemplateAll() throws Exception {
        List<TaskTemplateResp> taskTemplateRespList = taskTemplateService.selectTaskTemplateAll();
        return JsonResponse.okayWithData(taskTemplateRespList);
    }

    @ApiOperationSupport(order = 20)
    @GetMapping("/template/detail/{id}")
    @ApiOperation(value = "Task-Template-detail", notes = "需要传入id")
    public JsonResponse detail(@PathVariable("id") String id) throws Exception {
        TaskTemplateResp taskTemplateResp = taskTemplateService.selectTaskTemplateOne(id);
        return JsonResponse.okayWithData(taskTemplateResp);
    }

    @ApiOperationSupport(order = 21)
    @PostMapping("/search/{page}/{pageSize}")
    @ApiOperation(value = "Task-Info-search")
    public JsonResponse<QueryResponse<TaskInfoResp>> selectRequestInfo(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize") @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) SaveTaskInfoReq req)
            throws Exception {
        QueryResponse<TaskInfoResp> queryResponse = taskInfoService.selectTaskInfoService(page, pageSize, req);
        return JsonResponse.okayWithData(queryResponse);
    }

}

