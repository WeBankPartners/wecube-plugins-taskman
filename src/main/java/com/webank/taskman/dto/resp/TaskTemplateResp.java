package com.webank.taskman.dto.resp;


import com.webank.taskman.dto.RoleDTO;
import io.swagger.annotations.ApiModel;

import java.util.ArrayList;
import java.util.List;

@ApiModel
public class TaskTemplateResp {
    private String id;

    private String procDefId;

    private String procDefKey;

    private String procDefName;

    private String name;

    private String procNode;

    private String description;

    private List<RoleDTO> useRoles = new ArrayList<>();

    private List<RoleDTO> manageRoles = new ArrayList<>();

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getProcDefId() {
        return procDefId;
    }

    public void setProcDefId(String procDefId) {
        this.procDefId = procDefId;
    }

    public String getProcDefKey() {
        return procDefKey;
    }

    public void setProcDefKey(String procDefKey) {
        this.procDefKey = procDefKey;
    }

    public String getProcDefName() {
        return procDefName;
    }

    public void setProcDefName(String procDefName) {
        this.procDefName = procDefName;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getProcNode() {
        return procNode;
    }

    public void setProcNode(String procNode) {
        this.procNode = procNode;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
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
        return "TaskTemplateSVResp{" +
                "id='" + id + '\'' +
                ", procDefId='" + procDefId + '\'' +
                ", procDefKey='" + procDefKey + '\'' +
                ", procDefName='" + procDefName + '\'' +
                ", name='" + name + '\'' +
                ", procNode='" + procNode + '\'' +
                ", description='" + description + '\'' +
                ", useRoles=" + useRoles +
                ", manageRoles=" + manageRoles +
                '}';
    }
}
