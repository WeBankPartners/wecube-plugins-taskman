package com.webank.taskman.controller.x100;


import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.converter.RequestTemplateGroupConverter;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.TemplateGroupDTO;
import com.webank.taskman.dto.TemplateGroupReq;
import com.webank.taskman.dto.req.SaveAndUpdateTemplateGropReq;
import com.webank.taskman.dto.resp.RequestTemplateGroupResq;
import com.webank.taskman.service.RequestTemplateGroupService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.ApiParam;
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
    public JsonResponse<RequestTemplateGroupResq> createTemplateGroup(
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


    @PostMapping("/search/{page}/{pageSize}")
    @ApiOperationSupport(order = 24)
    @ApiOperation(value = "search RequestTemplateGroup ")
    public JsonResponse<QueryResponse<TemplateGroupDTO>> selectTemplateGroup(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize")  @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) TemplateGroupReq req
    ) throws Exception {
        QueryResponse<TemplateGroupDTO> queryResponse = requestTemplateGroupService.selectAllTemplateGroupService(page, pageSize, req);
        return JsonResponse.okayWithData(queryResponse);
    }
    @PostMapping("/available")
    @ApiOperationSupport(order = 25)
    @ApiOperation(value = "query available RequestTemplateGroup List")
    public JsonResponse<List<TemplateGroupDTO>> available( @RequestBody(required = false) TemplateGroupReq req) throws Exception {
        QueryWrapper<RequestTemplateGroup> wrapper = new QueryWrapper<>();
        wrapper.eq(!StringUtils.isEmpty(req.getId()),"id", req.getId());
        wrapper.eq(!StringUtils.isEmpty(req.getManageRole()),"manage_role_id", req.getManageRole());
        wrapper.like(!StringUtils.isEmpty(req.getName()),"name", req.getName());
        List<TemplateGroupDTO> dtoList = requestTemplateGroupConverter.toDto(requestTemplateGroupService.list(wrapper));
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

