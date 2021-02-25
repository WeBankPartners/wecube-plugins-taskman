package com.webank.taskman.dto.resp;

import java.util.ArrayList;
import java.util.List;

import com.webank.taskman.dto.RoleDto;

public class RequestTemplateRespDto  {

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

    private FormTemplateRespDto formTemplateResp;

    public String getId() {
        return id;
    }

    public RequestTemplateRespDto setId(String id) {
        this.id = id;
        return this;
    }

    public String getRequestTempGroup() {
        return requestTempGroup;
    }

    public RequestTemplateRespDto setRequestTempGroup(String requestTempGroup) {
        this.requestTempGroup = requestTempGroup;
        return this;
    }

    public String getRequestTempGroupName() {
        return requestTempGroupName;
    }

    public RequestTemplateRespDto setRequestTempGroupName(String requestTempGroupName) {
        this.requestTempGroupName = requestTempGroupName;
        return this;
    }

    public String getProcDefKey() {
        return procDefKey;
    }

    public RequestTemplateRespDto setProcDefKey(String procDefKey) {
        this.procDefKey = procDefKey;
        return this;
    }

    public String getProcDefId() {
        return procDefId;
    }

    public RequestTemplateRespDto setProcDefId(String procDefId) {
        this.procDefId = procDefId;
        return this;
    }

    public String getProcDefName() {
        return procDefName;
    }

    public RequestTemplateRespDto setProcDefName(String procDefName) {
        this.procDefName = procDefName;
        return this;
    }

    public String getName() {
        return name;
    }

    public RequestTemplateRespDto setName(String name) {
        this.name = name;
        return this;
    }

    public String getDescription() {
        return description;
    }

    public RequestTemplateRespDto setDescription(String description) {
        this.description = description;
        return this;
    }

    public String getVersion() {
        return version;
    }

    public RequestTemplateRespDto setVersion(String version) {
        this.version = version;
        return this;
    }

    public String getTags() {
        return tags;
    }

    public RequestTemplateRespDto setTags(String tags) {
        this.tags = tags;
        return this;
    }

    public String getStatus() {
        return status;
    }

    public RequestTemplateRespDto setStatus(String status) {
        this.status = status;
        return this;
    }

    public List<RoleDto> getUseRoles() {
        return useRoles;
    }

    public RequestTemplateRespDto setUseRoles(List<RoleDto> useRoles) {
        this.useRoles = useRoles;
        return this;
    }

    public List<RoleDto> getManageRoles() {
        return manageRoles;
    }

    public RequestTemplateRespDto setManageRoles(List<RoleDto> manageRoles) {
        this.manageRoles = manageRoles;
        return this;
    }

    public FormTemplateRespDto getFormTemplateResp() {
        return formTemplateResp;
    }

    public RequestTemplateRespDto setFormTemplateResp(FormTemplateRespDto formTemplateResp) {
        this.formTemplateResp = formTemplateResp;
        return this;
    }
}
