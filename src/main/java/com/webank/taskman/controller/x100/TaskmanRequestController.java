package com.webank.taskman.controller.x100;


import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.commons.ApplicationConstants;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanException;
import com.webank.taskman.constant.StatusCodeEnum;
import com.webank.taskman.constant.StatusEnum;
import com.webank.taskman.converter.RequestTemplateConverter;
import com.webank.taskman.converter.RequestTemplateGroupConverter;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.*;
import com.webank.taskman.dto.req.*;
import com.webank.taskman.dto.resp.*;
import com.webank.taskman.service.RequestInfoService;
import com.webank.taskman.service.RequestSynthesisService;
import com.webank.taskman.service.RequestTemplateGroupService;
import com.webank.taskman.service.RequestTemplateService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.ApiParam;
import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
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

    private static final Logger log = LoggerFactory.getLogger(TaskmanRequestController.class);

    @Autowired
    RequestSynthesisService requestSynthesisService;

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
    public JsonResponse saveRequestTemplate(@Valid @RequestBody SaveRequestTemplateReq req
            )throws Exception{
      RequestTemplateResp requestTemplateResp= requestTemplateService.saveRequestTemplate(req);
        return JsonResponse.okayWithData(requestTemplateResp);
    }

    @ApiOperationSupport(order = 11,ignoreParameters = {"requestTempGroup","procDefId","procDefName","description","name","tags","manageRoles","useRoles"})
    @PostMapping("/template/release")
    @ApiOperation(value = "request-template-release", notes = "Need to pass in object: ")
    public JsonResponse releaseRequestTemplate(@RequestBody SaveRequestTemplateReq req) throws Exception {
        if(StringUtils.isEmpty(req.getId())){
            throw new TaskmanException(StatusCodeEnum.PARAM_ISNULL);
        }
        RequestTemplate requestTemplate = requestTemplateService.getOne(new QueryWrapper<RequestTemplate>()
                    .eq("id",req.getId()).eq("status",StatusEnum._DEFAULT.ordinal()));
        if(null == requestTemplate){
            throw new TaskmanException(StatusCodeEnum.PARAM_ISNULL);
        }
        requestTemplate.setStatus(StatusEnum.ENABLE.ordinal());
        requestTemplate.setUpdatedBy(AuthenticationContextHolder.getCurrentUsername());
        requestTemplateService.updateById(requestTemplate);
        return okay();
    }

    @ApiOperationSupport(order = 12)
    @PostMapping("/template/search/{page}/{pageSize}")
    @ApiOperation(value = "request-template-search")
    public JsonResponse<QueryResponse<RequestTemplateResp>> selectRequestTemplatePage(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize")  @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) QueryRequestTemplateReq req)
            throws Exception {
        QueryResponse<RequestTemplateResp> queryResponse = requestTemplateService.selectRequestTemplatePage(page, pageSize, req);
        return JsonResponse.okayWithData(queryResponse);
    }

    @ApiOperationSupport(order = 13)
    @DeleteMapping("/template/delete/{id}")
    @ApiOperation(value = "request-template-delete", notes = "需要传入id")
    public JsonResponse deleteRequestTemplateByID(@PathVariable("id") String id) throws Exception {
        requestTemplateService.deleteRequestTemplateService(id);
        return okay();
    }

    @ApiOperationSupport(order = 14)
    @GetMapping("/template/detail/{id}")
    @ApiOperation(value = "request-template-detail", notes = "需要传入id")
    public JsonResponse<RequestTemplateResp> detail(@PathVariable("id") String id) throws Exception {
       RequestTemplateResp requestTemplateResp= requestTemplateService.detailRequestTemplate(id);
        return JsonResponse.okayWithData(requestTemplateResp);
    }

    @GetMapping(value = {"/template/available","/template/available/{all}"})
    @ApiOperationSupport(order = 15)
    @ApiOperation(value = "request-template-available")
    public JsonResponse<List<RequestTemplateResp>> availableRequest( @PathVariable(value = "all",required = false) String all,@ApiIgnore QueryRequestTemplateReq req) throws Exception {
        if(StringUtils.isEmpty(all)){
            req.setStatus(StatusEnum.ENABLE.ordinal());
        }
        AuthenticationContextHolder.getCurrentUsername();
        req.setSourceTableFix("rt");
        req.setUseRoleName(AuthenticationContextHolder.getCurrentUserRolesToString());
        List<RequestTemplateResp> dtoList = requestTemplateService.selectAvailableRequest(req);
        return okayWithData(dtoList);
    }

    @PostMapping("/template/group/save")
    @ApiOperationSupport(order = 16)
    @ApiOperation(value = "request-group-template-save", notes = "")
    public JsonResponse<RequestTemplateGroupResq> createTemplateGroup(
            @Valid @RequestBody SaveRequestTemplateGropReq req) throws Exception {
        RequestTemplateGroup requestTemplateGroup = requestTemplateGroupService.saveTemplateGroupByReq(req);
        RequestTemplateGroupResq groupResq =new RequestTemplateGroupResq();
        groupResq.setId(requestTemplateGroup.getId());
        return JsonResponse.okayWithData(groupResq);
    }

    @PostMapping("/template/group/search/{page}/{pageSize}")
    @ApiOperationSupport(order = 17)
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
    @ApiOperationSupport(order = 18)
    @ApiOperation(value = "request-group-template-available")
    public JsonResponse<List<TemplateGroupDTO>> available(@ApiIgnore @RequestBody(required = false) TemplateGroupReq req) throws Exception {
        QueryWrapper<RequestTemplateGroup> wrapper = new QueryWrapper<RequestTemplateGroup>();
        wrapper.eq("status",0);
        List<TemplateGroupDTO> dtoList = requestTemplateGroupConverter.toDto(requestTemplateGroupService.list(wrapper));
        return JsonResponse.okayWithData(dtoList);
    }

    @DeleteMapping("/template/group/delete/{id}")
    @ApiOperationSupport(order = 19)
    @ApiOperation(value = "request-group-template-delete", notes = "需要传入id")
    public JsonResponse deleteTemplateGroupByID(@PathVariable("id") String id) throws Exception {
        requestTemplateGroupService.deleteTemplateGroupByIDService(id);
        return okay();
    }

    @ApiOperationSupport(order = 20)
    @PostMapping("/save")
    @ApiOperation(value = "Request-Info-save")
    public JsonResponse<SaveRequestInfoReq> saveRequestInfo(@RequestBody SaveRequestInfoReq req) throws Exception {
        return okayWithData(requestInfoService.saveRequestInfo(req));
    }

    @ApiOperationSupport(order = 21)
    @PostMapping("/search/{page}/{pageSize}")
    @ApiOperation(value = "Request-Info-search")
    public JsonResponse<QueryResponse<SynthesisRequestInfoResp>>selectSynthesisRequestInfo(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize") @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) SynthesisRequestInfoReq req)
            throws Exception {
        QueryResponse<SynthesisRequestInfoResp> list = requestSynthesisService.selectSynthesisRequestInfoService(page, pageSize,req);
        return JsonResponse.okayWithData(list);
    }


    @ApiOperationSupport(order = 22)
    @PostMapping("/details")
    @ApiOperation(value = "Request-Info-details")
    public JsonResponse<SynthesisRequestInfoFormRequest> selectSynthesisRequestInfoForm(String id)
            throws Exception {
        SynthesisRequestInfoFormRequest synthesisRequestInfoFormRequest = requestSynthesisService.selectSynthesisRequestInfoFormService(id);
        return JsonResponse.okayWithData(synthesisRequestInfoFormRequest);
    }

    @ApiOperationSupport(order = 23)
    @ApiOperation(value = "Request-Info-done")
    @PostMapping(ApplicationConstants.ApiInfo.API_RESOURCE_SERVICE_REQUEST_DONE)
    public JsonResponse updateServiceRequest(@RequestBody DoneServiceRequestRequest request,
                                             @ApiIgnore HttpServletRequest httpRequest) throws Exception {
        requestInfoService.doneServiceRequest(request);
        return okay();
    }

}


