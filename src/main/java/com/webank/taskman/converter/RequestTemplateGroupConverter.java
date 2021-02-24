package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.RequestTemplateGroupDto;
import com.webank.taskman.dto.req.SaveRequestTemplateGroupReq;
import org.mapstruct.Mapper;
import org.mapstruct.ReportingPolicy;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface RequestTemplateGroupConverter extends BaseConverter<RequestTemplateGroupDto, RequestTemplateGroup> {

    RequestTemplateGroup saveReqToDomain(SaveRequestTemplateGroupReq req);

}