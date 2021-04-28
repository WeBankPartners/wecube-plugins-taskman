package com.webank.taskman.converter;

import java.util.ArrayList;
import java.util.List;

import org.springframework.beans.BeanUtils;
import org.springframework.stereotype.Service;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.TaskTemplate;
import com.webank.taskman.dto.req.TaskTemplateSaveReqDto;
import com.webank.taskman.dto.req.TemplateQueryReqDto;
import com.webank.taskman.dto.resp.TaskTemplateByRoleRespDto;
import com.webank.taskman.dto.resp.TaskTemplateRespDto;

@Service
public class TaskTemplateConverter implements BaseConverter<TaskTemplateRespDto, TaskTemplate> {

    public TaskTemplate convertToTaskTemplate(TaskTemplateSaveReqDto dto){
        if(dto == null){
            return null;
        }
        
        TaskTemplate entity = new TaskTemplate();
        //better to apply getter/setter instead,avoid to introduce bean utility property copy
        BeanUtils.copyProperties(dto, entity);
        return entity;
    }

    public TaskTemplate convertToTaskTemplate(TemplateQueryReqDto dto){
        if(dto == null){
            return null;
        }
        
        TaskTemplate entity = new TaskTemplate();
        BeanUtils.copyProperties(dto, entity);
        return entity;
    }

    public List<TaskTemplateByRoleRespDto> convertToTaskTemplateByRoleRespDtos(List<TaskTemplate> entities){
        if(entities == null){
            return null;
        }
        
        List<TaskTemplateByRoleRespDto> dtos = new ArrayList<>();
        for(TaskTemplate entity : entities){
            TaskTemplateByRoleRespDto dto = new TaskTemplateByRoleRespDto();
            BeanUtils.copyProperties(entity, dto);
            
            dtos.add(dto);
        }
        return dtos;
    }

    @Override
    public TaskTemplate convertToEntity(TaskTemplateRespDto dto) {
        if(dto == null){
            return null;
        }
        
        TaskTemplate entity = new TaskTemplate();
        BeanUtils.copyProperties(dto, entity);
        return entity;
    }

    @Override
    public TaskTemplateRespDto convertToDto(TaskTemplate entity) {
        if(entity == null){
            return null;
        }
        TaskTemplateRespDto dto = new TaskTemplateRespDto();
        
        BeanUtils.copyProperties(entity, dto);
        return dto;
    }

    @Override
    public List<TaskTemplate> convertToEntities(List<TaskTemplateRespDto> dtos) {
        if(dtos == null){
            return null;
        }
        
        List<TaskTemplate> entities = new ArrayList<>();
        for(TaskTemplateRespDto dto : dtos){
            entities.add(convertToEntity(dto));
        }
        return entities;
    }

    @Override
    public List<TaskTemplateRespDto> convertToDtos(List<TaskTemplate> entities) {
        if(entities == null){
            return null;
        }
        
        List<TaskTemplateRespDto> dtos = new ArrayList<>();
        for(TaskTemplate entity : entities){
            dtos.add(convertToDto(entity));
        }
        return dtos;
    }
}