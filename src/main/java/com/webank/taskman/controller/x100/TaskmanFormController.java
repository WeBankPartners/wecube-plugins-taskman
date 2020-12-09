package com.webank.taskman.controller.x100;


import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.SaveFormItemTemplateReq;
import com.webank.taskman.dto.req.SaveFormTemplateReq;
import com.webank.taskman.dto.req.SelectFormItemTemplateReq;
import com.webank.taskman.dto.resp.FormItemTemplateResq;
import com.webank.taskman.dto.resp.FormItemTemplateSVResq;
import com.webank.taskman.dto.resp.FormTemplateResp;
import com.webank.taskman.service.FormItemTemplateService;
import com.webank.taskman.service.FormTemplateService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.ApiParam;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.validation.BindingResult;
import org.springframework.validation.ObjectError;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;


@RestController
@RequestMapping("/v1/form")
@Api(tags = {"5、 Form inteface API"})
public class TaskmanFormController {

    @Autowired
    FormItemTemplateService formItemTemplateService;


    @Autowired
    FormTemplateService formTemplateService;

    @ApiOperationSupport(order = 1)
    @PostMapping("/template/save")
    @ApiOperation(value = "form-template-save", notes = "")
    public JsonResponse saveFormTemplate(@Valid @RequestBody SaveFormTemplateReq req, BindingResult bindingResult) throws Exception {

        if (bindingResult.hasErrors()){
            for (ObjectError error:bindingResult.getAllErrors()){
                return JsonResponse.okayWithData(error.getDefaultMessage());
            }
        }
        FormTemplateResp formTemplateResp= formTemplateService.saveFormTemplateByReq(req);
        return JsonResponse.okayWithData(formTemplateResp);
    }

    @ApiOperationSupport(order = 3)
    @DeleteMapping("/template/delete/{id}")
    @ApiOperation(value = "form-template-delete")
    public JsonResponse deleteFormTemplateByID(@PathVariable("id") String id) throws Exception {
        formTemplateService.deleteFormTemplate(id);
        return JsonResponse.okay();
    }

    @ApiOperationSupport(order = 4)
    @GetMapping("/template/detail/{tempType}/{tempId}")
    @ApiOperation(value = "form-template-detail")
    public JsonResponse<QueryResponse<FormTemplateResp>> FormTemplateDetail(@PathVariable("tempType") Integer tempType,@PathVariable("tempId") String tempId) throws Exception {
        SaveFormTemplateReq req = new SaveFormTemplateReq(tempId,tempType);
        FormTemplateResp formTemplateResp = formTemplateService.detailFormTemplate(req);
        return JsonResponse.okayWithData(formTemplateResp);
    }

    @ApiOperationSupport(order = 5)
    @PostMapping("/item/template/save")
    @ApiOperation(value = "form-item-template-save", notes = "")
    public JsonResponse saveFormItemTemplate(@Valid @RequestBody SaveFormItemTemplateReq req, BindingResult bindingResult) throws Exception {
        if (bindingResult.hasErrors()) {
            for (ObjectError error : bindingResult.getAllErrors()) {
                return JsonResponse.okayWithData(error.getDefaultMessage());
            }
        }
        FormItemTemplate formItemTemplate = formItemTemplateService.saveFormItemTemplateByReq(req);
        FormItemTemplateSVResq formItemTemplateSVResq=new FormItemTemplateSVResq();
        formItemTemplateSVResq.setId(formItemTemplate.getId());
        return JsonResponse.okayWithData(formItemTemplateSVResq);
    }

    @ApiOperationSupport(order = 6)
    @DeleteMapping("/item/delete/{id}")
    @ApiOperation(value = "form-item-template-delete", notes = "")
    public JsonResponse deleteRequestTemplate(@PathVariable("id") String id) throws Exception {
        formItemTemplateService.deleteRequestTemplateByID(id);
        return JsonResponse.okay();
    }

    @ApiOperationSupport(order = 7)
    @PostMapping("/item/template/search/{page}/{pageSize}")
    @ApiOperation(value = "form-item-template-serach ")
    public JsonResponse<QueryResponse<FormItemTemplateResq>> searchFormItemTemplate(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize")  @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) SelectFormItemTemplateReq req)
    {
        QueryResponse<FormItemTemplateResq> queryResponse= formItemTemplateService.selectAllFormItemTemplateService(page, pageSize, req);
        return JsonResponse.okayWithData(queryResponse);
    }



}
