package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.req.AddTemplateGropReq;
import com.webank.taskman.dto.TemplateGroupCreateVO;
import com.webank.taskman.dto.TemplateGroupDTO;
import com.webank.taskman.dto.TemplateGroupVO;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.ReportingPolicy;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface TemplateGroupConverter extends BaseConverter<TemplateGroupDTO, RequestTemplateGroup> {

    @Mapping(target = "version", ignore = true)
    RequestTemplateGroup voToDomain(TemplateGroupVO vo);

    RequestTemplateGroup cVoTODomain(TemplateGroupCreateVO vo);

    RequestTemplateGroup addReqDomain(AddTemplateGropReq req);
}