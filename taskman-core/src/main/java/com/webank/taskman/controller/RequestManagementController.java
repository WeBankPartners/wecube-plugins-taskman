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

import com.webank.taskman.base.JsonResponse;
import com.webank.taskman.base.LocalPageableQueryResult;
import com.webank.taskman.dto.CreateTaskDto;
import com.webank.taskman.dto.RequestTemplateDto;
import com.webank.taskman.dto.req.RequestInfoQueryReqDto;
import com.webank.taskman.dto.req.RequestTemplateQueryDto;
import com.webank.taskman.dto.resp.RequestInfoResqDto;
import com.webank.taskman.dto.resp.RequestTemplateQueryResultDto;
import com.webank.taskman.service.RequestInfoService;
import com.webank.taskman.service.RequestTemplateService;

@RestController
@RequestMapping("/v1/request")
public class RequestManagementController {

    @Autowired
    private RequestTemplateService requestTemplateService;

    @Autowired
    private RequestInfoService requestInfoService;

    @PostMapping("/template/save")
    public JsonResponse saveRequestTemplate(@Valid @RequestBody RequestTemplateDto requestTemplateDto) {
        RequestTemplateDto requestTemplateDTO = requestTemplateService.saveRequestTemplate(requestTemplateDto);
        return okayWithData(requestTemplateDTO);
    }

    /**
     * Release template.
     * 
     * @param reqDto
     * @return
     */
    @PostMapping("/template/release")
    public JsonResponse releaseRequestTemplate(@RequestBody RequestTemplateDto requestTemplateDto) {
        RequestTemplateDto respDto = requestTemplateService.releaseRequestTemplate(requestTemplateDto);
        return okayWithData(respDto);
    }

    /**
     * 
     * @param page
     * @param pageSize
     * @param requestTemplateQueryDto
     * @return
     */
    @PostMapping("/template/search/{page}/{page-size}")
    public JsonResponse searchRequestTemplates(@PathVariable("page") Integer page,
            @PathVariable("page-size") Integer pageSize,
            @RequestBody(required = false) RequestTemplateQueryDto requestTemplateQueryDto) {
        LocalPageableQueryResult<RequestTemplateDto> queryResponse = requestTemplateService.searchRequestTemplates(page,
                pageSize, requestTemplateQueryDto);
        return okayWithData(queryResponse);
    }

    /**
     * 
     * @param id
     * @return
     */
    @DeleteMapping("/template/delete/{id}")
    public JsonResponse deleteRequestTemplate(@PathVariable("id") String id) {
        requestTemplateService.deleteRequestTemplate(id);
        return okay();
    }

    /**
     * 
     * @param id
     * @return
     */
    @GetMapping("/template/detail/{id}")
    public JsonResponse fetchRequestTemplateDetail(@PathVariable("id") String id) {
        RequestTemplateQueryResultDto detailRequestTemplateResq = requestTemplateService.fetchRequestTemplateDetail(id);
        return okayWithData(detailRequestTemplateResq);
    }

    /**
     * 
     * @return
     */
    @GetMapping("/template/available")
    public JsonResponse fetchAvailableRequestTemplates() {
        List<RequestTemplateDto> retRequestTemplateDtos = requestTemplateService.fetchAvailableRequestTemplates();
        return okayWithData(retRequestTemplateDtos);
    }

    /**
     * Submit new request
     * 
     * @param req
     * @return
     */
    @PostMapping("/save")
    public JsonResponse createNewRequestInfo(@RequestBody CreateTaskDto req) {
        return okayWithData(requestInfoService.createNewRequestInfo(req));
    }

    @PostMapping("/search/{page}/{page-size}")
    public JsonResponse requestInfoSearch(@PathVariable("page") Integer page,
            @PathVariable("page-size") Integer pageSize, @RequestBody(required = false) RequestInfoQueryReqDto req) {
        LocalPageableQueryResult<RequestInfoResqDto> list = requestInfoService.selectRequestInfoPage(page, pageSize, req);
        return okayWithData(list);
    }

    @GetMapping("/details/{id}")
    public JsonResponse requestInfoDetail(@PathVariable("id") String id) {
        RequestInfoResqDto requestInfoResq = requestInfoService.selectDetail(id);
        return okayWithData(requestInfoResq);
    }

}
