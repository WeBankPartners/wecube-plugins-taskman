package com.webank.taskman.controller.x100;


import com.webank.taskman.converter.RequestTemplateGroupConverter;
import com.webank.taskman.domain.RoleInfo;
import com.webank.taskman.dto.*;
import com.webank.taskman.dto.req.AddTemplateGropReq;
import com.webank.taskman.service.RequestTemplateGroupService;
import io.swagger.annotations.*;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/v1/request-template-group")
@Api(tags = {"2、 TemplateGroup model"})
public class RequestTemplateGroupController {
    @Autowired
    RequestTemplateGroupService requestTemplateGroupService;

    @Autowired
    RequestTemplateGroupConverter requestTemplateGroupConverter;

    @PostMapping("/add")
    @ApiOperationSupport(order = 21)
    @ApiOperation(value = "add RequestTemplateGroup", notes = "")
    public JsonResponse createTemplateGroup(
            @RequestBody AddTemplateGropReq req) throws Exception {
        if(StringUtils.isEmpty(req.getName())){
            return  JsonResponse.error(" manageRoleId is null");
        }
        if(StringUtils.isEmpty(req.getName())){
            return  JsonResponse.error(" name is null");
        }
        requestTemplateGroupService.addTemplateGroup(requestTemplateGroupConverter.addReqDomain(req));
        return JsonResponse.okay();
    }

    @PostMapping("edit")
    @ApiOperationSupport(order = 23)
    @ApiOperation(value = "edit RequestTemplateGroup", notes = "")
    public JsonResponse updateTemplateGroup(
            @RequestBody TemplateGroupVO templateGroupVO) throws Exception {
        requestTemplateGroupService.updateTemplateGroupService(templateGroupVO);
        return JsonResponse.okay();
    }

    @PostMapping("/search/{current}/{limit}")
    @ApiOperationSupport(order = 24)
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
        queryResponse.getContents().forEach(c->{
            c.setManageRole(new RoleInfo("2c9280827019695c017019ac974f001c","SUPER_ADMIN"));


        });
        return JsonResponse.okayWithData(queryResponse);
    }

    @PostMapping("/available")
    @ApiOperationSupport(order = 25)
    @ApiOperation(value = "query available RequestTemplateGroup List")
    public JsonResponse<List<TemplateGroupDTO>> available() throws Exception {
        List<TemplateGroupDTO> dtoList = requestTemplateGroupConverter.toDto(requestTemplateGroupService.list());
        return JsonResponse.okayWithData(dtoList);
    }

    @DeleteMapping("/delete/{id}")
    @ApiOperationSupport(order = 22)
    @ApiOperation(value = "delete RequestTemplateGroup", notes = "需要传入id")
    public JsonResponse deleteTemplateGroupByID(@PathVariable("id") String id) throws Exception {
        requestTemplateGroupService.deleteTemplateGroupByIDService(id);
        return JsonResponse.okay();
    }
}

