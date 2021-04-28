package com.webank.taskman.converter;

import java.util.List;

import org.mapstruct.Mapping;
import org.mapstruct.Mappings;
import org.springframework.stereotype.Service;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.dto.TaskInfoDto;
import com.webank.taskman.dto.platform.PlatformTaskCreationReqDto.TaskInfoReq;
import com.webank.taskman.dto.req.TaskInfoQueryReqDto;
import com.webank.taskman.dto.resp.TaskInfoInstanceRespDto;
import com.webank.taskman.dto.resp.TaskInfoRespDto;

@Service
public class TaskInfoConverter implements BaseConverter<TaskInfoDto, TaskInfo> {


    @Mappings({
            @Mapping(target = "nodeName",source = "taskName"),
            @Mapping(target = "description",source = "taskDescription"),
            @Mapping(target = "nodeDefId",source = "taskNodeId"),
    })
    TaskInfo toEntityByReq(TaskInfoReq req);

    TaskInfo toEntityByQuery(TaskInfoQueryReqDto req);

    TaskInfoInstanceRespDto toInstanceResp(TaskInfo taskInfo);

    TaskInfoRespDto toResp(TaskInfo taskInfo);

    @Override
    public TaskInfo convertToEntity(TaskInfoDto dto) {
        // TODO Auto-generated method stub
        return null;
    }

    @Override
    public TaskInfoDto convertToDto(TaskInfo entity) {
        // TODO Auto-generated method stub
        return null;
    }

    @Override
    public List<TaskInfo> convertToEntities(List<TaskInfoDto> dtos) {
        // TODO Auto-generated method stub
        return null;
    }

    @Override
    public List<TaskInfoDto> convertToDtos(List<TaskInfo> entities) {
        // TODO Auto-generated method stub
        return null;
    }


}