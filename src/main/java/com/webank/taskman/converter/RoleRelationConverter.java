package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.RoleRelation;
import com.webank.taskman.dto.RoleDTO;
import com.webank.taskman.support.core.dto.RolesDataDTO;
import com.webank.taskman.support.core.dto.RolesDataResponse;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.Mappings;
import org.mapstruct.ReportingPolicy;

import java.util.List;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface RoleRelationConverter extends BaseConverter<RoleDTO, RoleRelation> {


    @Mappings({
            @Mapping(target = "displayName",source ="description" ),
    })
    RoleDTO rolesDataResponseToDto(RolesDataResponse rolesDataResponse);

    List<RoleDTO> rolesDataResponseToDtoList(List<RolesDataResponse> rolesDataResponseList);

    @Mappings({
            @Mapping(target = "description",source ="displayName" ),
            @Mapping(target = "roleName",source ="name" ),
            @Mapping(target = "roleId",source ="id" ),
    })
    RolesDataResponse  roleDTOToRolesDataResponse(RolesDataDTO rolesDataDTO);

    List<RolesDataResponse> roleDTOToRolesDataResponseList(List<RolesDataDTO> rolesDataDTOList);
}