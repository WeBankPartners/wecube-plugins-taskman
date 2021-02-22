package com.webank.taskman.controller;

import static com.webank.taskman.base.JsonResponse.okay;
import static com.webank.taskman.base.JsonResponse.okayWithData;

import java.util.Date;
import java.util.List;

import javax.validation.Valid;

import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
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
import com.webank.taskman.dto.CreateTaskDto;
import com.webank.taskman.dto.RequestTemplateDTO;
import com.webank.taskman.dto.RequestTemplateGroupDTO;
import com.webank.taskman.dto.req.QueryRequestInfoReq;
import com.webank.taskman.dto.req.QueryRequestTemplateReq;
import com.webank.taskman.dto.req.QueryRoleRelationBaseReq;
import com.webank.taskman.dto.req.SaveRequestTemplateGropReq;
import com.webank.taskman.dto.req.SaveRequestTemplateReq;
import com.webank.taskman.dto.resp.RequestInfoResq;
import com.webank.taskman.dto.resp.RequestTemplateResp;
import com.webank.taskman.service.RequestInfoService;
import com.webank.taskman.service.RequestTemplateGroupService;
import com.webank.taskman.service.RequestTemplateService;

@RestController
@RequestMapping("/v1/request")
public class TaskmanRequestController {

    @Autowired
    private RequestTemplateService requestTemplateService;

    @Autowired
    private RequestTemplateGroupService requestTemplateGroupService;

    @Autowired
    private RequestTemplateGroupConverter requestTemplateGroupConverter;

    @Autowired
    private RequestInfoService requestInfoService;

    @PostMapping("/template/group/save")
    public JsonResponse requestGroupTemplateSave(@Valid @RequestBody SaveRequestTemplateGropReq req)
            throws TaskmanException {
        return JsonResponse.okayWithData(requestTemplateGroupService.saveTemplateGroupByReq(req));
    }

    @PostMapping("/template/group/search/{page}/{pageSize}")
    public JsonResponse requestGroupTemplateSearch(@PathVariable("page") Integer page,
            @PathVariable("pageSize") Integer pageSize, @RequestBody(required = false) RequestTemplateGroupDTO req)
            throws TaskmanRuntimeException {
        return JsonResponse
                .okayWithData(requestTemplateGroupService.selectRequestTemplateGroupPage(page, pageSize, req));
    }

    @GetMapping("/template/group/available")
    public JsonResponse requestGroupTemplateAvailable() throws TaskmanRuntimeException {
        LambdaQueryWrapper lambdaQueryWrapper = new RequestTemplateGroup().setStatus(StatusEnum.DEFAULT.toString())
                .getLambdaQueryWrapper();
        List<RequestTemplateGroupDTO> dtoList = requestTemplateGroupConverter
                .toDto(requestTemplateGroupService.list(lambdaQueryWrapper));
        return JsonResponse.okayWithData(dtoList);
    }

    @DeleteMapping("/template/group/delete/{id}")
    public JsonResponse requestGroupTemplateDelete(@PathVariable("id") String id) throws TaskmanRuntimeException {
        requestTemplateGroupService.deleteTemplateGroupByIDService(id);
        return okay();
    }

    @PostMapping("/template/save")
    public JsonResponse requestTemplateSave(@Valid @RequestBody SaveRequestTemplateReq req)
            throws TaskmanRuntimeException {
        RequestTemplateDTO requestTemplateDTO = requestTemplateService.saveRequestTemplate(req);
        return JsonResponse.okayWithData(requestTemplateDTO);
    }

    @PostMapping("/template/release")
    public JsonResponse requestTemplateRelease(@RequestBody SaveRequestTemplateReq req) throws TaskmanRuntimeException {
        if (StringUtils.isEmpty(req.getId())) {
            return JsonResponse.customError(StatusCodeEnum.PARAM_ISNULL);
        }
        RequestTemplate requestTemplate = requestTemplateService
                .getOne(new RequestTemplate().setId(req.getId()).getLambdaQueryWrapper());
        if (null == requestTemplate) {
            return JsonResponse.customError(StatusCodeEnum.NOT_FOUND_RECORD);
        }
        requestTemplate.setStatus(StatusEnum.UNRELEASED.toString().equals(requestTemplate.getStatus())
                ? StatusEnum.RELEASED.toString() : StatusEnum.UNRELEASED.toString());
        requestTemplate.setUpdatedBy(AuthenticationContextHolder.getCurrentUsername());
        requestTemplate.setUpdatedTime(new Date());
        requestTemplateService.updateById(requestTemplate);
        return okayWithData(
                new RequestTemplateDTO().setId(requestTemplate.getId()).setStatus(requestTemplate.getStatus()));
    }

    @PostMapping("/template/search/{page}/{pageSize}")
    public JsonResponse requestTemplateSearch(@PathVariable("page") Integer page,
            @PathVariable("pageSize") Integer pageSize, @RequestBody(required = false) QueryRequestTemplateReq req)
            throws TaskmanRuntimeException {
        QueryResponse<RequestTemplateDTO> queryResponse = requestTemplateService.selectRequestTemplatePage(page,
                pageSize, req);
        return JsonResponse.okayWithData(queryResponse);
    }

    @DeleteMapping("/template/delete/{id}")
    public JsonResponse requestTemplateDelete(@PathVariable("id") String id) throws TaskmanRuntimeException {
        requestTemplateService.deleteRequestTemplateService(id);
        return okay();
    }

    @GetMapping("/template/detail/{id}")
    public JsonResponse requestTemplateDetail(@PathVariable("id") String id) throws TaskmanRuntimeException {
        RequestTemplateResp detailRequestTemplateResq = requestTemplateService.detailRequestTemplate(id);
        return JsonResponse.okayWithData(detailRequestTemplateResq);
    }

    @Autowired
    RequestTemplateConverter requestTemplateConverter;

    @GetMapping(value = "/template/available")
    public JsonResponse requestTemplateAvailable() {
        RequestTemplate requestTemplate = new RequestTemplate().setStatus(StatusEnum.RELEASED.toString());
        String inSql = QueryRoleRelationBaseReq.getEqUseRole();
        LambdaQueryWrapper<RequestTemplate> queryWrapper = requestTemplate.getLambdaQueryWrapper()
                .inSql(!StringUtils.isEmpty(inSql), RequestTemplate::getId, QueryRoleRelationBaseReq.getEqUseRole());
        return okayWithData(requestTemplateConverter.toDto(requestTemplateService.list(queryWrapper)));
    }

    @PostMapping("/save")
    public JsonResponse requestInfoSave(@RequestBody CreateTaskDto req) throws TaskmanRuntimeException {
        return okayWithData(requestInfoService.saveRequestInfoByDto(req));
    }

    @PostMapping("/search/{page}/{pageSize}")
    public JsonResponse requestInfoSearch(@PathVariable("page") Integer page,
            @PathVariable("pageSize") Integer pageSize, @RequestBody(required = false) QueryRequestInfoReq req)
            throws TaskmanRuntimeException {
        QueryResponse<RequestInfoResq> list = requestInfoService.selectRequestInfoPage(page, pageSize, req);
        return JsonResponse.okayWithData(list);
    }

    @GetMapping("/details/{id}")
    public JsonResponse requestInfoDetail(@PathVariable("id") String id) throws TaskmanRuntimeException {
        RequestInfoResq requestInfoResq = requestInfoService.selectDetail(id);
        return JsonResponse.okayWithData(requestInfoResq);
    }

}
