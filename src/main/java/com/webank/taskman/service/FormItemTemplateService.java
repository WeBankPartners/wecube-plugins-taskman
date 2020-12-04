package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.SaveAndUpdateFormItemTemplateReq;
import com.webank.taskman.dto.req.SelectFormItemTemplateReq;
import com.webank.taskman.dto.resp.FormItemTemplateResq;


public interface FormItemTemplateService extends IService<FormItemTemplate> {


    FormItemTemplate addOrUpdateFormItemTemplate(SaveAndUpdateFormItemTemplateReq templateReq) throws Exception;

    void deleteRequestTemplateByID(String id);

    QueryResponse<FormItemTemplateResq> selectAllFormItemTemplateService(Integer current, Integer limit, SelectFormItemTemplateReq req);
}
