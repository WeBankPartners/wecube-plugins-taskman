package com.webank.taskman.converter;

import java.util.ArrayList;
import java.util.List;

import org.springframework.beans.BeanUtils;
import org.springframework.stereotype.Service;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.RoleRelation;
import com.webank.taskman.dto.RoleDto;
import com.webank.taskman.support.platform.dto.RolesDataResponse;

@Service
public class RoleRelationConverter implements BaseConverter<RoleDto, RoleRelation> {


//    @Mappings({
//            @Mapping(target = "roleName",source ="name" ),
//    })
    public RoleDto rolesDataResponseToDto(RolesDataResponse rolesDataResponse){
        RoleDto dto = new RoleDto();
        dto.setRoleName(rolesDataResponse.getName());
        dto.setDisplayName(rolesDataResponse.getDisplayName());
        
        return dto;
    }

    public List<RoleDto> rolesDataResponseToDtoList(List<RolesDataResponse> rolesDataResponseList){
        List<RoleDto> dtos = new ArrayList<>();
        for(RolesDataResponse entity : rolesDataResponseList){
            RoleDto dto = new RoleDto();
            dto.setRoleName(entity.getName());
            dto.setDisplayName(entity.getDisplayName());
            
            dtos.add(dto);
        }
        
        return dtos;
    }

    @Override
    public RoleRelation convertToEntity(RoleDto dto) {
        RoleRelation entity = new RoleRelation();
        BeanUtils.copyProperties(dto, entity);
        return entity;
    }

    @Override
    public RoleDto convertToDto(RoleRelation entity) {
        RoleDto dto = new RoleDto();
        BeanUtils.copyProperties(entity, dto);
        return dto;
    }

    @Override
    public List<RoleRelation> convertToEntities(List<RoleDto> dtos) {
        if(dtos == null){
            return null;
        }
        
        List<RoleRelation> entities = new ArrayList<>();
        for(RoleDto dto : dtos){
            RoleRelation entity = convertToEntity(dto);
            entities.add(entity);
        }
        return entities;
    }

    @Override
    public List<RoleDto> convertToDtos(List<RoleRelation> entities) {
        if(entities == null){
            return null;
        }
        
        List<RoleDto> dtos = new ArrayList<>();
        for(RoleRelation entity : entities){
            RoleDto dto = convertToDto(entity);
            dtos.add(dto);
        }
        return dtos;
    }

}