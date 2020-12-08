package com.webank.taskman.controller.x100;


import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.QueryRequestTemplateReq;
import com.webank.taskman.dto.req.SaveRequestTemplateReq;
import com.webank.taskman.dto.resp.RequestTemplateResp;
import com.webank.taskman.service.RequestTemplateService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.ApiParam;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.validation.BindingResult;
import org.springframework.validation.ObjectError;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;


@Api(tags = {"4、 RequestTemplate inteface API"})
@RestController
@RequestMapping("/v1/request/template")
public class RequestTemplateController {

    @Autowired
    RequestTemplateService requestTemplateService;

    //TODO implemented   insert or update
    @ApiOperationSupport(order = 1)
    @PostMapping("/save")
    @ApiOperation(value = "save RequestTemplate", notes = "Need to pass in object: ")
    public JsonResponse saveRequestTemplate(@Valid @RequestBody SaveRequestTemplateReq req, BindingResult bindingResult) throws Exception {
        if (bindingResult.hasErrors()){
            for (ObjectError error:bindingResult.getAllErrors()){
                return JsonResponse.okayWithData(error.getDefaultMessage());
            }
        }
      RequestTemplateResp requestTemplateResp= requestTemplateService.saveRequestTemplate(req);
        return JsonResponse.okayWithData(requestTemplateResp);
    }

    //TODO Not implemented
    @ApiOperationSupport(order = 2)
    @PostMapping("/search/{page}/{pageSize}")
    @ApiOperation(value = "search RequestTemplate ")
    public JsonResponse<QueryResponse<RequestTemplateResp>> selectRequestTemplate(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize")  @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) QueryRequestTemplateReq req)
            throws Exception {
        QueryResponse<RequestTemplateResp> queryResponse = requestTemplateService.selectAllequestTemplateService(page, pageSize, req);
        return JsonResponse.okayWithData(queryResponse);
    }

    //TODO Not implemented
    @ApiOperationSupport(order = 3)
    @DeleteMapping("/delete/{id}")
    @ApiOperation(value = "delete RequestTemplate", notes = "需要传入id")
    public JsonResponse deleteRequestTemplateByID(@PathVariable("id") String id) throws Exception {
        requestTemplateService.deleteRequestTemplateService(id);
        return JsonResponse.okay();
    }

    //TODO Not implemented
    @ApiOperationSupport(order = 4)
    @GetMapping("/detail/{id}")
    @ApiOperation(value = "detail RequestTemplate ", notes = "需要传入id")
    public JsonResponse<RequestTemplateResp> detail(@PathVariable("id") String id) throws Exception {
       RequestTemplateResp requestTemplateResp= requestTemplateService.detailRequestTemplate(id);
        return JsonResponse.okayWithData(requestTemplateResp);
    }
}


