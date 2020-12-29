package com.webank.taskman.dto;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import java.util.ArrayList;
import java.util.List;

@ApiModel
public class RequestTemplateDTO {



    @ApiModelProperty(value = "",position = 100)
    private String id;

    @ApiModelProperty(value = "",position = 101)
    private String requestTempGroup;

    private String requestTempGroupName;

    @ApiModelProperty(value = "",position = 102)
    private String procDefKey;


    @ApiModelProperty(value = "",position = 103)
    private String procDefId;


    @ApiModelProperty(value = "",position = 104)
    private String procDefName;


    @ApiModelProperty(value = "",position = 105)
    private String name;


    @ApiModelProperty(value = "",position = 106)
    private String description;
    @ApiModelProperty(value = "",position = 106)
    private String version;


    @ApiModelProperty(value = "",position = 107)
    private String tags;


    @ApiModelProperty(value = "",position = 108)
    private String status;

    @ApiModelProperty(value = "",position = 109)
    private List<RoleDTO> useRoles = new ArrayList<>();

    @ApiModelProperty(value = "",position = 100)
    private List<RoleDTO> manageRoles = new ArrayList<>();


    public String getId() {
        return id;
    }

    public RequestTemplateDTO setId(String id) {
        this.id = id;
        return this;
    }

    public String getRequestTempGroup() {
        return requestTempGroup;
    }

    public RequestTemplateDTO setRequestTempGroup(String requestTempGroup) {
        this.requestTempGroup = requestTempGroup;
        return this;
    }

    public String getRequestTempGroupName() {
        return requestTempGroupName;
    }

    public void setRequestTempGroupName(String requestTempGroupName) {
        this.requestTempGroupName = requestTempGroupName;
    }

    public String getProcDefKey() {
        return procDefKey;
    }

    public RequestTemplateDTO setProcDefKey(String procDefKey) {
        this.procDefKey = procDefKey;
        return this;
    }

    public String getProcDefId() {
        return procDefId;
    }

    public RequestTemplateDTO setProcDefId(String procDefId) {
        this.procDefId = procDefId;
        return this;
    }

    public String getProcDefName() {
        return procDefName;
    }

    public RequestTemplateDTO setProcDefName(String procDefName) {
        this.procDefName = procDefName;
        return this;
    }

    public String getName() {
        return name;
    }

    public RequestTemplateDTO setName(String name) {
        this.name = name;
        return this;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getVersion() {
        return version;
    }

    public RequestTemplateDTO setVersion(String version) {
        this.version = version;
        return this;
    }

    public String getTags() {
        return tags;
    }

    public RequestTemplateDTO setTags(String tags) {
        this.tags = tags;
        return this;
    }

    public String getStatus() {
        return status;
    }

    public RequestTemplateDTO setStatus(String status) {
        this.status = status;
        return this;
    }

    public List<RoleDTO> getUseRoles() {
        return useRoles;
    }

    public void setUseRoles(List<RoleDTO> useRoles) {
        this.useRoles = useRoles;
    }

    public List<RoleDTO> getManageRoles() {
        return manageRoles;
    }

    public void setManageRoles(List<RoleDTO> manageRoles) {
        this.manageRoles = manageRoles;
    }

    @Override
    public String toString() {
        return "RequestTemplateDTO{" +
                "id='" + id + '\'' +
                ", requestTempGroup='" + requestTempGroup + '\'' +
                ", requestTempGroupName='" + requestTempGroupName + '\'' +
                ", procDefKey='" + procDefKey + '\'' +
                ", procDefId='" + procDefId + '\'' +
                ", procDefName='" + procDefName + '\'' +
                ", name='" + name + '\'' +
                ", description='" + description + '\'' +
                ", version='" + version + '\'' +
                ", tags='" + tags + '\'' +
                ", status='" + status + '\'' +
                ", useRoles=" + useRoles +
                ", manageRoles=" + manageRoles +
                '}';
    }
}
