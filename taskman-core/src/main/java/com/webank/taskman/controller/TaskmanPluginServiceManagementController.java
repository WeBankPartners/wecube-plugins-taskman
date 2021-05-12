package com.webank.taskman.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import com.webank.taskman.dto.platform.CoreCancelTaskDto;
import com.webank.taskman.service.TaskInfoService;
import com.webank.taskman.service.TaskmanPluginServiceManagementService;
import com.webank.taskman.support.platform.dto.CommonPlatformResponseDto;
import com.webank.taskman.support.platform.dto.PlatformPluginRequestDto;
import com.webank.taskman.support.platform.dto.PlatformPluginResponseDto;
import com.webank.taskman.support.platform.dto.TaskFormMetaDto;
import com.webank.taskman.support.platform.dto.TaskFormMetaResponseDto;

@RestController
@RequestMapping("/v1")
public class TaskmanPluginServiceManagementController {

    @Autowired
    private TaskInfoService taskInfoService;
    
    @Autowired
    private TaskmanPluginServiceManagementService taskmanPluginService;

    
    @GetMapping("/task/create/meta")
    public TaskFormMetaResponseDto fetchTaskCreationMeta(@RequestParam("procInstId") String procInstId,
            @RequestParam("nodeDefId") String nodeDefId) {
        TaskFormMetaDto taskFormMetaDto = taskmanPluginService.fetchTaskCreationMeta(procInstId, nodeDefId);
        
        TaskFormMetaResponseDto respDto = new TaskFormMetaResponseDto();
        respDto.setData(taskFormMetaDto);
        respDto.setStatus(TaskFormMetaResponseDto.STATUS_OK);
        return respDto;
    }

    @PostMapping("/task/create")
    public PlatformPluginResponseDto createTask(@RequestBody PlatformPluginRequestDto platformPluginRequestDto) {
        
        taskmanPluginService.createTask(platformPluginRequestDto);
        
        PlatformPluginResponseDto respDto = new PlatformPluginResponseDto();
        respDto.setResultCode(PlatformPluginResponseDto.RESULT_CODE_OK);
        
        return respDto;
    }

    //TODO 
    @PostMapping("/task/cancel")
    public CommonPlatformResponseDto taskCancel(@RequestBody CoreCancelTaskDto req) {
        return taskInfoService.cancelTask(req);
    }

}
