package com.webank.taskman.controller.x100;


import com.github.xiaoymin.knife4j.annotations.DynamicParameter;
import com.github.xiaoymin.knife4j.annotations.DynamicParameters;
import com.webank.taskman.dto.*;
import com.webank.taskman.dto.req.SaveRequestTemplateReq;
import com.webank.taskman.dto.req.SaveTaskTemplateReq;
import com.webank.taskman.dto.resp.RequestTemplateResp;
import com.webank.taskman.service.RequestTemplateService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;


@RestController
@RequestMapping("/v1/task/template")
@Api(tags = {"5、  TaskTemplate inteface API"})
public class TaskTemplateController {

    @Autowired
    RequestTemplateService requestTemplateService;


    //TODO implemented   insert or update
    @PostMapping("/save")
    @ApiOperation(value = "add TaskTemplate", notes = "Need to pass in object: ")
    public JsonResponse createRequestTemplate(@RequestBody SaveTaskTemplateReq taskTemplateReq) throws Exception {
        return JsonResponse.okay();
    }
    @DeleteMapping("/delete/{id}")
    @ApiOperation(value = "delete TaskTemplate", notes = "需要传入id")
    public JsonResponse deleteRequestTemplateByID(@PathVariable("id") String id) throws Exception {
        requestTemplateService.deleteRequestTemplateService(id);
        return JsonResponse.okay();
    }

    @PostMapping("/search/{current}/{limit}")
    @ApiOperation(value = "search TaskTemplate ")
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

    @GetMapping("/detail/{id}")
    @ApiOperation(value = "detail TaskTemplate", notes = "需要传入id")
    public JsonResponse detail(@PathVariable("id") String id) throws Exception {
        requestTemplateService.deleteRequestTemplateService(id);
        return JsonResponse.okay();
    }
}

