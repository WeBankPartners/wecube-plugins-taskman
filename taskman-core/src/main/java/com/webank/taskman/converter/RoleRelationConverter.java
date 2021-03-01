package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.RoleRelation;
import com.webank.taskman.dto.RoleDto;
import com.webank.taskman.support.core.dto.RolesDataResponse;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.Mappings;
import org.mapstruct.ReportingPolicy;

import java.util.List;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface RoleRelationConverter extends BaseConverter<RoleDto, RoleRelation> {


    @Mappings({
            @Mapping(target = "roleName",source ="name" ),
    })
    RoleDto rolesDataResponseToDto(RolesDataResponse rolesDataResponse);

    List<RoleDto> rolesDataResponseToDtoList(List<RolesDataResponse> rolesDataResponseList);

}