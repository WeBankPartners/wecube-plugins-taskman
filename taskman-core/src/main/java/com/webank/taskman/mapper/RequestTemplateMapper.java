package com.webank.taskman.mapper;

import java.util.List;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.RequestTemplateDto;
import com.webank.taskman.dto.req.RequestTemplateQueryReqDto;


public interface RequestTemplateMapper extends BaseMapper<RequestTemplate> {


    List<RequestTemplateDto>  selectDTOListByParam(RequestTemplateQueryReqDto req);


}
