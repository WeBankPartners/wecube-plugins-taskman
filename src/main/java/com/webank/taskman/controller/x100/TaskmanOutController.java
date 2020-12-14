package com.webank.taskman.controller.x100;


import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.WorkflowJsonResponse;
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
    public JsonResponse<List<CreateTaskServiceMetaResp>> taskCreateServiceMeta(
            @ApiParam(value = "proc-def-id",name = "procDefId",required = true) @PathVariable("proc-def-id") String procDefId,
            @ApiParam(value = "proc-def-id",name = "nodeDefId",required = true) @PathVariable("node-def-id") String nodeDefId)
    {
        List<CreateTaskServiceMetaResp> list = new LinkedList<>();
        list.add(new CreateTaskServiceMetaResp());
        return JsonResponse.okayWithData(new FormTemplateResp());
    }

    @ApiOperationSupport(order = 8)
    @PostMapping("/task/create")
    @ApiOperation(value = "create")
    public WorkflowJsonResponse createTask(@RequestBody CoreCreateTaskReq req)
    {
        List<FormItemTemplate> list = new LinkedList<>();
        return WorkflowJsonResponse.okay();
    }

    @ApiOperationSupport(order = 9)
    @PostMapping("/task/cancel")
    @ApiOperation(value = "cancel")
    public JsonResponse cancelTask(CoreCreateTaskReq req)
    {
        return JsonResponse.okay();
    }

}
