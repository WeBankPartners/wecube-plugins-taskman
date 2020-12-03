package com.webank.taskman.controller.x100;


import com.github.xiaoymin.knife4j.annotations.DynamicParameter;
import com.github.xiaoymin.knife4j.annotations.DynamicParameters;
import com.webank.taskman.dto.*;
import com.webank.taskman.dto.req.SaveRequestTemplateReq;
import com.webank.taskman.dto.resp.RequestTemplateResp;
import com.webank.taskman.service.RequestTemplateService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;


@Api(tags = {"4、 RequestTemplate inteface API"})
@RestController
@RequestMapping("/v1/request/template")
public class RequestTemplateController {

    @Autowired
    RequestTemplateService requestTemplateService;

    //TODO implemented   insert or update
    @PostMapping("/save")
    @ApiOperation(value = "save RequestTemplate", notes = "Need to pass in object: ")
    public JsonResponse saveRequestTemplate(@RequestBody SaveRequestTemplateReq req) throws Exception {
        requestTemplateService.saveRequestTemplate(req);
        return JsonResponse.okay();
    }

    //TODO Not implemented
    @PostMapping("/search/{current}/{limit}")
    @ApiOperation(value = "search RequestTemplate ")
    @DynamicParameters(name = "req", properties = {
            @DynamicParameter(name = "page", value = "页码", example = "", required = true, dataTypeClass = Integer.class),
            @DynamicParameter(name = "pageSize", value = "每页行数", example = "100", required = true, dataTypeClass = Integer.class),
            @DynamicParameter(name = "id", value = "主键", example = "", dataTypeClass = String.class),
            @DynamicParameter(name = "name", value = "模板名称"),
    })
    public JsonResponse<QueryResponse<RequestTemplateResp>> selectRequestTemplate(
            @PathVariable("current") Integer current,
            @PathVariable("limit") Integer limit,
            @RequestBody(required = false) SaveRequestTemplateReq req)
            throws Exception {
        QueryResponse<RequestTemplateResp> queryResponse = requestTemplateService.selectAllequestTemplateService(current, limit, req);
        return JsonResponse.okayWithData(queryResponse);
    }

    //TODO Not implemented
    @DeleteMapping("/delete/{id}")
    @ApiOperation(value = "delete RequestTemplate", notes = "需要传入id")
    public JsonResponse deleteRequestTemplateByID(@PathVariable("id") String id) throws Exception {
        requestTemplateService.deleteRequestTemplateService(id);
        return JsonResponse.okay();
    }

    //TODO Not implemented
    @GetMapping("/detail/{id}")
    @ApiOperation(value = "detail RequestTemplate ", notes = "需要传入id")
    public JsonResponse<RequestTemplateResp> detail(@PathVariable("id") String id) throws Exception {
       RequestTemplateResp requestTemplateResp= requestTemplateService.detailRequestTemplate(id);
        return JsonResponse.okayWithData(requestTemplateResp);
    }
}


