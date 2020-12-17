package com.webank.taskman.dto.req;

import com.webank.taskman.dto.RoleDTO;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import javax.validation.constraints.NotBlank;
import java.util.List;

@ApiModel(value = "AddRequestTemplateReq",description = "add RequestTemplate req")
public class SaveRequestTemplateReq {

    @ApiModelProperty(value = "主键",required = false,dataType = "String",position = 100)
    private String id;

    @NotBlank(message = "模板组编号不能为空")
    @ApiModelProperty(value = "模板组编号",required = true,dataType = "String",position = 101)
    private String requestTempGroup;

    @NotBlank(message = "流程编排id不能为空")
    @ApiModelProperty(value = "流程编排Id",required = true,dataType = "String",position = 102)
    private String procDefId;

    @NotBlank(message = "流程编排key不能为空")
    @ApiModelProperty(value = "流程编排key",required = true,dataType = "String",position = 103)
    private String procDefKey;

    @NotBlank(message = "流程编排名称不能为空")
    @ApiModelProperty(value = "流程编排名称",required = true,dataType = "String",position = 104)
    private String procDefName;

    @NotBlank(message = "名称不能为空")
    @ApiModelProperty(value = "请求模板名称",required = true,dataType = "String",position = 105)
    private String name;
    @ApiModelProperty(value = "描述",required = true,dataType = "String",position = 106)
    private String description;

    @ApiModelProperty(value = "标签",required = true,dataType = "String",position = 107)
    private String tags;


    @ApiModelProperty(value = "使用角色集",required = false,dataType = "String",position = 108)
    private List<RoleDTO> useRoles;

    @ApiModelProperty(value = "管理角色集",required = false,dataType = "String",position = 109)
    private List<RoleDTO> manageRoles;


    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getRequestTempGroup() {
        return requestTempGroup;
    }

    public void setRequestTempGroup(String requestTempGroup) {
        this.requestTempGroup = requestTempGroup;
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

    public String getTags() {
        return tags;
    }

    public void setTags(String tags) {
        this.tags = tags;
    }





}
