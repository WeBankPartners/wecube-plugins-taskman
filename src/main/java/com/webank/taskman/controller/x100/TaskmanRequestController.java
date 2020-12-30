package com.webank.taskman.controller.x100;


import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.base.JsonResponse;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanException;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.StatusCodeEnum;
import com.webank.taskman.constant.StatusEnum;
import com.webank.taskman.converter.RequestTemplateConverter;
import com.webank.taskman.converter.RequestTemplateGroupConverter;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.RequestTemplateDTO;
import com.webank.taskman.dto.RequestTemplateGroupDTO;
import com.webank.taskman.dto.req.*;
import com.webank.taskman.dto.resp.*;
import com.webank.taskman.service.*;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.ApiParam;
import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import springfox.documentation.annotations.ApiIgnore;

import javax.validation.Valid;
import java.util.Date;
import java.util.List;

import static com.webank.taskman.base.JsonResponse.okay;
import static com.webank.taskman.base.JsonResponse.okayWithData;


@Api(tags = {"3、 Request inteface API"})
@RestController
@RequestMapping("/v1/request")
public class TaskmanRequestController {

    private static final Logger log = LoggerFactory.getLogger(TaskmanRequestController.class);

    @Autowired
    RequestTemplateService requestTemplateService;

    @Autowired
    RequestTemplateGroupService requestTemplateGroupService;

    @Autowired
    RequestTemplateGroupConverter requestTemplateGroupConverter;

    @Autowired
    TaskInfoService taskInfoService;



    @Autowired
    RequestInfoService requestInfoService;

    @ApiOperationSupport(order = 1)
    @PostMapping("/template/group/save")
    @ApiOperation(value = "request-group-template-save", notes = "")
    public JsonResponse<RequestTemplateGroupDTO> requestGroupTemplateSave(
            @Valid @RequestBody SaveRequestTemplateGropReq req) throws TaskmanException {
        return JsonResponse.okayWithData(requestTemplateGroupService.saveTemplateGroupByReq(req));
    }

    @ApiOperationSupport(order = 2,ignoreParameters = {"manageRoleName"})
    @PostMapping("/template/group/search/{page}/{pageSize}")
    @ApiOperation(value = "request-group-template-search")
    public JsonResponse<QueryResponse<RequestTemplateGroupDTO>> requestGroupTemplateSearch(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize")  @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) RequestTemplateGroupDTO req
    ) throws TaskmanRuntimeException
    {
        return JsonResponse.okayWithData(requestTemplateGroupService.selectRequestTemplateGroupPage(page, pageSize, req));
    }

    @ApiOperationSupport(order = 3)
    @GetMapping("/template/group/available")
    @ApiOperation(value = "request-group-template-available")
    public JsonResponse<List<RequestTemplateGroupDTO>> requestGroupTemplateAvailable() throws TaskmanRuntimeException
    {
        LambdaQueryWrapper lambdaQueryWrapper = new RequestTemplateGroup().setStatus(StatusEnum.DEFAULT.toString()).getLambdaQueryWrapper();
        List<RequestTemplateGroupDTO> dtoList = requestTemplateGroupConverter.toDto(requestTemplateGroupService.list(lambdaQueryWrapper));
        return JsonResponse.okayWithData(dtoList);
    }

    @ApiOperationSupport(order = 4)
    @DeleteMapping("/template/group/delete/{id}")
    @ApiOperation(value = "request-group-template-delete", notes = "")
    public JsonResponse requestGroupTemplateDelete(@PathVariable("id") String id) throws TaskmanRuntimeException
    {
        requestTemplateGroupService.deleteTemplateGroupByIDService(id);
        return okay();
    }


    @ApiOperationSupport(order = 5)
    @PostMapping("/template/save")
    @ApiOperation(value = "request-template-save", notes = "Need to pass in object: ")
    public JsonResponse requestTemplateSave(@Valid @RequestBody SaveRequestTemplateReq req
    )throws TaskmanRuntimeException {
        RequestTemplateDTO requestTemplateDTO = requestTemplateService.saveRequestTemplate(req);
        return JsonResponse.okayWithData(requestTemplateDTO);
    }

