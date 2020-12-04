package com.webank.taskman.controller.x100;


import com.github.xiaoymin.knife4j.annotations.DynamicParameter;
import com.github.xiaoymin.knife4j.annotations.DynamicParameters;
import com.webank.taskman.converter.FormItemTemplateConverter;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.QueryResponse;

import com.webank.taskman.dto.TemplateGroupDTO;
import com.webank.taskman.dto.req.*;
import com.webank.taskman.dto.resp.FormItemTemplateResq;
import com.webank.taskman.dto.resp.FormItemTemplateSVResq;
import com.webank.taskman.dto.resp.RequestTemplateResp;
import com.webank.taskman.service.FormItemTemplateService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiModelProperty;
import io.swagger.annotations.ApiOperation;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.validation.BindingResult;
import org.springframework.validation.ObjectError;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;


@RestController
@RequestMapping("/v1/request/form/item")
@Api(tags = {"3、 FormItemTemplate model"})
public class FormItemTemplateController {

    @Autowired
    FormItemTemplateService formItemTemplateService;


    //TODO implemented   insert or update
    @PostMapping("/save")
    @ApiOperation(value = "save or update FormItemTemplate", notes = "Need to pass in object: ")
    public JsonResponse createFormItemTemplate(@Valid @RequestBody SaveAndUpdateFormItemTemplateReq req, BindingResult bindingResult) throws Exception {
        if (bindingResult.hasErrors()) {
            for (ObjectError error : bindingResult.getAllErrors()) {
                return JsonResponse.okayWithData(error.getDefaultMessage());
            }
        }
        FormItemTemplate formItemTemplate = formItemTemplateService.addOrUpdateFormItemTemplate(req);
        FormItemTemplateSVResq formItemTemplateSVResq=new FormItemTemplateSVResq();
        formItemTemplateSVResq.setId(formItemTemplate.getId());
        return JsonResponse.okayWithData(formItemTemplateSVResq);
    }

    //TODO Not implemented
    @PostMapping("/search/{current}/{limit}")
    @ApiOperation(value = "search RequestTemplate ")
    @DynamicParameters(name = "req", properties = {
            @DynamicParameter(name = "page", value = "页码", example = "", required = true, dataTypeClass = Integer.class),
            @DynamicParameter(name = "pageSize", value = "每页行数", example = "100", required = true, dataTypeClass = Integer.class),
    })
    public JsonResponse<QueryResponse<FormItemTemplateResq>> selectFormItemTemplate(
            @PathVariable("current") Integer current,
            @PathVariable("limit") Integer limit,
            @RequestBody(required = false) SelectFormItemTemplateReq req)
            throws Exception {
        QueryResponse<FormItemTemplateResq> queryResponse= formItemTemplateService.selectAllFormItemTemplateService(current, limit, req);
        return JsonResponse.okayWithData(queryResponse);
    }


    //TODO Not implemented
    @DeleteMapping("/delete/{id}")
    @ApiOperation(value = "delete RequestTemplate", notes = "需要传入id")
    public JsonResponse deleteRequestTemplateByID(@PathVariable("id") String id) throws Exception {
        formItemTemplateService.deleteRequestTemplateByID(id);
        return JsonResponse.okay();
    }

    //TODO Not implemented
    @GetMapping("/detail/{tempType}/{tempId}")
    @ApiOperation(value = " RequestTemplate ", notes = "需要传入id")
    public JsonResponse<RequestTemplateResp> detail(@PathVariable("id") String id) throws Exception {
        return JsonResponse.okay();
    }
}
