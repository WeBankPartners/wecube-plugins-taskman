package com.webank.taskman.dto.resp;

import com.webank.taskman.dto.RoleDTO;
import io.swagger.annotations.ApiModelProperty;

import java.util.List;

public class CreateRequestServiceMetaResp {

    @ApiModelProperty(value = "request-template-available-group",position = 1)
    List<RequestTemplateGroupResq> requestTemplateGroups;

    @ApiModelProperty(value = "core-resources-role-all",position = 1)
    List<RoleDTO> roles;


}
