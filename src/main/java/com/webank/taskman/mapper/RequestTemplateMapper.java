package com.webank.taskman.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.req.QueryRequestTemplateReq;


public interface RequestTemplateMapper extends BaseMapper<RequestTemplate> {


    IPage<RequestTemplate> selectPageByParam(Page page, QueryRequestTemplateReq queryReq);
}
