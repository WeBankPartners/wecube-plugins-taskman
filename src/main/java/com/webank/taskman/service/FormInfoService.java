package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.domain.FormInfo;
import com.webank.taskman.domain.FormItemInfo;

import java.util.List;


public interface FormInfoService extends IService<FormInfo> {


    FormInfo saveFormInfoByExists(String requestTempId, String requestInfoId) throws TaskmanRuntimeException;

    void saveFormInfoAndItems(List<FormItemInfo> formItems, String templateId, String recordId);

}
