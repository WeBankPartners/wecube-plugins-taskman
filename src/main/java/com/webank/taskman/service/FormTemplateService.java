package com.webank.taskman.service;


import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.domain.FormTemplate;
import com.webank.taskman.dto.req.SaveFormTemplateReq;
import com.webank.taskman.dto.resp.FormTemplateResp;

public interface FormTemplateService extends IService<FormTemplate> {
    QueryResponse<FormTemplateResp> selectFormTemplate(Integer current, Integer limit, SaveFormTemplateReq req);

    void deleteFormTemplate(String id);

    FormTemplateResp detailFormTemplate(SaveFormTemplateReq saveFormTemplateReq);

    FormTemplateResp saveFormTemplateByReq(SaveFormTemplateReq formTemplateReq);

    FormTemplateResp queryDetailByTemp(Integer tempType,String tempId);

}
