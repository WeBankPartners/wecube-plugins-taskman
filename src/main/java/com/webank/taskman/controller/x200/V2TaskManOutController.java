package com.webank.taskman.controller.x200;


import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.FormItemTemplateDTO;
import com.webank.taskman.dto.JsonResponse;
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
@RequestMapping("/v2/taskman")
public class V2TaskManOutController {

    // 获取任务表单模板
    @GetMapping("/task/form/item/templates")
    @ApiOperation(value = "query from item template list")
    @ApiImplicitParams({
            @ApiImplicitParam(name = "procDefId", value = "流程id",required = true, dataTypeClass = String.class),
            @ApiImplicitParam(name = "procDefNode", value = "流程节点id",required = true,dataTypeClass = String.class),

    })
    public JsonResponse<List<FormItemTemplateDTO>> queryTaskmanFormItemTemplateList(Map<String,Object> params)
    {
        List<FormItemTemplateDTO> list = new LinkedList<>();
        list.add(new FormItemTemplateDTO());
        return okayWithData(list);
    }

    // 创建任务
    @GetMapping("/task/create")
    @ApiOperation(value = "create")
    public JsonResponse createTask(@RequestBody SaveTaskInfoReq req)
    {
        List<FormItemTemplate> list = new LinkedList<>();
        return okay();
    }

    // 取消任务
    @GetMapping("/task/cancel/{procDefId}/{procDefNodeId}")
    @ApiOperation(value = "cancel")
    public JsonResponse cancelTask( @PathVariable("procDefId") Integer procDefId,
            @PathVariable("procDefNodeId") Integer procDefNodeId)
    {
        return okay();
    }

}
