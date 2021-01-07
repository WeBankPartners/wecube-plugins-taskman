package com.webank.taskman.dto.resp;

import com.webank.taskman.dto.RoleDTO;
import io.swagger.annotations.ApiModelProperty;

import java.util.ArrayList;
import java.util.List;

public class RequestTemplateResp  {

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

    private FormTemplateResp formTemplateResp;

    public String getId() {
        return id;
    }

    public RequestTemplateResp setId(String id) {
        this.id = id;
        return this;
    }

    public String getRequestTempGroup() {
        return requestTempGroup;
    }

    public RequestTemplateResp setRequestTempGroup(String requestTempGroup) {
        this.requestTempGroup = requestTempGroup;
        return this;
    }

    public String getRequestTempGroupName() {
        return requestTempGroupName;
    }

    public RequestTemplateResp setRequestTempGroupName(String requestTempGroupName) {
        this.requestTempGroupName = requestTempGroupName;
        return this;
    }

    public String getProcDefKey() {
        return procDefKey;
    }

    public RequestTemplateResp setProcDefKey(String procDefKey) {
        this.procDefKey = procDefKey;
        return this;
    }

    public String getProcDefId() {
        return procDefId;
    }

    public RequestTemplateResp setProcDefId(String procDefId) {
        this.procDefId = procDefId;
        return this;
    }

    public String getProcDefName() {
        return procDefName;
    }

    public RequestTemplateResp setProcDefName(String procDefName) {
        this.procDefName = procDefName;
        return this;
    }

    public String getName() {
        return name;
    }

    public RequestTemplateResp setName(String name) {
        this.name = name;
        return this;
    }

    public String getDescription() {
        return description;
    }

    public RequestTemplateResp setDescription(String description) {
        this.description = description;
        return this;
    }

    public String getVersion() {
        return version;
    }

    public RequestTemplateResp setVersion(String version) {
        this.version = version;
        return this;
    }

    public String getTags() {
        return tags;
    }

    public RequestTemplateResp setTags(String tags) {
        this.tags = tags;
        return this;
    }

    public String getStatus() {
        return status;
    }

    public RequestTemplateResp setStatus(String status) {
        this.status = status;
        return this;
    }

    public List<RoleDTO> getUseRoles() {
        return useRoles;
    }

    public RequestTemplateResp setUseRoles(List<RoleDTO> useRoles) {
        this.useRoles = useRoles;
        return this;
    }

    public List<RoleDTO> getManageRoles() {
        return manageRoles;
    }

    public RequestTemplateResp setManageRoles(List<RoleDTO> manageRoles) {
        this.manageRoles = manageRoles;
        return this;
    }

    public FormTemplateResp getFormTemplateResp() {
        return formTemplateResp;
    }

    public RequestTemplateResp setFormTemplateResp(FormTemplateResp formTemplateResp) {
        this.formTemplateResp = formTemplateResp;
        return this;
    }
}
