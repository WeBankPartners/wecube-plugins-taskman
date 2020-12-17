package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.TemplateGroupDTO;
import com.webank.taskman.dto.TemplateGroupVO;
import com.webank.taskman.dto.req.SaveRequestTemplateGropReq;
import org.mapstruct.Mapper;
import org.mapstruct.ReportingPolicy;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface RequestTemplateGroupConverter extends BaseConverter<TemplateGroupDTO, RequestTemplateGroup> {

    RequestTemplateGroup voToDomain(TemplateGroupVO vo);


    RequestTemplateGroup addOrUpdateDomain(SaveRequestTemplateGropReq req);

}