package com.webank.taskman.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.webank.taskman.dto.CoreCancelTaskDto;
import com.webank.taskman.dto.CoreCreateTaskDto;
import com.webank.taskman.service.FormItemTemplateService;
import com.webank.taskman.service.TaskInfoService;
import com.webank.taskman.support.core.CommonResponseDto;

@RestController
@RequestMapping("/v1")
public class TaskmanOutController {

    @Autowired
    private FormItemTemplateService formItemTemplateService;

    @Autowired
    private TaskInfoService taskInfoService;

    
    //TODO
    @GetMapping("/task/create/service-meta/{proc-inst-id}/{node-def-id}")
    public CommonResponseDto taskCreateServiceMeta(@PathVariable("proc-inst-id") String procInstId,
            @PathVariable("node-def-id") String nodeDefId) {
        return CommonResponseDto.okayWithData(formItemTemplateService.getTaskCreateServiceMeta(procInstId, nodeDefId));
    }

    //TODO
    @PostMapping("/task/create")
    public CommonResponseDto taskCreate(@RequestBody CoreCreateTaskDto req) {
        return taskInfoService.createTask(req);
    }

    //TODO
    @PostMapping("/task/cancel")
    public CommonResponseDto taskCancel(@RequestBody CoreCancelTaskDto req) {
        return taskInfoService.cancelTask(req);
    }

}
