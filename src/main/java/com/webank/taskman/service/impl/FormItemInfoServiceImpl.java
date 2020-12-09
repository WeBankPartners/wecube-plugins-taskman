package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.domain.FormItemInfo;
import com.webank.taskman.dto.req.SaveFormTemplateReq;
import com.webank.taskman.mapper.FormItemInfoMapper;
import com.webank.taskman.service.FormItemInfoService;
import org.springframework.stereotype.Service;


@Service
public class FormItemInfoServiceImpl extends ServiceImpl<FormItemInfoMapper, FormItemInfo> implements FormItemInfoService {


    @Override
    public void saveFormItemInfoByReq(SaveFormTemplateReq req) {

    }
}
