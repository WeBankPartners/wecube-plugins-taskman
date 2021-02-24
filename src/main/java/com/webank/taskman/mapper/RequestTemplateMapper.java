package com.webank.taskman.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.req.QueryRequestTemplateReq;
import com.webank.taskman.dto.RequestTemplateDto;
import org.apache.ibatis.annotations.Param;

import java.util.List;


public interface RequestTemplateMapper extends BaseMapper<RequestTemplate> {


    List<RequestTemplateDto>  selectDTOListByParam(QueryRequestTemplateReq req);


}
