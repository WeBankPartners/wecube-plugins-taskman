package com.webank.taskman.controller.x100;


import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.*;
import com.webank.taskman.dto.req.CoreCreateTaskReq;
import com.webank.taskman.dto.req.CreateTaskRequestDto;
import com.webank.taskman.dto.resp.CreateTaskServiceMetaResp;
import com.webank.taskman.dto.resp.FormTemplateResp;
import io.swagger.annotations.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.web.bind.annotation.*;


import java.util.LinkedList;
import java.util.List;


@Api(tags = {"2、 Taskman open inteface API"})
@RestController
@RequestMapping("/v1")
public class TaskmanOutController {



    private static final Logger log = LoggerFactory.getLogger(TaskmanOutController.class);


    @ApiOperationSupport(order = 7)
    @GetMapping("/task/create/service-meta/{proc-def-id}/{node-def-id}")
    @ApiOperation(value = "service-meta")
    public JsonResponse<CreateTaskServiceMetaResp> taskCreateServiceMeta(
            @ApiParam(value = "proc-def-id",name = "procDefId",required = true) @PathVariable("proc-def-id") String procDefId,
            @ApiParam(value = "node-def-id",name = "nodeDefId",required = true) @PathVariable("node-def-id") String nodeDefId)
    {
        return JsonResponse.okayWithData(new CreateTaskServiceMetaResp());
    }

    @ApiOperationSupport(order = 8)
    @PostMapping("/task/create")
    @ApiOperation(value = "create")
    public CoreCreateTaskResp createTask(@RequestBody CoreCreateTaskDTO req)
    {
        List<FormItemTemplate> list = new LinkedList<>();
        return new CoreCreateTaskResp();
    }

    @ApiOperationSupport(order = 9)
    @PostMapping("/task/cancel")
    @ApiOperation(value = "cancel")
    public CoreCreateTaskResp cancelTask(@RequestBody CoreCancelTaskDTO req)
    {
        return new CoreCreateTaskResp();
    }

}
