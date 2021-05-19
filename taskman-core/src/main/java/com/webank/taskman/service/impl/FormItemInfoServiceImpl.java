package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.domain.FormItemInfo;
import com.webank.taskman.dto.resp.FormItemInfoQueryResultDto;
import com.webank.taskman.mapper.FormItemInfoMapper;
import com.webank.taskman.mapper.FormItemTemplateMapper;
import com.webank.taskman.service.FormItemInfoService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class FormItemInfoServiceImpl extends ServiceImpl<FormItemInfoMapper, FormItemInfo>
        implements FormItemInfoService {

    @Autowired
    private FormItemTemplateMapper formItemTemplateMapper;

    @Override
    public List<FormItemInfoQueryResultDto> returnDetail(String id) {
        List<FormItemInfoQueryResultDto> formItemInfoResps = formItemTemplateMapper.selectDetail(id);
        return formItemInfoResps;
    }

    @Override
    public void saveItemInfoByList(List<FormItemInfo> formItems, String recordId, String formId) {
        if (formItems == null || formItems.isEmpty()) {
            return;
        }

        formItems.forEach(item -> {
            item.setFormId(formId);
            item.setRecordId(recordId);
            save(item);
        });
    }
}
