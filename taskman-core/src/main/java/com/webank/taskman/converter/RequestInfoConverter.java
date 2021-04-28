package com.webank.taskman.converter;

import java.util.ArrayList;
import java.util.List;

import org.springframework.beans.BeanUtils;
import org.springframework.stereotype.Service;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.RequestInfo;
import com.webank.taskman.dto.CreateTaskDto;
import com.webank.taskman.dto.RequestInfoDto;
import com.webank.taskman.dto.req.RequestInfoSaveReqDto;
import com.webank.taskman.dto.resp.RequestInfoInstanceResqDto;
import com.webank.taskman.dto.resp.RequestInfoQueryResultDto;

@Service
public class RequestInfoConverter implements BaseConverter<RequestInfoDto, RequestInfo> {

    public RequestInfo reqToDomain(RequestInfoSaveReqDto dto){
        RequestInfo entity = new RequestInfo();
        BeanUtils.copyProperties(dto, entity);
        return entity;
    }

    public RequestInfo createDtoToDomain(CreateTaskDto dto){
        RequestInfo entity = new RequestInfo();
        BeanUtils.copyProperties(dto, entity);
        return entity;
    }

    public RequestInfoQueryResultDto toResp(RequestInfo requestInfo){
        RequestInfoQueryResultDto dto = new RequestInfoQueryResultDto();
        BeanUtils.copyProperties(requestInfo, dto);
        return dto;
    }

    public RequestInfoInstanceResqDto toInstanceResp(RequestInfo requestInfo){
        RequestInfoInstanceResqDto dto = new RequestInfoInstanceResqDto();
        BeanUtils.copyProperties(requestInfo, dto);
        return dto;
    }

    @Override
    public RequestInfo convertToEntity(RequestInfoDto dto) {
        RequestInfo entity = new RequestInfo();
        BeanUtils.copyProperties(dto, entity);
        return entity;
    }

    @Override
    public RequestInfoDto convertToDto(RequestInfo entity) {
        RequestInfoDto dto = new RequestInfoDto();
        BeanUtils.copyProperties(entity, dto);
        return dto;
    }

    @Override
    public List<RequestInfo> convertToEntities(List<RequestInfoDto> dtos) {
        if(dtos == null){
            return null;
        }
        
        List<RequestInfo> entities = new ArrayList<>();
        for(RequestInfoDto dto : dtos){
            RequestInfo entity = convertToEntity(dto);
            entities.add(entity);
        }
        return entities;
    }

    @Override
    public List<RequestInfoDto> convertToDtos(List<RequestInfo> entities) {
        if(entities == null){
            return null;
        }
        
        List<RequestInfoDto> dtos = new ArrayList<>();
        for(RequestInfo entity : entities){
            RequestInfoDto dto = convertToDto(entity);
            dtos.add(dto);
        }
        return dtos;
    }
}
