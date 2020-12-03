package com.webank.taskman.controller.x100;


import com.github.xiaoymin.knife4j.annotations.DynamicParameter;
import com.github.xiaoymin.knife4j.annotations.DynamicParameters;
import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.RequestTemplateReq;
import com.webank.taskman.dto.req.SaveFormItemTemplateReq;
import com.webank.taskman.dto.resp.RequestTemplateResp;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/v1/request/form")
@Api(tags = {"3、 TemplateGroup model"})
public class FormTemplateController {


    //TODO implemented   insert or update
    @PostMapping("/save")
    @ApiOperation(value = "save formTemplate", notes = "Need to pass in object: ")
    public JsonResponse saveRequestTemplate(@RequestBody SaveFormItemTemplateReq req) throws Exception {


        return JsonResponse.okay();
    }

    //TODO Not implemented
    @PostMapping("/search/{current}/{limit}")
    @ApiOperation(value = "search RequestTemplate ")
    @DynamicParameters(name = "req", properties = {
            @DynamicParameter(name = "page", value = "页码", example = "", required = true, dataTypeClass = Integer.class),
            @DynamicParameter(name = "pageSize", value = "每页行数", example = "100", required = true, dataTypeClass = Integer.class),
    })
    public JsonResponse<QueryResponse<RequestTemplateResp>> selectRequestTemplate(
            @PathVariable("current") Integer current,
            @PathVariable("limit") Integer limit,
            @RequestBody(required = false) RequestTemplateReq req)
            throws Exception {
        return JsonResponse.okayWithData(null);
    }

    //TODO Not implemented
    @DeleteMapping("/delete/{id}")
    @ApiOperation(value = "delete RequestTemplate", notes = "需要传入id")
    public JsonResponse deleteRequestTemplateByID(@PathVariable("id") String id) throws Exception {
        return JsonResponse.okay();
    }

    //TODO Not implemented
    @GetMapping("/detail/{id}")
    @ApiOperation(value = " RequestTemplate ", notes = "需要传入id")
    public JsonResponse<RequestTemplateResp> detail(@PathVariable("id") String id) throws Exception {
        return JsonResponse.okay();
    }
}
