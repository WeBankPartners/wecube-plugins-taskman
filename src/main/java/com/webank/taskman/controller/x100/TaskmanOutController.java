package com.webank.taskman.controller.x100;


import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.FormItemTemplateDTO;
import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.WorkflowJsonResponse;
import com.webank.taskman.dto.req.CoreTaskCreateServiceMetaReq;
import com.webank.taskman.dto.req.CreateTaskRequestDto;
import com.webank.taskman.dto.req.SaveTaskInfoReq;
import com.webank.taskman.dto.resp.FormTemplateResp;
import io.swagger.annotations.*;
import org.springframework.web.bind.annotation.*;
import springfox.documentation.annotations.ApiIgnore;

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
    public JsonResponse<FormTemplateResp> queryTaskFormItemTemplateList(CoreTaskCreateServiceMetaReq req)
    {
        List<FormItemTemplateDTO> list = new LinkedList<>();
        list.add(new FormItemTemplateDTO());

        return JsonResponse.okayWithData(new FormTemplateResp());
    }

    @ApiOperationSupport(order = 8)
    @PostMapping("/task/create")
    @ApiOperation(value = "create")
    public WorkflowJsonResponse createTask(@RequestBody CreateTaskRequestDto req)
    {
        List<FormItemTemplate> list = new LinkedList<>();
        return WorkflowJsonResponse.okay();
    }

    @ApiOperationSupport(order = 9)
    @PostMapping("/task/cancel")
    @ApiOperation(value = "cancel")
    public JsonResponse cancelTask(CoreTaskCreateServiceMetaReq req)
    {
        return JsonResponse.okay();
    }

}
