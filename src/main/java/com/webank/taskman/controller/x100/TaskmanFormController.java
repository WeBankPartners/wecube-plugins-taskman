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
@Api(tags = {"3、 FormModel"})
public class TaskmanFormController {

    @Autowired
    FormItemTemplateService formItemTemplateService;


    @Autowired
    FormTemplateService formTemplateService;

    @ApiOperationSupport(order = 1)
    @PostMapping("/save")
    @ApiOperation(value = "save FormTemplate", notes = "Need to pass in object: ")
    public JsonResponse saveRequestTemplate(@Valid @RequestBody SaveFormTemplateReq req, BindingResult bindingResult) throws Exception {

        if (bindingResult.hasErrors()){
            for (ObjectError error:bindingResult.getAllErrors()){
                return JsonResponse.okayWithData(error.getDefaultMessage());
            }
        }
        FormTemplateResp formTemplateResp= formTemplateService.saveFormTemplateByReq(req);
        return JsonResponse.okayWithData(formTemplateResp);
    }

    @ApiOperationSupport(order = 3)
    @DeleteMapping("/delete/{id}")
    @ApiOperation(value = "delete RequestTemplate")
    public JsonResponse deleteFormTemplateByID(@PathVariable("id") String id) throws Exception {
        formTemplateService.deleteFormTemplate(id);
        return JsonResponse.okay();
    }

    @ApiOperationSupport(order = 4)
    @GetMapping("/detail/{tempType}/{tempId}")
    @ApiOperation(value = "detail FormTemplate ")
    public JsonResponse<QueryResponse<FormTemplateResp>> detail(@PathVariable("tempType") Integer tempType,@PathVariable("tempId") String tempId) throws Exception {
        SaveFormTemplateReq req = new SaveFormTemplateReq(tempId,tempType);
        FormTemplateResp formTemplateResp = formTemplateService.detailFormTemplate(req);
        return JsonResponse.okayWithData(formTemplateResp);
    }

    @PostMapping("/item/save")
    @ApiOperation(value = "save or update FormItemTemplate", notes = "Need to pass in object: ")
    public JsonResponse createFormItemTemplate(@Valid @RequestBody SaveFormItemTemplateReq req, BindingResult bindingResult) throws Exception {
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

    //TODO Not implemented
    @PostMapping("/item/search/{page}/{pageSize}")
    @ApiOperation(value = "search RequestTemplate ")
    public JsonResponse<QueryResponse<FormItemTemplateResq>> selectFormItemTemplate(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize")  @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) SelectFormItemTemplateReq req)
            throws Exception {
        QueryResponse<FormItemTemplateResq> queryResponse= formItemTemplateService.selectAllFormItemTemplateService(page, pageSize, req);
        return JsonResponse.okayWithData(queryResponse);
    }


    //TODO Not implemented
    @DeleteMapping("/item/delete/{id}")
    @ApiOperation(value = "delete RequestTemplate", notes = "需要传入id")
    public JsonResponse deleteRequestTemplateByID(@PathVariable("id") String id) throws Exception {
        formItemTemplateService.deleteRequestTemplateByID(id);
        return JsonResponse.okay();
    }

}
