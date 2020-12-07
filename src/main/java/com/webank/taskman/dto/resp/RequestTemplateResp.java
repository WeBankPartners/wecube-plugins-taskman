package com.webank.taskman.dto.resp;

import com.webank.taskman.domain.RoleInfo;
import com.webank.taskman.dto.RoleDTO;
import io.swagger.annotations.ApiModel;

import java.util.ArrayList;
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

    private List<RoleDTO> useRoles = new ArrayList<>();

    private List<RoleDTO> manageRoles = new ArrayList<>();


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
}
