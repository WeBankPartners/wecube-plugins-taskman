package com.webank.taskman.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import com.webank.taskman.dto.platform.CoreCancelTaskDto;
import com.webank.taskman.service.FormItemTemplateService;
import com.webank.taskman.service.TaskInfoService;
import com.webank.taskman.support.platform.dto.CommonPlatformResponseDto;
import com.webank.taskman.support.platform.dto.PlatformPluginRequestDto;
import com.webank.taskman.support.platform.dto.PlatformPluginResponseDto;
import com.webank.taskman.support.platform.dto.TaskFormMetaResponseDto;

@RestController
@RequestMapping("/v1")
public class TaskmanOutController {

    @Autowired
    private FormItemTemplateService formItemTemplateService;

    @Autowired
    private TaskInfoService taskInfoService;

    
    //TODO
    @GetMapping("/task/create/meta")
    public TaskFormMetaResponseDto taskCreateServiceMeta(@RequestParam("procInstId") String procInstId,
            @RequestParam("nodeDefId") String nodeDefId) {
//        return CommonResponseDto.okayWithData(formItemTemplateService.getTaskCreateServiceMeta(procInstId, nodeDefId));
        
        return null;
    }

    //TODO
    @PostMapping("/task/create")
    public PlatformPluginResponseDto taskCreate(@RequestBody PlatformPluginRequestDto platformPluginRequestDto) {
        //TODO #48
//        return taskInfoService.createTask(req);
        
        return null;
    }

    //TODO
    @PostMapping("/task/cancel")
    public CommonPlatformResponseDto taskCancel(@RequestBody CoreCancelTaskDto req) {
        return taskInfoService.cancelTask(req);
    }

}
