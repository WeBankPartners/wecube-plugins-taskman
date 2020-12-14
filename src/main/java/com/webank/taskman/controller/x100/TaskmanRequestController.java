package com.webank.taskman.controller.x100;


import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.commons.ApplicationConstants;
import com.webank.taskman.converter.RequestTemplateConverter;
import com.webank.taskman.converter.RequestTemplateGroupConverter;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.*;
import com.webank.taskman.dto.req.QueryRequestTemplateReq;
import com.webank.taskman.dto.req.SaveRequestInfoReq;
import com.webank.taskman.dto.req.SaveRequestTemplateGropReq;
import com.webank.taskman.dto.req.SaveRequestTemplateReq;
import com.webank.taskman.dto.resp.RequestInfoResq;
import com.webank.taskman.dto.resp.RequestTemplateGroupResq;
import com.webank.taskman.dto.resp.RequestTemplateResp;
import com.webank.taskman.service.RequestInfoService;
import com.webank.taskman.service.RequestTemplateGroupService;
import com.webank.taskman.service.RequestTemplateService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.ApiParam;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.validation.BindingResult;
import org.springframework.validation.ObjectError;
import org.springframework.web.bind.annotation.*;
import springfox.documentation.annotations.ApiIgnore;

import javax.servlet.http.HttpServletRequest;
import javax.validation.Valid;
import java.util.List;

import static com.webank.taskman.dto.JsonResponse.okay;
import static com.webank.taskman.dto.JsonResponse.okayWithData;


@Api(tags = {"3、 Request inteface API"})
@RestController
@RequestMapping("/v1/request")
public class TaskmanRequestController {

    @Autowired
    RequestTemplateService requestTemplateService;

    @Autowired
    RequestTemplateGroupService requestTemplateGroupService;


    @Autowired
    RequestTemplateGroupConverter requestTemplateGroupConverter;

    @Autowired
    RequestTemplateConverter requestTemplateConverter;


    @Autowired
    RequestInfoService requestInfoService;

    @ApiOperationSupport(order = 10)
    @PostMapping("/template/save")
    @ApiOperation(value = "request-template-save", notes = "Need to pass in object: ")
    public JsonResponse saveRequestTemplate(@Valid @RequestBody SaveRequestTemplateReq req, BindingResult bindingResult) throws Exception {
        if (bindingResult.hasErrors()){
            for (ObjectError error:bindingResult.getAllErrors()){
                return JsonResponse.okayWithData(error.getDefaultMessage());
            }
        }
      RequestTemplateResp requestTemplateResp= requestTemplateService.saveRequestTemplate(req);
        return JsonResponse.okayWithData(requestTemplateResp);
    }

