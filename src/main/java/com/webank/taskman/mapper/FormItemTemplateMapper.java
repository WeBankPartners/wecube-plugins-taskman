package com.webank.taskman.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.webank.taskman.domain.FormItemTemplate;


public interface FormItemTemplateMapper extends BaseMapper<FormItemTemplate> {

    void deleteRequestTemplateByIDMapper(String id);

    int deleteByDomain(FormItemTemplate formItemTemplate);
}
