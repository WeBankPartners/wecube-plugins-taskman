package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.dto.req.SaveTaskInfoReq;
import com.webank.taskman.dto.resp.TaskInfoGetResp;
import org.mapstruct.Mapper;
import org.mapstruct.ReportingPolicy;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface TaskInfoGetConverter extends BaseConverter<TaskInfoGetResp, TaskInfo> {

}