package com.webank.taskman.converter;

import java.util.ArrayList;
import java.util.List;

import org.springframework.beans.BeanUtils;
import org.springframework.stereotype.Service;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.RequestTemplateGroupDto;

@Service
public class RequestTemplateGroupConverter implements BaseConverter<RequestTemplateGroupDto, RequestTemplateGroup> {

    @Override
    public RequestTemplateGroup convertToEntity(RequestTemplateGroupDto dto) {
        RequestTemplateGroup entity = new RequestTemplateGroup();
        BeanUtils.copyProperties(dto, entity);
        return entity;
    }

    @Override
    public RequestTemplateGroupDto convertToDto(RequestTemplateGroup entity) {
        RequestTemplateGroupDto dto = new RequestTemplateGroupDto();
        BeanUtils.copyProperties(entity, dto);
        return dto;
    }

    @Override
    public List<RequestTemplateGroup> convertToEntities(List<RequestTemplateGroupDto> dtos) {
        if(dtos == null){
            return null;
        }
        
        List<RequestTemplateGroup> entities = new ArrayList<>();
        for(RequestTemplateGroupDto dto : dtos){
            RequestTemplateGroup entity = convertToEntity(dto);
            entities.add(entity);
        }
        return entities;
    }

    @Override
    public List<RequestTemplateGroupDto> convertToDtos(List<RequestTemplateGroup> entities) {
        if(entities == null){
            return null;
        }
        
        List<RequestTemplateGroupDto> dtos = new ArrayList<>();
        for(RequestTemplateGroup entity : entities){
            RequestTemplateGroupDto dto = convertToDto(entity);
            dtos.add(dto);
        }
        return dtos;
    }

}