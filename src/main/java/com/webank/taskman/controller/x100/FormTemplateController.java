package com.webank.taskman.controller.x100;


import com.github.xiaoymin.knife4j.annotations.DynamicParameter;
import com.github.xiaoymin.knife4j.annotations.DynamicParameters;
import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.SaveFormTemplateReq;
import com.webank.taskman.dto.resp.FormTemplateResp;
import com.webank.taskman.service.FormTemplateService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/v1/request/form")
@Api(tags = {"3、 FormTemplate model"})
public class FormTemplateController {

    @Autowired
    FormTemplateService formTemplateService;

    //TODO implemented   insert or update
    @PostMapping("/save")
    @ApiOperation(value = "save FormTemplate", notes = "Need to pass in object: ")
    public JsonResponse saveRequestTemplate(@RequestBody SaveFormTemplateReq req) throws Exception {
        formTemplateService.saveFormTemplate(req);
        return JsonResponse.okay();
    }

    //TODO Not implemented
    @PostMapping("/search/{current}/{limit}")
    @ApiOperation(value = "search FormTemplate ")
    @DynamicParameters(name = "req", properties = {
            @DynamicParameter(name = "page", value = "页码", example = "", required = true, dataTypeClass = Integer.class),
            @DynamicParameter(name = "pageSize", value = "每页行数", example = "100", required = true, dataTypeClass = Integer.class),
    })
    public JsonResponse<QueryResponse<FormTemplateResp>> selectFormTemplate(
            @PathVariable("current") Integer current,
            @PathVariable("limit") Integer limit,
            @RequestBody(required = false) SaveFormTemplateReq req)
            throws Exception {
       QueryResponse<FormTemplateResp> queryResponse= formTemplateService.selectFormTemplate(current,limit,req);
        return JsonResponse.okayWithData(queryResponse);
    }

    //TODO Not implemented
    @DeleteMapping("/delete/{id}")
    @ApiOperation(value = "delete RequestTemplate", notes = "需要传入id")
    public JsonResponse deleteFormTemplateByID(@PathVariable("id") String id) throws Exception {
        formTemplateService.deleteFormTemplate(id);
        return JsonResponse.okay();
    }

    //TODO Not implemented
    @GetMapping("/detail/{id}")
    @ApiOperation(value = "detail FormTemplate ", notes = "需要传入id")
    public JsonResponse<QueryResponse<FormTemplateResp>> detail(@PathVariable("id") String id) throws Exception {
        FormTemplateResp formTemplateResp=formTemplateService.detailFormTemplate(id);
        return JsonResponse.okayWithData(formTemplateResp);
    }
}
