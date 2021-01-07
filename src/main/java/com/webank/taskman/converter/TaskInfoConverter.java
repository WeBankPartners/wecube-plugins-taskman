package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.dto.CoreCreateTaskDTO.TaskInfoReq;
import com.webank.taskman.dto.TaskInfoDTO;
import com.webank.taskman.dto.req.QueryTaskInfoReq;
import com.webank.taskman.dto.resp.TaskInfoResp;
import com.webank.taskman.dto.resp.TaskInfoInstanceResp;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.Mappings;
import org.mapstruct.ReportingPolicy;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface TaskInfoConverter extends BaseConverter<TaskInfoDTO, TaskInfo> {


    @Mappings({
            @Mapping(target = "nodeName",source = "taskName"),
            @Mapping(target = "description",source = "taskDescription"),
            @Mapping(target = "nodeDefId",source = "taskNodeId"),
    })
    TaskInfo toEntityByReq(TaskInfoReq req);

    TaskInfo toEntityByQuery(QueryTaskInfoReq req);

    TaskInfoInstanceResp toInstanceResp(TaskInfo taskInfo);

    TaskInfoResp toResp(TaskInfo taskInfo);


}