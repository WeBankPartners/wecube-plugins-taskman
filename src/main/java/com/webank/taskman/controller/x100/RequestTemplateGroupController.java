package com.webank.taskman.controller.x100;


import com.webank.taskman.dto.*;
import com.webank.taskman.service.RequestTemplateGroupService;
import io.swagger.annotations.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import springfox.documentation.annotations.ApiIgnore;

import java.util.Map;

@RestController
@RequestMapping("/v1/request-template-group")
@Api(tags = {"V1.0.0 TemplateGroup model"})
public class RequestTemplateGroupController {
    @Autowired
    RequestTemplateGroupService requestTemplateGroupService;

    @PostMapping("/add")
    @ApiOperation(value = "add RequestTemplateGroup", notes = "")
    public JsonResponse createTemplateGroup(
            @RequestBody TemplateGroupCreateVO templateGroupCreateVO) throws Exception {
        requestTemplateGroupService.createTemplateGroupService(templateGroupCreateVO);
        return JsonResponse.okay();
    }

    @PostMapping("edit")
    @ApiOperation(value = "edit RequestTemplateGroup", notes = "")
    public JsonResponse updateTemplateGroup(
            @RequestBody TemplateGroupVO templateGroupVO) throws Exception {
        requestTemplateGroupService.updateTemplateGroupService(templateGroupVO);
        return JsonResponse.okay();
    }

    @PostMapping("/search/{current}/{limit}")
    @ApiOperation(value = "search RequestTemplateGroup ")
    @DynamicParameters(name = "req", properties = {
            @DynamicParameter(name = "id", value = "主键", example = "", dataTypeClass = String.class),
            @DynamicParameter(name = "name", value = "模板组名称"),
            @DynamicParameter(name = "manageRoleId", value = "管理角色Id"),
    })
    public JsonResponse<QueryResponse<TemplateGroupDTO>> selectTemplateGroup(
            @PathVariable("current") Integer current,
            @PathVariable("limit") Integer limit,
            @RequestBody(required = false) TemplateGroupReq req
    ) throws Exception {
        QueryResponse<TemplateGroupDTO> queryResponse = requestTemplateGroupService.selectAllTemplateGroupService(current, limit, req);
        return JsonResponse.okayWithData(queryResponse);
    }

    @DeleteMapping("/delete/{id}")
    @ApiOperation(value = "delete RequestTemplateGroup", notes = "需要传入id")
    public JsonResponse deleteTemplateGroupByID(@PathVariable("id") String id) throws Exception {
        requestTemplateGroupService.deleteTemplateGroupByIDService(id);
        return JsonResponse.okay();
    }
}

