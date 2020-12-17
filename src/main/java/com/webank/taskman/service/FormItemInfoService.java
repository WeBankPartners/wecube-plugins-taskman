package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.FormItemInfo;
import com.webank.taskman.dto.req.SaveFormTemplateReq;

import java.util.List;


public interface FormItemInfoService extends IService<FormItemInfo> {

    void saveFormItemInfoByReq(SaveFormTemplateReq req);

    List<FormItemInfo> selectFormItemInfo(String requestTempId);
}
