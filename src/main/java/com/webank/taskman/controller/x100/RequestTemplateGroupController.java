package com.webank.taskman.controller.x100;


import com.webank.taskman.converter.TemplateGroupConverter;
import com.webank.taskman.dto.*;
import com.webank.taskman.dto.req.AddTemplateGropReq;
import com.webank.taskman.service.RequestTemplateGroupService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.DynamicParameter;
import io.swagger.annotations.DynamicParameters;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/v1/request-template-group")
@Api(tags = {"V1.0.0 TemplateGroup model"})
public class RequestTemplateGroupController {
    @Autowired
    RequestTemplateGroupService requestTemplateGroupService;

    @Autowired
    TemplateGroupConverter templateGroupConverter;

    @PostMapping("/add")
    @ApiOperation(value = "add RequestTemplateGroup", notes = "")
    public JsonResponse createTemplateGroup(
            @RequestBody AddTemplateGropReq req) throws Exception {
        if(StringUtils.isEmpty(req.getName())){
            return  JsonResponse.error(" manageRoleId is null");
        }
        if(StringUtils.isEmpty(req.getName())){
            return  JsonResponse.error(" name is null");
        }
        requestTemplateGroupService.addTemplateGroup(templateGroupConverter.addReqDomain(req));
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

