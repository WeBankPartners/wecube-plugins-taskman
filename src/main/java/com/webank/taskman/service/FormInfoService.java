package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.FormInfo;
import com.webank.taskman.dto.req.SaveFormInfoReq;
import com.webank.taskman.dto.req.SaveFormTemplateReq;


public interface FormInfoService extends IService<FormInfo> {

    void saveFormInfoByReq(SaveFormInfoReq req);
}
