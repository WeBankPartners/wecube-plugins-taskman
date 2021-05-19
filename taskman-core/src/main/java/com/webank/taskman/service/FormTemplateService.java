package com.webank.taskman.service;


import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.LocalPageableQueryResult;
import com.webank.taskman.domain.FormTemplate;
import com.webank.taskman.dto.req.FormTemplateSaveReqDto;
import com.webank.taskman.dto.resp.FormTemplateQueryResultDto;

public interface FormTemplateService extends IService<FormTemplate> {
    LocalPageableQueryResult<FormTemplateQueryResultDto> selectFormTemplate(Integer current, Integer limit, FormTemplateSaveReqDto req);

    void deleteFormTemplate(String id);

    FormTemplateQueryResultDto detailFormTemplate(FormTemplateSaveReqDto saveFormTemplateReq);

    FormTemplateQueryResultDto saveOrUpdateFormTemplate(FormTemplateSaveReqDto formTemplateReq);

    FormTemplateQueryResultDto queryDetailByTemplate(Integer tempType,String tempId);

}
