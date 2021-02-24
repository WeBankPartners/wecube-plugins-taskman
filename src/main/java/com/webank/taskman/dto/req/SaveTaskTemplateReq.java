package com.webank.taskman.dto.req;

import java.util.List;

import com.webank.taskman.constant.TemplateTypeEnum;
import com.webank.taskman.dto.RoleDto;

public class SaveTaskTemplateReq {

    private String id;
    private String tempId;

    private String procDefId;

    private String procDefKey;

    private String procDefName;

    private String nodeDefId;

    private String nodeName;

    private String name;

    private String description;

    // @ApiModelProperty(value = "使用角色集",required = false,position = 9)
    private List<RoleDto> useRoles;

    // @ApiModelProperty(value = "管理角色集",required = false,position = 10)
    private List<RoleDto> manageRoles;

    // @ApiModelProperty(value = "任务表单模板", required = false, position = 11)
    private SaveFormTemplateReq form;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getTempId() {
        return tempId;
    }

    public SaveTaskTemplateReq setTempId(String tempId) {
        this.tempId = tempId;
        return this;
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

    public String getNodeDefId() {
        return nodeDefId;
    }

    public void setNodeDefId(String nodeDefId) {
        this.nodeDefId = nodeDefId;
    }

    public String getNodeName() {
        return nodeName;
    }

    public void setNodeName(String nodeName) {
        this.nodeName = nodeName;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public List<RoleDto> getUseRoles() {
        return useRoles;
    }

    public void setUseRoles(List<RoleDto> useRoles) {
        this.useRoles = useRoles;
    }

    public List<RoleDto> getManageRoles() {
        return manageRoles;
    }

    public void setManageRoles(List<RoleDto> manageRoles) {
        this.manageRoles = manageRoles;
    }

    public SaveFormTemplateReq getForm() {
        this.form.setTempType(TemplateTypeEnum.TASK.getType());
        return form;

    }

    public void setForm(SaveFormTemplateReq form) {
        this.form = form;
    }
}