    @ApiOperationSupport(order = 6,ignoreParameters = {"requestTempGroup","procDefId","procDefName","description","name","tags","manageRoles","useRoles"})
    @PostMapping("/template/release")
    @ApiOperation(value = "request-template-release", notes = "Need to pass in object: ")
    public JsonResponse<RequestTemplateDTO> requestTemplateRelease(@RequestBody SaveRequestTemplateReq req) throws TaskmanRuntimeException {
        if(StringUtils.isEmpty(req.getId())){
            return  JsonResponse.customError(StatusCodeEnum.PARAM_ISNULL);
        }
        RequestTemplate requestTemplate = requestTemplateService.getOne(new RequestTemplate().setId(req.getId()).getLambdaQueryWrapper());
        if(null == requestTemplate){
            return JsonResponse.customError(StatusCodeEnum.NOT_FOUND_RECORD);
        }
        requestTemplate.setStatus(StatusEnum.UNRELEASED.toString().equals(requestTemplate.getStatus())?
                StatusEnum.RELEASED.toString() :StatusEnum.UNRELEASED.toString());
        requestTemplate.setUpdatedBy(AuthenticationContextHolder.getCurrentUsername());
        requestTemplate.setUpdatedTime(new Date());
        requestTemplateService.updateById(requestTemplate);
        return okayWithData(new RequestTemplateDTO().setId(requestTemplate.getId()).setStatus(requestTemplate.getStatus()));
    }

    @ApiOperationSupport(order = 7)
    @PostMapping("/template/search/{page}/{pageSize}")
    @ApiOperation(value = "request-template-search")
    public JsonResponse<QueryResponse<RequestTemplateDTO>> requestTemplateSearch(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize")  @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) QueryRequestTemplateReq req)
            throws TaskmanRuntimeException {
        QueryResponse<RequestTemplateDTO> queryResponse = requestTemplateService.selectRequestTemplatePage(page, pageSize, req);
        return JsonResponse.okayWithData(queryResponse);
    }

    @ApiOperationSupport(order = 8)
    @DeleteMapping("/template/delete/{id}")
    @ApiOperation(value = "request-template-delete", notes = "需要传入id")
    public JsonResponse requestTemplateDelete(@PathVariable("id") String id) throws TaskmanRuntimeException {
        requestTemplateService.deleteRequestTemplateService(id);
        return okay();
    }

    @ApiOperationSupport(order = 9)
    @GetMapping("/template/detail/{id}")
    @ApiOperation(value = "request-template-detail", notes = "需要传入id")
    public JsonResponse<DetailRequestTemplateResq> requestTemplateDetail(@PathVariable("id") String id) throws TaskmanRuntimeException {
        DetailRequestTemplateResq detailRequestTemplateResq= requestTemplateService.detailRequestTemplate(id);
        return JsonResponse.okayWithData(detailRequestTemplateResq);
    }

    @Autowired
    RequestTemplateConverter requestTemplateConverter;

    @GetMapping(value = "/template/available")
    @ApiOperationSupport(order = 10)
    @ApiOperation(value = "request-template-available")
    public JsonResponse<List<RequestTemplateDTO>> requestTemplateAvailable(@ApiIgnore QueryRoleRelationBaseReq req) throws TaskmanRuntimeException {
        RequestTemplate requestTemplate = new RequestTemplate().setStatus(StatusEnum.RELEASED.toString());
        LambdaQueryWrapper<RequestTemplate> queryWrapper = requestTemplate.getLambdaQueryWrapper()
                .inSql(RequestTemplate::getId,req.getEqUseRole());
        return okayWithData(requestTemplateConverter.toDto(requestTemplateService.list(queryWrapper)));
    }



    @ApiOperationSupport(order = 22)
    @PostMapping("/save")
    @ApiOperation(value = "request-info-save")
    public JsonResponse<SaveRequestInfoReq> requestInfoSave(@RequestBody SaveRequestInfoReq req) throws TaskmanRuntimeException {
        return okayWithData(requestInfoService.saveRequestInfo(req));
    }

    @ApiOperationSupport(order = 12)
    @PostMapping("/search/{page}/{pageSize}")
    @ApiOperation(value = "request-info-search")
    public JsonResponse<QueryResponse<RequestInfoResq>> requestInfoSearch(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize") @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) QueryRequestInfoReq req)
            throws TaskmanRuntimeException {
        QueryResponse<RequestInfoResq> list = requestInfoService.selectRequestInfoPage(page, pageSize,req);
        return JsonResponse.okayWithData(list);
    }


    @ApiOperationSupport(order = 13)
    @GetMapping("/details/{id}")
    @ApiOperation(value = "request-info-detail")
    public JsonResponse<RequestInfoResq> requestInfoDetail(@PathVariable("id") String id)
            throws TaskmanRuntimeException {
        RequestInfoResq requestInfoResq = requestInfoService.selectDetail(id);
        return JsonResponse.okayWithData(requestInfoResq);
    }

}


