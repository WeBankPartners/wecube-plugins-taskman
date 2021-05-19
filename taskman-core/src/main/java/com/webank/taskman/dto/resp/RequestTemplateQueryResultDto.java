package com.webank.taskman.dto.resp;

import java.util.ArrayList;
import java.util.List;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.webank.taskman.dto.RoleDto;

@JsonInclude(JsonInclude.Include.NON_NULL)
public class RequestTemplateQueryResultDto  {

    private String id;
    private String requestTempGroup;
    private String requestTempGroupName;
    private String procDefKey;
    private String procDefId;
    private String procDefName;
    private String name;
    private String description;
    private String version;
    private String tags;
    private String status;
    private List<RoleDto> useRoles = new ArrayList<>();
    private List<RoleDto> manageRoles = new ArrayList<>();

    private FormTemplateQueryResultDto formTemplateResp;

    public String getId() {
        return id;
    }

    public RequestTemplateQueryResultDto setId(String id) {
        this.id = id;
        return this;
    }

    public String getRequestTempGroup() {
        return requestTempGroup;
    }

    public RequestTemplateQueryResultDto setRequestTempGroup(String requestTempGroup) {
        this.requestTempGroup = requestTempGroup;
        return this;
    }

    public String getRequestTempGroupName() {
        return requestTempGroupName;
    }

    public RequestTemplateQueryResultDto setRequestTempGroupName(String requestTempGroupName) {
        this.requestTempGroupName = requestTempGroupName;
        return this;
    }

    public String getProcDefKey() {
        return procDefKey;
    }

    public RequestTemplateQueryResultDto setProcDefKey(String procDefKey) {
        this.procDefKey = procDefKey;
        return this;
    }

    public String getProcDefId() {
        return procDefId;
    }

    public RequestTemplateQueryResultDto setProcDefId(String procDefId) {
        this.procDefId = procDefId;
        return this;
    }

    public String getProcDefName() {
        return procDefName;
    }

    public RequestTemplateQueryResultDto setProcDefName(String procDefName) {
        this.procDefName = procDefName;
        return this;
    }

    public String getName() {
        return name;
    }

    public RequestTemplateQueryResultDto setName(String name) {
        this.name = name;
        return this;
    }

    public String getDescription() {
        return description;
    }

    public RequestTemplateQueryResultDto setDescription(String description) {
        this.description = description;
        return this;
    }

    public String getVersion() {
        return version;
    }

    public RequestTemplateQueryResultDto setVersion(String version) {
        this.version = version;
        return this;
    }

    public String getTags() {
        return tags;
    }

    public RequestTemplateQueryResultDto setTags(String tags) {
        this.tags = tags;
        return this;
    }

    public String getStatus() {
        return status;
    }

    public RequestTemplateQueryResultDto setStatus(String status) {
        this.status = status;
        return this;
    }

    public List<RoleDto> getUseRoles() {
        return useRoles;
    }

    public RequestTemplateQueryResultDto setUseRoles(List<RoleDto> useRoles) {
        this.useRoles = useRoles;
        return this;
    }

    public List<RoleDto> getManageRoles() {
        return manageRoles;
    }

    public RequestTemplateQueryResultDto setManageRoles(List<RoleDto> manageRoles) {
        this.manageRoles = manageRoles;
        return this;
    }

    public FormTemplateQueryResultDto getFormTemplateResp() {
        return formTemplateResp;
    }

    public RequestTemplateQueryResultDto setFormTemplateResp(FormTemplateQueryResultDto formTemplateResp) {
        this.formTemplateResp = formTemplateResp;
        return this;
    }
}