    @ApiOperationSupport(order = 11)
    @PostMapping("/template/search/{page}/{pageSize}")
    @ApiOperation(value = "request-template-search")
    public JsonResponse<QueryResponse<RequestTemplateResp>> selectRequestTemplate(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize")  @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) QueryRequestTemplateReq req)
            throws Exception {
        QueryResponse<RequestTemplateResp> queryResponse = requestTemplateService.selectAllequestTemplateService(page, pageSize, req);
        return JsonResponse.okayWithData(queryResponse);
    }

    @ApiOperationSupport(order = 12)
    @DeleteMapping("/template/delete/{id}")
    @ApiOperation(value = "request-template-delete", notes = "需要传入id")
    public JsonResponse deleteRequestTemplateByID(@PathVariable("id") String id) throws Exception {
        requestTemplateService.deleteRequestTemplateService(id);
        return okay();
    }

    @ApiOperationSupport(order = 13)
    @GetMapping("/template/detail/{id}")
    @ApiOperation(value = "request-template-detail", notes = "需要传入id")
    public JsonResponse<RequestTemplateResp> detail(@PathVariable("id") String id) throws Exception {
       RequestTemplateResp requestTemplateResp= requestTemplateService.detailRequestTemplate(id);
        return JsonResponse.okayWithData(requestTemplateResp);
    }

    @GetMapping("/template/available")
    @ApiOperationSupport(order = 14)
    @ApiOperation(value = "request-template-available")
    public JsonResponse<List<RequestTemplateResp>> availableRequest(@ApiIgnore @RequestBody(required = false) TemplateGroupReq req) throws Exception {
        QueryWrapper<RequestTemplate> wrapper = new QueryWrapper<RequestTemplate>();
        wrapper.eq("status",1);
        List<RequestTemplateResp> dtoList = requestTemplateConverter.toDto(requestTemplateService.list(wrapper));
        return okayWithData(dtoList);
    }

    @PostMapping("/template/group/save")
    @ApiOperationSupport(order = 14)
    @ApiOperation(value = "request-group-template-save", notes = "")
    public JsonResponse<RequestTemplateGroupResq> createTemplateGroup(
            @Valid @RequestBody SaveRequestTemplateGropReq req, BindingResult bindingResult) throws Exception {
        if (bindingResult.hasErrors()) {
            for (ObjectError error : bindingResult.getAllErrors()) {
                return JsonResponse.okayWithData(error.getDefaultMessage());
            }
        }
        RequestTemplateGroup requestTemplateGroup = requestTemplateGroupService.saveTemplateGroupByReq(req);
        RequestTemplateGroupResq groupResq =new RequestTemplateGroupResq();
        groupResq.setId(requestTemplateGroup.getId());
        return JsonResponse.okayWithData(groupResq);
    }


    @PostMapping("/template/group/search/{page}/{pageSize}")
    @ApiOperationSupport(order = 15)
    @ApiOperation(value = "request-group-template-search")
    public JsonResponse<QueryResponse<TemplateGroupDTO>> selectTemplateGroup(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize")  @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) TemplateGroupReq req
    ) throws Exception {
        QueryResponse<TemplateGroupDTO> queryResponse = requestTemplateGroupService.selectAllTemplateGroupService(page, pageSize, req);
        return JsonResponse.okayWithData(queryResponse);
    }
    @GetMapping("/template/group/available")
    @ApiOperationSupport(order = 16)
    @ApiOperation(value = "request-group-template-available")
    public JsonResponse<List<TemplateGroupDTO>> available(@ApiIgnore @RequestBody(required = false) TemplateGroupReq req) throws Exception {
        QueryWrapper<RequestTemplateGroup> wrapper = new QueryWrapper<RequestTemplateGroup>();
        wrapper.eq("status",0);
        List<TemplateGroupDTO> dtoList = requestTemplateGroupConverter.toDto(requestTemplateGroupService.list(wrapper));
        return JsonResponse.okayWithData(dtoList);
    }

    @DeleteMapping("/template/group/delete/{id}")
    @ApiOperationSupport(order = 17)
    @ApiOperation(value = "request-group-template-delete", notes = "需要传入id")
    public JsonResponse deleteTemplateGroupByID(@PathVariable("id") String id) throws Exception {
        requestTemplateGroupService.deleteTemplateGroupByIDService(id);
        return okay();
    }

    @ApiOperationSupport(order = 18)
    @PostMapping("/save")

    @ApiOperation(value = "Request-Info-save")
    public JsonResponse<SaveRequestInfoReq> saveRequestInfo(@RequestBody SaveRequestInfoReq req) throws Exception {
        return okayWithData(requestInfoService.saveRequestInfo(req));
    }
    
    @ApiOperationSupport(order = 19)
    @PostMapping("/search/{page}/{pageSize}")
    @ApiOperation(value = "Request-Info-search")
    public JsonResponse<QueryResponse<RequestInfoResq>> selectRequestInfo(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize")  @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) SaveRequestInfoReq req)
            throws Exception {
        QueryResponse<RequestInfoResq> queryResponse = requestInfoService.selectRequestInfoService(page, pageSize, req);
        return JsonResponse.okayWithData(queryResponse);
    }





    @ApiOperation(value = "Request-Info-done")
    @PostMapping(ApplicationConstants.ApiInfo.API_RESOURCE_SERVICE_REQUEST_DONE)
    public JsonResponse updateServiceRequest(@RequestBody DoneServiceRequestRequest request,
                                             @ApiIgnore HttpServletRequest httpRequest) throws Exception {
        requestInfoService.doneServiceRequest(request);
        return okay();
    }

}


