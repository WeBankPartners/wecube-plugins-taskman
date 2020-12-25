package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.FormInfo;
import com.webank.taskman.dto.resp.RequestFormResq;
import com.webank.taskman.dto.resp.TaskFormResq;
import org.mapstruct.Mapper;
import org.mapstruct.ReportingPolicy;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface TaskFormInfoConverter extends BaseConverter<TaskFormResq, FormInfo> {

}
