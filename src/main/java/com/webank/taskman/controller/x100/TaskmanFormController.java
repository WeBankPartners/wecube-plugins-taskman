package com.webank.taskman.controller.x100;


import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.converter.FormItemTemplateConverter;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.base.JsonResponse;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.dto.req.SaveFormTemplateReq;
import com.webank.taskman.dto.resp.FormTemplateResp;
import com.webank.taskman.service.FormItemTemplateService;
import com.webank.taskman.service.FormTemplateService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;


@RestController
@RequestMapping("/v1/form")
@Api(tags = {"5、 Form inteface API"})
public class TaskmanFormController {



    private static final Logger log = LoggerFactory.getLogger(TaskmanFormController.class);

    @Autowired
    FormItemTemplateService formItemTemplateService;


    @Autowired
    FormTemplateService formTemplateService;

    @Autowired
    FormItemTemplateConverter formItemTemplateConverter;

    @ApiOperationSupport(order = 1)
    @PostMapping("/template/save")
    @ApiOperation(value = "form-template-save", notes = "")
    public JsonResponse saveFormTemplate(@Valid @RequestBody SaveFormTemplateReq req) throws Exception {

        FormTemplateResp formTemplateResp= formTemplateService.saveFormTemplateByReq(req);
        return JsonResponse.okayWithData(formTemplateResp);
    }

    @ApiOperationSupport(order = 2)
    @DeleteMapping("/template/delete/{id}")
    @ApiOperation(value = "form-template-delete")
    public JsonResponse deleteFormTemplateByID(@PathVariable("id") String id) throws Exception {
        formTemplateService.deleteFormTemplate(id);
        return JsonResponse.okay();
    }

    @ApiOperationSupport(order = 3)
    @GetMapping("/template/detail/{tempType}/{tempId}")
    @ApiOperation(value = "form-template-detail")
    public JsonResponse<FormTemplateResp> FormTemplateDetail(@PathVariable("tempType") Integer tempType,@PathVariable("tempId") String tempId) throws Exception {
        return JsonResponse.okayWithData(
                formTemplateService.detailFormTemplate(
                        new SaveFormTemplateReq(tempId,tempType)));
    }

    @ApiOperationSupport(order = 4)
    @DeleteMapping("/item/delete/{id}")
    @ApiOperation(value = "form-item-template-delete", notes = "")
    public JsonResponse deleteRequestTemplate(@PathVariable("id") String id) throws Exception {
        formItemTemplateService.deleteRequestTemplateByID(id);
        return JsonResponse.okay();
    }


    @ApiOperationSupport(order = 5)
    @PostMapping("/item/template/currency")
    @ApiOperation(value = "form-item-template-currency")
    public JsonResponse<QueryResponse<FormItemTemplate>> formItemTemplateAvailable()
    {
        QueryWrapper<FormItemTemplate> wrapper = new QueryWrapper<FormItemTemplate>();
        wrapper.eq("status",1);
        return JsonResponse.okayWithData(formItemTemplateService.list(wrapper));
    }


}
