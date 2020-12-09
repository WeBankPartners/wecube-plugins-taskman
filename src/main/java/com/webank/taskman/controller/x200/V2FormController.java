package com.webank.taskman.controller.x200;

import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.req.SaveFormTemplateReq;
import com.webank.taskman.dto.resp.FormTemplateResp;
import com.webank.taskman.service.FormTemplateService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.ApiParam;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.validation.BindingResult;
import org.springframework.validation.ObjectError;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;

@Api(tags = {"4、 Taskman Form model"})
@RestController
@RequestMapping("/v2/form")
public class V2FormController {


    @Autowired
    FormTemplateService formTemplateService;

    //TODO implemented   insert or update
    @PostMapping("/template/save")
    @ApiOperation(value = "save FormTemplate", notes = "Need to pass in object: ")
    public JsonResponse saveRequestTemplate(@Valid @RequestBody SaveFormTemplateReq req, BindingResult bindingResult) throws Exception {

        if (bindingResult.hasErrors()){
            for (ObjectError error:bindingResult.getAllErrors()){
                return JsonResponse.okayWithData(error.getDefaultMessage());
            }
        }
        formTemplateService.saveFormTemplateByReq(req);
        return JsonResponse.okay();
    }

    //TODO Not implemented
    @GetMapping("/template/detail/{tempType}/{tempId}")
    @ApiOperation(value = "//")
    public JsonResponse<FormTemplateResp> formTemplateDetail(
             @ApiParam(value = "模板类型:(1.请求 2.任务)",name = "tempType") @PathVariable("tempType") Integer tempType,
             @ApiParam(value = "模板id",name = "tempId") @PathVariable("tempId") String tempId) throws Exception {

        FormTemplateResp formTemplateResp=formTemplateService.queryDetailByTemp(tempType,tempId);
        return JsonResponse.okayWithData(formTemplateResp);
    }
}
