package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.commons.TaskmanException;
import com.webank.taskman.domain.FormTemplate;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.SaveFormTemplateReq;
import com.webank.taskman.dto.resp.FormTemplateResp;


public interface FormTemplateService extends IService<FormTemplate> {
    QueryResponse<FormTemplateResp> selectFormTemplate(Integer current, Integer limit, SaveFormTemplateReq req) throws Exception;

    void deleteFormTemplate(String id) throws TaskmanException;

    FormTemplateResp detailFormTemplate(SaveFormTemplateReq saveFormTemplateReq) throws TaskmanException;

    FormTemplateResp saveFormTemplateByReq(SaveFormTemplateReq formTemplateReq) throws TaskmanException;

    FormTemplateResp queryDetailByTemp(Integer tempType,String tempId) throws TaskmanException;

}
