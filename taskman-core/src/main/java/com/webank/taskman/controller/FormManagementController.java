package com.webank.taskman.controller;

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

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.webank.taskman.base.JsonResponse;
import com.webank.taskman.converter.FormItemTemplateConverter;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.FormItemTemplateDto;
import com.webank.taskman.dto.req.FormTemplateSaveReqDto;
import com.webank.taskman.dto.resp.FormTemplateRespDto;
import com.webank.taskman.service.FormItemTemplateService;
import com.webank.taskman.service.FormTemplateService;

@RestController
@RequestMapping("/v1/form")
public class FormManagementController {

    @Autowired
    private FormItemTemplateService formItemTemplateService;

    @Autowired
    private FormTemplateService formTemplateService;

    @Autowired
    private FormItemTemplateConverter formItemTemplateConverter;

    @PostMapping("/template/save")
    public JsonResponse createFormTemplate(@Valid @RequestBody FormTemplateSaveReqDto reqDto) {

        FormTemplateRespDto formTemplateResp = formTemplateService.saveOrUpdateFormTemplate(reqDto);
        return JsonResponse.okayWithData(formTemplateResp);
    }

    @DeleteMapping("/template/delete/{id}")
    public JsonResponse formTemplateDelete(@PathVariable("id") String id) {
        formTemplateService.deleteFormTemplate(id);
        return JsonResponse.okay();
    }

    @GetMapping("/template/detail/{temp-type}/{temp-id}")
    public JsonResponse fetchFormTemplateDetail(@PathVariable("temp-type") String tempType,
            @PathVariable("temp-id") String tempId) {
        
        FormTemplateSaveReqDto criteriaDto = new FormTemplateSaveReqDto(tempId, tempType);
        FormTemplateRespDto formTemplateDto = formTemplateService.detailFormTemplate(criteriaDto);
        return JsonResponse
                .okayWithData(formTemplateDto);
    }

    @DeleteMapping("/item/delete/{id}")
    public JsonResponse formItemTemplateDelete(@PathVariable("id") String id) {
        formItemTemplateService.deleteRequestTemplateByID(id);
        return JsonResponse.okay();
    }

    @PostMapping("/item/template/currency")
    public JsonResponse formItemTemplateAvailable() {
        QueryWrapper<FormItemTemplate> wrapper = new QueryWrapper<FormItemTemplate>();
        wrapper.eq("status", 1);
        
        List<FormItemTemplate> formItemTemplateEntities = formItemTemplateService.list(wrapper);
        List<FormItemTemplateDto> queryResponse = formItemTemplateConverter
                .convertToDtos(formItemTemplateEntities);
        return JsonResponse.okayWithData(queryResponse);
    }

}
