package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.TemplateGroup;
import com.webank.taskman.dto.TemplateGroupDTO;
import com.webank.taskman.dto.TemplateGroupVO;
import org.mapstruct.Mapper;
import org.mapstruct.ReportingPolicy;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface TemplateGroupConverter extends BaseConverter<TemplateGroupDTO, TemplateGroup> {

    TemplateGroup voToDomain(TemplateGroupVO vo);
}