package com.webank.taskman.service.impl;


import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.domain.FormItemInfo;
import com.webank.taskman.dto.resp.FormItemInfoResp;
import com.webank.taskman.mapper.FormItemInfoMapper;
import com.webank.taskman.mapper.FormItemTemplateMapper;
import com.webank.taskman.service.FormItemInfoService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;


@Service
public class FormItemInfoServiceImpl extends ServiceImpl<FormItemInfoMapper, FormItemInfo> implements FormItemInfoService {

    @Autowired
    FormItemTemplateMapper formItemTemplateMapper;


    @Override
    public List<FormItemInfoResp> returnDetail(String id) {
        List<FormItemInfoResp> formItemInfoResps = formItemTemplateMapper.selectDetail(id);
        return formItemInfoResps;
    }


    @Override
    public void saveItemInfoByList(List<FormItemInfo> formItems, String recordId, String formId) {
        if (null != formItems && formItems.size() > 0) {
            formItems.stream().forEach(item -> {
                item.setFormId(formId);
                item.setRecordId(recordId);
                save(item);
            });
        }
    }
}
