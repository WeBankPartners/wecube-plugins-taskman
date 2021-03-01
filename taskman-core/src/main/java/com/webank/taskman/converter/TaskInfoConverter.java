package com.webank.taskman.converter;

import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.Mappings;
import org.mapstruct.ReportingPolicy;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.dto.CoreCreateTaskDto.TaskInfoReq;
import com.webank.taskman.dto.TaskInfoDto;
import com.webank.taskman.dto.req.TaskInfoQueryReqDto;
import com.webank.taskman.dto.resp.TaskInfoInstanceRespDto;
import com.webank.taskman.dto.resp.TaskInfoRespDto;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface TaskInfoConverter extends BaseConverter<TaskInfoDto, TaskInfo> {


    @Mappings({
            @Mapping(target = "nodeName",source = "taskName"),
            @Mapping(target = "description",source = "taskDescription"),
            @Mapping(target = "nodeDefId",source = "taskNodeId"),
    })
    TaskInfo toEntityByReq(TaskInfoReq req);

    TaskInfo toEntityByQuery(TaskInfoQueryReqDto req);

    TaskInfoInstanceRespDto toInstanceResp(TaskInfo taskInfo);

    TaskInfoRespDto toResp(TaskInfo taskInfo);


}