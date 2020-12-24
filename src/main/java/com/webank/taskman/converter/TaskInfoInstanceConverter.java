package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.dto.req.SaveTaskInfoReq;
import com.webank.taskman.dto.resp.TaskInfoInstanceResp;
import org.mapstruct.Mapper;
import org.mapstruct.ReportingPolicy;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface TaskInfoInstanceConverter extends BaseConverter<TaskInfoInstanceResp, TaskInfo> {

}