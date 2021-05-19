package com.webank.taskman.converter;

import java.util.ArrayList;
import java.util.List;

import org.springframework.beans.BeanUtils;
import org.springframework.stereotype.Service;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.dto.TaskInfoDto;
import com.webank.taskman.dto.req.TaskInfoQueryReqDto;
import com.webank.taskman.dto.resp.TaskInfoRespDto;

@Service
public class TaskInfoConverter implements BaseConverter<TaskInfoDto, TaskInfo> {


//    @Mappings({
//            @Mapping(target = "nodeName",source = "taskName"),
//            @Mapping(target = "description",source = "taskDescription"),
//            @Mapping(target = "nodeDefId",source = "taskNodeId"),
//    })
//    TaskInfo toEntityByReq(TaskInfoReq req);

    public TaskInfo convertToTaskInfoByQuery(TaskInfoQueryReqDto reqDto){
        TaskInfo info = new TaskInfo();
        BeanUtils.copyProperties(reqDto, info);
        return info;
        
    }

//    TaskInfoInstanceRespDto toInstanceResp(TaskInfo taskInfo);

    public TaskInfoRespDto convertToTaskInfoRespDto(TaskInfo taskInfo){
        TaskInfoRespDto dto = new TaskInfoRespDto();
        BeanUtils.copyProperties(taskInfo, dto);
        return dto;
    }

    @Override
    public TaskInfo convertToEntity(TaskInfoDto dto) {
        TaskInfo entity = new TaskInfo();
        BeanUtils.copyProperties(dto, entity);
        
        return entity;
    }

    @Override
    public TaskInfoDto convertToDto(TaskInfo entity) {
        TaskInfoDto dto = new TaskInfoDto();
        BeanUtils.copyProperties(entity, dto);
        return dto;
    }

    @Override
    public List<TaskInfo> convertToEntities(List<TaskInfoDto> dtos) {
        if(dtos == null){
            return null;
            
        }
        
        List<TaskInfo> entities = new ArrayList<>();
        for(TaskInfoDto dto : dtos){
            TaskInfo entity = convertToEntity(dto);
            entities.add(entity);
        }
        return entities;
    }

    @Override
    public List<TaskInfoDto> convertToDtos(List<TaskInfo> entities) {
        if(entities == null){
            return null;
        }
        
        List<TaskInfoDto> dtos = new ArrayList<>();
        for(TaskInfo entity : entities){
            TaskInfoDto dto = convertToDto(entity);
            dtos.add(dto);
        }
        return dtos;
    }


}