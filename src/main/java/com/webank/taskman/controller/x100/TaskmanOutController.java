package com.webank.taskman.controller.x100;


import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.FormItemTemplateDTO;
import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.WorkflowJsonResponse;
import com.webank.taskman.dto.req.SaveTaskInfoReq;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiImplicitParam;
import io.swagger.annotations.ApiImplicitParams;
import io.swagger.annotations.ApiOperation;
import org.springframework.web.bind.annotation.*;

import java.util.LinkedList;
import java.util.List;
import java.util.Map;

import static com.webank.taskman.dto.JsonResponse.okay;
import static com.webank.taskman.dto.JsonResponse.okayWithData;

@Api(tags = {"2、 Taskman open inteface API"})
@RestController
@RequestMapping("/v1")
public class TaskmanOutController {

    @ApiOperationSupport(order = 7)
    @GetMapping("/task/create/service-meta")
    @ApiOperation(value = "service-meta")
    @ApiImplicitParams({
            @ApiImplicitParam(name = "procDefId", value = "流程id",required =true, dataTypeClass = String.class),
            @ApiImplicitParam(name = "nodeDefId", value = "流程节点id",required = true,dataTypeClass = String.class),

    })
    public WorkflowJsonResponse queryTaskFormItemTemplateList(Map<String,Object> params)
    {
        List<FormItemTemplateDTO> list = new LinkedList<>();
        list.add(new FormItemTemplateDTO());

        return WorkflowJsonResponse.okayWithData(null);
    }

    // 创建任务
    @ApiOperationSupport(order = 8)
    @PostMapping("/task/create")
    @ApiOperation(value = "create")
    public WorkflowJsonResponse createTask(@RequestBody SaveTaskInfoReq req)
    {
        List<FormItemTemplate> list = new LinkedList<>();
        return WorkflowJsonResponse.okay();
    }

    // 取消任务
    @ApiOperationSupport(order = 9)
    @PostMapping("/task/cancel")
    @ApiOperation(value = "cancel")
    public WorkflowJsonResponse cancelTask(String procDefId,String procDefNodeId)
    {
        return WorkflowJsonResponse.okay();
    }

}
