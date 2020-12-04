package com.webank.taskman.controller.x100;


import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.github.xiaoymin.knife4j.annotations.DynamicParameter;
import com.github.xiaoymin.knife4j.annotations.DynamicParameters;
import com.webank.taskman.converter.RequestTemplateGroupConverter;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.domain.RoleInfo;
import com.webank.taskman.dto.*;
import com.webank.taskman.dto.req.SaveAndUpdateTemplateGropReq;
import com.webank.taskman.dto.req.SaveTemplateGropReq;
import com.webank.taskman.dto.resp.RequestTemplateGroupResq;
import com.webank.taskman.service.RequestTemplateGroupService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.validation.BindingResult;
import org.springframework.validation.ObjectError;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;
import java.util.List;

@RestController
@RequestMapping("/v1/request/template/group")
@Api(tags = {"3、 TemplateGroup model"})
public class RequestTemplateGroupController {
    @Autowired
    RequestTemplateGroupService requestTemplateGroupService;

    @Autowired
    RequestTemplateGroupConverter requestTemplateGroupConverter;

    //FIXME implemented   insert or update
    @PostMapping("/save")
    @ApiOperationSupport(order = 21)
    @ApiOperation(value = "save 0r update RequestTemplateGroup", notes = "")
    public JsonResponse createTemplateGroup(
            @Valid @RequestBody SaveAndUpdateTemplateGropReq req, BindingResult bindingResult) throws Exception {
        if (bindingResult.hasErrors()) {
            for (ObjectError error : bindingResult.getAllErrors()) {
                return JsonResponse.okayWithData(error.getDefaultMessage());
            }
        }
        RequestTemplateGroup requestTemplateGroup = requestTemplateGroupService.addOrUpdateTemplateGroup(req);
        RequestTemplateGroupResq groupResq =new RequestTemplateGroupResq();
        groupResq.setId(requestTemplateGroup.getId());
        return JsonResponse.okayWithData(groupResq);
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

