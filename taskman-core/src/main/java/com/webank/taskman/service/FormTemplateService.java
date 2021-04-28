package com.webank.taskman.service;


import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.PageableQueryResult;
import com.webank.taskman.domain.FormTemplate;
import com.webank.taskman.dto.req.FormTemplateSaveReqDto;
import com.webank.taskman.dto.resp.FormTemplateRespDto;

public interface FormTemplateService extends IService<FormTemplate> {
    PageableQueryResult<FormTemplateRespDto> selectFormTemplate(Integer current, Integer limit, FormTemplateSaveReqDto req);

    void deleteFormTemplate(String id);

    FormTemplateRespDto detailFormTemplate(FormTemplateSaveReqDto saveFormTemplateReq);

    FormTemplateRespDto saveOrUpdateFormTemplate(FormTemplateSaveReqDto formTemplateReq);

    FormTemplateRespDto queryDetailByTemplate(Integer tempType,String tempId);

}
