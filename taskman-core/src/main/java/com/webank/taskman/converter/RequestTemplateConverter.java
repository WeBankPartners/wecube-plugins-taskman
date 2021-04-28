package com.webank.taskman.converter;

import java.util.ArrayList;
import java.util.List;

import org.springframework.beans.BeanUtils;
import org.springframework.stereotype.Service;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.RequestTemplateDto;
import com.webank.taskman.dto.req.RequestTemplateSaveReqDto;
import com.webank.taskman.dto.resp.RequestTemplateRespDto;

@Service
public class RequestTemplateConverter implements BaseConverter<RequestTemplateDto, RequestTemplate> {

    @Override
    public RequestTemplate convertToEntity(RequestTemplateDto dto) {
        RequestTemplate entity = new RequestTemplate();
        BeanUtils.copyProperties(dto, entity);
        return entity;
    }

    @Override
    public RequestTemplateDto convertToDto(RequestTemplate entity) {
        RequestTemplateDto dto = new RequestTemplateDto();
        BeanUtils.copyProperties(entity, dto);
        return dto;
    }

    @Override
    public List<RequestTemplate> convertToEntities(List<RequestTemplateDto> dtos) {
        if(dtos == null){
            return null;
        }
        
        List<RequestTemplate> entities = new ArrayList<>();
        for(RequestTemplateDto dto : dtos){
            RequestTemplate entity = convertToEntity(dto);
            entities.add(entity);
        }
        return entities;
    }

    @Override
    public List<RequestTemplateDto> convertToDtos(List<RequestTemplate> entities) {
        if(entities == null){
            return null;
        }
        
        List<RequestTemplateDto> dtos = new ArrayList<>();
        for(RequestTemplate entity : entities){
            RequestTemplateDto dto = convertToDto(entity);
            dtos.add(dto);
        }
        return dtos;
    }

    public RequestTemplate saveReqToEntity(RequestTemplateSaveReqDto dto){
        RequestTemplate entity = new RequestTemplate();
        BeanUtils.copyProperties(dto, entity);
        return entity;
        
    }

    public RequestTemplateRespDto toRespByEntity(RequestTemplate entity){
        RequestTemplateRespDto dto = new RequestTemplateRespDto();
        BeanUtils.copyProperties(entity, dto);
        return dto;
    }

}