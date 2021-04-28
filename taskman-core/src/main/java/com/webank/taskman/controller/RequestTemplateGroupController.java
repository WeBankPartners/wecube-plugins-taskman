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
import com.webank.taskman.dto.RequestTemplateGroupDto;
import com.webank.taskman.service.RequestTemplateGroupService;

@RestController
@RequestMapping("/v1/request")
public class RequestTemplateGroupController {
    @Autowired
    private RequestTemplateGroupService requestTemplateGroupService;

    /**
     * 
     * @param req
     * @return
     */
    @PostMapping("/template/group/save")
    public JsonResponse saveOrUpdateRequestTemplateGroup(@Valid @RequestBody RequestTemplateGroupDto requestTemplateGroupDto) {
        RequestTemplateGroupDto savedRequestTemplateGroupDto = requestTemplateGroupService
                .saveOrUpdateTemplateGroup(requestTemplateGroupDto);
        return okayWithData(savedRequestTemplateGroupDto);
    }

    /**
     * 
     * @param page
     * @param pageSize
     * @param req
     * @return
     */
    @PostMapping("/template/group/search/{page}/{page-size}")
    public JsonResponse searchRequestTemplateGroups(@PathVariable("page") Integer page,
            @PathVariable("page-size") Integer pageSize, @RequestBody(required = false) RequestTemplateGroupDto requestTemplateGroupSearchCriteriaDto) {
        return okayWithData(requestTemplateGroupService.searchRequestTemplateGroups(page, pageSize, requestTemplateGroupSearchCriteriaDto));
    }

    /**
     * 
     * @return
     */
    @GetMapping("/template/group/available")
    public JsonResponse fetchAvailableGroupTemplates() {
        List<RequestTemplateGroupDto> requestTemplateGroupDtos = requestTemplateGroupService
                .fetchAvailableGroupTemplates();
        return okayWithData(requestTemplateGroupDtos);
    }

    /**
     * 
     * @param id
     * @return
     */
    @DeleteMapping("/template/group/delete/{id}")
    public JsonResponse requestGroupTemplateDelete(@PathVariable("id") String id) {
        requestTemplateGroupService.deleteRequestTemplateGroup(id);
        return okay();
    }

}
