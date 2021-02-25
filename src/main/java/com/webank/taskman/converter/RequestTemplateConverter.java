package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.req.RequestTemplateSaveReqDto;
import com.webank.taskman.dto.RequestTemplateDto;
import com.webank.taskman.dto.resp.RequestTemplateRespDto;
import org.mapstruct.Mapper;
import org.mapstruct.ReportingPolicy;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface RequestTemplateConverter extends BaseConverter<RequestTemplateDto, RequestTemplate> {


    RequestTemplate saveReqToEntity(RequestTemplateSaveReqDto req);

    RequestTemplateRespDto toRespByEntity(RequestTemplate req);




}