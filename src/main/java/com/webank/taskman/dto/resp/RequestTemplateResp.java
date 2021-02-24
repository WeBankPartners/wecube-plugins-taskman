package com.webank.taskman.dto.resp;

import java.util.ArrayList;
import java.util.List;

import com.webank.taskman.dto.RoleDTO;

public class RequestTemplateResp  {

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
    private List<RoleDTO> useRoles = new ArrayList<>();
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
