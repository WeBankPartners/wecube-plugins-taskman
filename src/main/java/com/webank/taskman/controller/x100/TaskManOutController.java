package com.webank.taskman.controller.x100;


import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.FormItemTemplateDTO;
import com.webank.taskman.dto.JsonResponse;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.DynamicParameter;
import io.swagger.annotations.DynamicParameters;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.LinkedList;
import java.util.List;
import java.util.Map;

import static com.webank.taskman.dto.JsonResponse.okay;
import static com.webank.taskman.dto.JsonResponse.okayWithData;

@Api(tags = {"4、 process model"})
@RestController
@RequestMapping("/v1/taskman")
public class TaskManOutController {

    // 获取任务表单模板
    @GetMapping("/task/form/item-templates")
    @ApiOperation(value = "item-templates")
    @DynamicParameters(name = "params", properties = {
            @DynamicParameter(name = "procDefId", value = "流程id",required = true, dataTypeClass = String.class),
            @DynamicParameter(name = "procDefNode", value = "流程节点id",required = true,dataTypeClass = String.class),
    })
    public JsonResponse<List<FormItemTemplateDTO>> queryTaskFormItemTemplateList(Map<String,Object> params) {
        List<FormItemTemplateDTO> list = new LinkedList<>();
        list.add(new FormItemTemplateDTO());
        return okayWithData(list);
    }

    // 创建任务
    @GetMapping("/task/create")
    @ApiOperation(value = "create")
    public JsonResponse createTask() {
        List<FormItemTemplate> list = new LinkedList<>();
        return okay();
    }

    // 取消任务
    @GetMapping("/task/cancel")
    @ApiOperation(value = "cancel")
    public JsonResponse cancelTask() {
        List<FormItemTemplate> list = new LinkedList<>();
        return okay();
    }

}
