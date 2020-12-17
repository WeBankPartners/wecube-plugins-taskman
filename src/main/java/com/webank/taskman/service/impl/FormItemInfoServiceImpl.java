package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.domain.FormItemInfo;
import com.webank.taskman.dto.req.SaveFormTemplateReq;
import com.webank.taskman.mapper.FormItemInfoMapper;
import com.webank.taskman.service.FormItemInfoService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;


@Service
public class FormItemInfoServiceImpl extends ServiceImpl<FormItemInfoMapper, FormItemInfo> implements FormItemInfoService {

    @Autowired
    FormItemInfoMapper formItemInfoMapper;

    @Override
    public void saveFormItemInfoByReq(SaveFormTemplateReq req) {

    }

    @Override
    public List<FormItemInfo> selectFormItemInfo(String requestTempId) {
        List<FormItemInfo> list=formItemInfoMapper.selectFormItemInfo(requestTempId);
        return list;
    }
}
