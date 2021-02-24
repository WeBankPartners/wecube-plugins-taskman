package com.webank.taskman.controller;

import javax.validation.Valid;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.util.StringUtils;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import com.webank.taskman.base.JsonResponse;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.dto.TaskInfoDto;
import com.webank.taskman.dto.req.ProcessingTasksReq;
import com.webank.taskman.dto.req.QueryTaskInfoReq;
import com.webank.taskman.dto.req.QueryTemplateReq;
import com.webank.taskman.dto.req.TaskTemplateSaveReqDto;
import com.webank.taskman.dto.resp.RequestInfoInstanceResq;
import com.webank.taskman.dto.resp.TaskTemplateByRoleResp;
import com.webank.taskman.dto.resp.TaskTemplateResp;
import com.webank.taskman.service.TaskInfoService;
import com.webank.taskman.service.TaskTemplateService;


@RestController
@RequestMapping("/v1/task")
public class TaskmanTaskController {

    @Autowired
    private TaskTemplateService taskTemplateService;

    @Autowired
    private TaskInfoService taskInfoService;

    @PostMapping("/template/save")
    public JsonResponse taskTemplateSave(@Valid @RequestBody TaskTemplateSaveReqDto taskTemplateReq){

        TaskTemplateResp taskTemplateResp = taskTemplateService.saveTaskTemplateByReq(taskTemplateReq);
        return JsonResponse.okayWithData(taskTemplateResp);
    }

    @PostMapping("/template/search/{page}/{pageSize}")
    public JsonResponse taskTemplateSearch( @PathVariable("page") Integer page,
            @PathVariable("pageSize") Integer pageSize,
            @RequestBody QueryTemplateReq req) {
        QueryResponse<TaskTemplateByRoleResp> queryResponse = taskTemplateService.selectTaskTemplatePage(page, pageSize,
                req);
        return JsonResponse.okayWithData(queryResponse);
    }

    @GetMapping("/template/detail/{id}")
    public JsonResponse taskTemplateDetail(@PathVariable("id") String id){
        TaskTemplateResp taskTemplateResp = taskTemplateService.taskTemplateDetail(id);
        return JsonResponse.okayWithData(taskTemplateResp);
    }

    @PostMapping("/search/{page}/{pageSize}")
    public JsonResponse taskInfoSearch( @PathVariable("page") Integer page,
             @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) QueryTaskInfoReq req) {
        if (!StringUtils.isEmpty(req.getIsMy())) {
            req.setReporter(AuthenticationContextHolder.getCurrentUsername());
        }
        QueryResponse<TaskInfoDto> queryResponse = taskInfoService.selectTaskInfo(page, pageSize, req);
        return JsonResponse.okayWithData(queryResponse);
    }

    @PostMapping("/details")
    public JsonResponse taskInfoDetail(@RequestBody TaskInfoDto req) {
        return JsonResponse.okayWithData(taskInfoService.taskInfoDetail(req.getId()));
    }

    @PostMapping("/receive")
    public JsonResponse taskInfoReceive(@RequestBody TaskInfoDto req) {
        TaskInfoDto taskDTO = taskInfoService.taskInfoReceive(req.getId());
        if (null == taskDTO.getId()) {
            return JsonResponse.error("The task is not in an unclaimed state");
        }
        return JsonResponse.okayWithData(taskDTO);
    }

    @GetMapping("/instance")
    public JsonResponse taskInfoInstance(@RequestParam("requestId") String requestId,
            @RequestParam("taskId") String taskId) {
        RequestInfoInstanceResq requestInfoInstanceResq = taskInfoService.selectTaskInfoInstanceService(requestId,
                taskId);
        return JsonResponse.okayWithData(requestInfoInstanceResq);
    }

    @PostMapping("/processing")
    public JsonResponse taskInfoProcessing(@Valid @RequestBody ProcessingTasksReq req) {
        return taskInfoService.taskInfoProcessing(req);
    }

}
