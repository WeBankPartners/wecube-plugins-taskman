package com.webank.taskman.dto.resp;

import com.webank.taskman.dto.RoleDTO;
import io.swagger.annotations.ApiModelProperty;

import java.util.List;

public class CreateRequestServiceMetaResp {

    @ApiModelProperty(value = "request-template-available-group",position = 1)
    List<RequestTemplateGroupResq> requestTemplateGroups;

    @ApiModelProperty(value = "form-item-template-currency",position = 2)
    List<FormItemTemplateResq>  currencyFormItems;

    @ApiModelProperty(value = "core-resources-role-all",position = 3)
    List<RootEntityResp> rootEntity;

    @ApiModelProperty(value = "core-resources-role-all",position = 4)
    List<RoleDTO> roles;

    public List<RequestTemplateGroupResq> getRequestTemplateGroups() {
        return requestTemplateGroups;
    }

    public void setRequestTemplateGroups(List<RequestTemplateGroupResq> requestTemplateGroups) {
        this.requestTemplateGroups = requestTemplateGroups;
    }

    public List<FormItemTemplateResq> getCurrencyFormItems() {
        return currencyFormItems;
    }

    public void setCurrencyFormItems(List<FormItemTemplateResq> currencyFormItems) {
        this.currencyFormItems = currencyFormItems;
    }

    public List<RootEntityResp> getRootEntity() {
        return rootEntity;
    }

    public void setRootEntity(List<RootEntityResp> rootEntity) {
        this.rootEntity = rootEntity;
    }

    public List<RoleDTO> getRoles() {
        return roles;
    }

    public void setRoles(List<RoleDTO> roles) {
        this.roles = roles;
    }

}
