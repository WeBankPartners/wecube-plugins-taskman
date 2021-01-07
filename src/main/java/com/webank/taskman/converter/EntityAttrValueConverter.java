package com.webank.taskman.converter;


import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.dto.CreateTaskDto;
import com.webank.taskman.support.core.dto.DynamicEntityValueDto;
import org.mapstruct.Mapper;
import org.mapstruct.ReportingPolicy;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface EntityAttrValueConverter extends BaseConverter<CreateTaskDto.EntityAttrValueDto, DynamicEntityValueDto.DynamicEntityAttrValueDto> {
}
