package com.webank.taskman.controller.x100;


import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.base.JsonResponse;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.*;
import com.webank.taskman.dto.resp.TaskServiceMetaResp;
import com.webank.taskman.service.FormItemTemplateService;
import io.swagger.annotations.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;


import java.util.LinkedList;
import java.util.List;


@Api(tags = {"2、 Taskman open inteface API"})
@RestController
@RequestMapping("/v1")
public class TaskmanOutController {



    private static final Logger log = LoggerFactory.getLogger(TaskmanOutController.class);

    @Autowired
    FormItemTemplateService formItemTemplateService;

    @ApiOperationSupport(order = 1)
    @GetMapping("/task/create/service-meta/{proc-def-id}/{node-def-id}")
    @ApiOperation(value = "service-meta")
    public JsonResponse<TaskServiceMetaResp> taskCreateServiceMeta(
            @PathVariable("proc-def-id") String procDefId,@PathVariable("node-def-id") String nodeDefId)
    {
        return JsonResponse.okayWithData(formItemTemplateService.getTaskCreateServiceMeta(procDefId,nodeDefId));
    }

    @ApiOperationSupport(order = 2)
    @PostMapping("/task/create")
    @ApiOperation(value = "create")
    public CoreCreateTaskResp createTask(@RequestBody CoreCreateTaskDTO req)
    {
        List<FormItemTemplate> list = new LinkedList<>();
        return new CoreCreateTaskResp();
    }

    @ApiOperationSupport(order = 3)
    @PostMapping("/task/cancel")
    @ApiOperation(value = "cancel")
    public CoreCreateTaskResp cancelTask(@RequestBody CoreCancelTaskDTO req)
    {
        return new CoreCreateTaskResp();
    }

}
