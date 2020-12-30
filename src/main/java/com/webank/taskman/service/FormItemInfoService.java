package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.FormItemInfo;
import com.webank.taskman.dto.req.SaveFormItemInfoReq;

import java.util.List;


public interface FormItemInfoService extends IService<FormItemInfo> {

    void saveItemInfoByList(List<FormItemInfo> formItems, String recordId, String formId);
}
