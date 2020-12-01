package com.webank.taskman.controller.x100;


import com.webank.taskman.dto.*;
import com.webank.taskman.dto.req.AddRequestTemplateReq;
import com.webank.taskman.service.RequestTemplateService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.DynamicParameter;
import io.swagger.annotations.DynamicParameters;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;


@RestController
@RequestMapping("/v1/request-template")
@Api(tags = {"V1.0.0 RequestTemplate model"})
public class RequestTemplateController {

    @Autowired
    RequestTemplateService requestTemplateService;


    @PostMapping("/add")
    @ApiOperation(value = "add RequestTemplate", notes = "Need to pass in object: ")
    @DynamicParameters(name = "req", properties = {
            @DynamicParameter(name = "page", value = "页码", example = "", required = true, dataTypeClass = Integer.class),
            @DynamicParameter(name = "pageSize", value = "每页行数", example = "100", required = true, dataTypeClass = Integer.class),
            @DynamicParameter(name = "id", value = "主键", example = "", dataTypeClass = String.class),
            @DynamicParameter(name = "name", value = "模板组名称"),
    })
    public JsonResponse createRequestTemplate(@RequestBody RequestTemplateVO requestTemplateVO) throws Exception {
        requestTemplateService.createRequestTemplateService(requestTemplateVO);
        return JsonResponse.okay();
    }

    @PostMapping("edit")
    @ApiOperation(value = "edit RequestTemplate", notes = "")
    public JsonResponse updateRequestTemplate(
            @RequestBody RequestTemplateVO requestTemplateVO) throws Exception {
        requestTemplateService.updateRequestTemplateService(requestTemplateVO);
        return JsonResponse.okay();
    }

    @PostMapping("/search/{current}/{limit}")
    @ApiOperation(value = "search RequestTemplate ")
    @DynamicParameters(name = "req", properties = {
            @DynamicParameter(name = "page", value = "页码", example = "", required = true, dataTypeClass = Integer.class),
            @DynamicParameter(name = "pageSize", value = "每页行数", example = "100", required = true, dataTypeClass = Integer.class),
            @DynamicParameter(name = "id", value = "主键", example = "", dataTypeClass = String.class),
            @DynamicParameter(name = "name", value = "模板名称"),
    })
    public JsonResponse<QueryResponse<RequestTemplateDTO>> selectRequestTemplate(
            @PathVariable("current") Integer current,
            @PathVariable("limit") Integer limit,
            @RequestBody(required = false) RequestTemplateReq req)
            throws Exception {
        QueryResponse<RequestTemplateDTO> queryResponse = requestTemplateService.selectAllequestTemplateService(current, limit, req);
        return JsonResponse.okayWithData(queryResponse);
    }

    @DeleteMapping("/delete/{id}")
    @ApiOperation(value = "delete RequestTemplate", notes = "需要传入id")
    public JsonResponse deleteRequestTemplateByID(@PathVariable("id") String id) throws Exception {
        requestTemplateService.deleteRequestTemplateService(id);
        return JsonResponse.okay();
    }
}

