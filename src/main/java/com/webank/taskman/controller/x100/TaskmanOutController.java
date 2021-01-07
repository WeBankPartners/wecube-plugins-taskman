package com.webank.taskman.controller.x100;


import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.dto.CoreCancelTaskDTO;
import com.webank.taskman.dto.CoreCreateTaskDTO;
import com.webank.taskman.service.FormItemTemplateService;
import com.webank.taskman.service.TaskInfoService;
import com.webank.taskman.support.core.CommonResponseDto;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;


@Api(tags = {"2、 Taskman open inteface API"})
@RestController
@RequestMapping("/v1")
public class TaskmanOutController {



    private static final Logger log = LoggerFactory.getLogger(TaskmanOutController.class);

    @Autowired
    FormItemTemplateService formItemTemplateService;

    @ApiOperationSupport(order = 1)
    @GetMapping("/task/create/service-meta/{proc-inst-id}/{node-def-id}")
    @ApiOperation(value = "task-create-service-meta")
    public CommonResponseDto taskCreateServiceMeta(
            @PathVariable("proc-inst-id") String procInstId,@PathVariable("node-def-id") String nodeDefId)
    {
        return CommonResponseDto.okayWithData(formItemTemplateService.getTaskCreateServiceMeta(procInstId,nodeDefId));
    }

    @Autowired
    TaskInfoService taskInfoService;

    @ApiOperationSupport(order = 2)
    @PostMapping("/task/create")
    @ApiOperation(value = "task-create")
    public CommonResponseDto taskCreate(@RequestBody CoreCreateTaskDTO req)
    {
        return taskInfoService.createTask(req);
    }

    @ApiOperationSupport(order = 3)
    @PostMapping("/task/cancel")
    @ApiOperation(value = "task-cancel")
    public CommonResponseDto taskCancel(@RequestBody CoreCancelTaskDTO req)
    {
        return taskInfoService.cancelTask(req);
    }

}
