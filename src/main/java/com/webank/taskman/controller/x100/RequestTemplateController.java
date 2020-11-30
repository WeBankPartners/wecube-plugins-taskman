package com.webank.taskman.controller.x100;


import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.TemplateGroupDTO;
import com.webank.taskman.dto.TemplateGroupVO;
import com.webank.taskman.service.RequestTemplateService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.DynamicParameter;
import io.swagger.annotations.DynamicParameters;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.Map;


@RestController
@RequestMapping("/v1/request-template")
@Api(tags = {"V1.0.0 RequestTemplate model"})
public class RequestTemplateController {

    @Autowired
    RequestTemplateService requestTemplateService;


    @PostMapping("/add")
    @ApiOperation(value = "add RequestTemplate",notes = "Need to pass in object: ")
    @DynamicParameters(name = "req",properties = {
            @DynamicParameter(name = "page",value = "页码",example = "",required = true,dataTypeClass = Integer.class),
            @DynamicParameter(name = "pageSize",value = "每页行数",example = "100",required = true,dataTypeClass = Integer.class),
            @DynamicParameter(name = "id",value = "主键",example = "",dataTypeClass = String.class),
            @DynamicParameter(name = "name",value = "模板组名称"),
            @DynamicParameter(name = "manageRoleId",value = "管理角色Id"),
    })
    public JsonResponse createTemplateGroup(@RequestBody TemplateGroupVO templateGroupVO) throws Exception {
        return JsonResponse.okay();
    }

    @PostMapping("edit")
    @ApiOperation(value = "edit RequestTemplate",notes = "")
    public JsonResponse updateTemplateGroup(
            @RequestBody TemplateGroupVO templateGroupVO) throws Exception {
        return JsonResponse.okay();
    }

    @PostMapping("/search")
    @ApiOperation(value = "search RequestTemplate ")
    @DynamicParameters(name = "req",properties = {
            @DynamicParameter(name = "page",value = "页码",example = "",required = true,dataTypeClass = Integer.class),
            @DynamicParameter(name = "pageSize",value = "每页行数",example = "100",required = true,dataTypeClass = Integer.class),
            @DynamicParameter(name = "id",value = "主键",example = "",dataTypeClass = String.class),
            @DynamicParameter(name = "name",value = "模板名称"),
    })
    public JsonResponse<QueryResponse<TemplateGroupDTO>> selectTemplateGroup(@RequestBody Map<String,Object> req) throws Exception {
        return JsonResponse.okayWithData(null);
    }

    @DeleteMapping("/delete/{id}")
    @ApiOperation(value = "delete RequestTemplate",notes = "需要传入id")
    public JsonResponse deleteTemplateGroupByID(@PathVariable("id") String id) throws Exception {
        return JsonResponse.okay();
    }
}

