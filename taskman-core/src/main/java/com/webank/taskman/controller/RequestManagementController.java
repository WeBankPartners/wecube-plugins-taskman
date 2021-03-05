package com.webank.taskman.controller;

import static com.webank.taskman.base.JsonResponse.okay;
import static com.webank.taskman.base.JsonResponse.okayWithData;

import java.util.List;

import javax.validation.Valid;

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
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.converter.RequestTemplateGroupConverter;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.CreateTaskDto;
import com.webank.taskman.dto.RequestTemplateDto;
import com.webank.taskman.dto.RequestTemplateGroupDto;
import com.webank.taskman.dto.req.RequestInfoQueryReqDto;
import com.webank.taskman.dto.req.RequestTemplateGroupSaveReqDto;
import com.webank.taskman.dto.req.RequestTemplateQueryReqDto;
import com.webank.taskman.dto.req.RequestTemplateSaveReqDto;
import com.webank.taskman.dto.resp.RequestInfoResqDto;
import com.webank.taskman.dto.resp.RequestTemplateRespDto;
import com.webank.taskman.service.RequestInfoService;
import com.webank.taskman.service.RequestTemplateGroupService;
import com.webank.taskman.service.RequestTemplateService;

@RestController
@RequestMapping("/v1/request")
public class RequestManagementController {

    @Autowired
    private RequestTemplateService requestTemplateService;

    @Autowired
    private RequestTemplateGroupService requestTemplateGroupService;

    @Autowired
    private RequestTemplateGroupConverter requestTemplateGroupConverter;

    @Autowired
    private RequestInfoService requestInfoService;

    @PostMapping("/template/group/save")
    public JsonResponse requestGroupTemplateSave(@Valid @RequestBody RequestTemplateGroupSaveReqDto req) {
        return okayWithData(requestTemplateGroupService.saveTemplateGroupByReq(req));
    }

    @PostMapping("/template/group/search/{page}/{pageSize}")
    public JsonResponse requestGroupTemplateSearch(@PathVariable("page") Integer page,
            @PathVariable("pageSize") Integer pageSize, @RequestBody(required = false) RequestTemplateGroupDto req)
            throws TaskmanRuntimeException {
        return okayWithData(requestTemplateGroupService.selectRequestTemplateGroupPage(page, pageSize, req));
    }

    @GetMapping("/template/group/available")
    public JsonResponse requestGroupTemplateAvailable() {
        LambdaQueryWrapper lambdaQueryWrapper = new RequestTemplateGroup()
                .setStatus(RequestTemplateGroup.STATUS_AVAILABLE).getLambdaQueryWrapper();
        List<RequestTemplateGroupDto> dtoList = requestTemplateGroupConverter
                .toDto(requestTemplateGroupService.list(lambdaQueryWrapper));
        return okayWithData(dtoList);
    }

    @DeleteMapping("/template/group/delete/{id}")
    public JsonResponse requestGroupTemplateDelete(@PathVariable("id") String id) {
        requestTemplateGroupService.deleteTemplateGroupByIDService(id);
        return okay();
    }

    @PostMapping("/template/save")
    public JsonResponse requestTemplateSave(@Valid @RequestBody RequestTemplateSaveReqDto req) {
        RequestTemplateDto requestTemplateDTO = requestTemplateService.saveRequestTemplate(req);
        return okayWithData(requestTemplateDTO);
    }

    /**
     * Release template.
     * 
     * @param reqDto
     * @return
     */
    @PostMapping("/template/release")
    public JsonResponse releaseRequestTemplate(@RequestBody RequestTemplateSaveReqDto reqDto) {
        RequestTemplateDto respDto = requestTemplateService.releaseRequestTemplate(reqDto);
        return okayWithData(respDto);
    }

    @PostMapping("/template/search/{page}/{pageSize}")
    public JsonResponse requestTemplateSearch(@PathVariable("page") Integer page,
            @PathVariable("pageSize") Integer pageSize, @RequestBody(required = false) RequestTemplateQueryReqDto req) {
        QueryResponse<RequestTemplateDto> queryResponse = requestTemplateService.selectRequestTemplatePage(page,
                pageSize, req);
        return okayWithData(queryResponse);
    }

    @DeleteMapping("/template/delete/{id}")
    public JsonResponse requestTemplateDelete(@PathVariable("id") String id) {
        requestTemplateService.deleteRequestTemplateService(id);
        return okay();
    }

    @GetMapping("/template/detail/{id}")
    public JsonResponse requestTemplateDetail(@PathVariable("id") String id) {
        RequestTemplateRespDto detailRequestTemplateResq = requestTemplateService.detailRequestTemplate(id);
        return okayWithData(detailRequestTemplateResq);
    }

    @GetMapping(value = "/template/available")
    public JsonResponse requestTemplateAvailable() {
        List<RequestTemplateDto> retRequestTemplateDtos = requestTemplateService.fetchAvailableRequestTemplates();
        return okayWithData(retRequestTemplateDtos);
    }

    @PostMapping("/save")
    public JsonResponse requestInfoSave(@RequestBody CreateTaskDto req) {
        return okayWithData(requestInfoService.saveRequestInfoByDto(req));
    }

    @PostMapping("/search/{page}/{pageSize}")
    public JsonResponse requestInfoSearch(@PathVariable("page") Integer page,
            @PathVariable("pageSize") Integer pageSize, @RequestBody(required = false) RequestInfoQueryReqDto req) {
        QueryResponse<RequestInfoResqDto> list = requestInfoService.selectRequestInfoPage(page, pageSize, req);
        return okayWithData(list);
    }

    @GetMapping("/details/{id}")
    public JsonResponse requestInfoDetail(@PathVariable("id") String id) {
        RequestInfoResqDto requestInfoResq = requestInfoService.selectDetail(id);
        return okayWithData(requestInfoResq);
    }

}
