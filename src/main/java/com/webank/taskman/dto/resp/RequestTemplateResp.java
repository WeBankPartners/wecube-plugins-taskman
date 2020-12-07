package com.webank.taskman.dto.resp;

import com.webank.taskman.domain.RoleInfo;
import io.swagger.annotations.ApiModel;

import java.util.List;

@ApiModel
public class RequestTemplateResp {



    private String id;

    private String requestTempGroup;


    private String procDefKey;


    private String procDefId;


    private String procDefName;


    private String name;


    private String version;


    private String tags;


    private Integer status;

    private List<RoleInfo> roleIds;

    private List<RoleInfo> ManagementRole;


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

    public Integer getStatus() {
        return status;
    }

    public RequestTemplateResp setStatus(Integer status) {
        this.status = status;
        return this;
    }

    public List<RoleInfo> getRoleIds() {
        return roleIds;
    }

    public void setRoleIds(List<RoleInfo> roleIds) {
        this.roleIds = roleIds;
    }

    public List<RoleInfo> getManagementRole() {
        return ManagementRole;
    }

    public void setManagementRole(List<RoleInfo> managementRole) {
        ManagementRole = managementRole;
    }
}
