package com.webank.taskman.dto.req;

import com.webank.taskman.dto.RoleDTO;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import javax.validation.constraints.NotBlank;
import java.util.List;

@ApiModel(value = "AddRequestTemplateReq",description = "add RequestTemplate req")
public class QueryRequestTemplateReq {

    @ApiModelProperty(value = "主键",required = false,dataType = "String",position = 100)
    private String id;

    @ApiModelProperty(value = "模板组编号",required = false,dataType = "String",position = 101)
    private String requestTempGroup;

    @ApiModelProperty(value = "流程编排Id",required = false,dataType = "String",position = 102)
    private String procDefId;

    @ApiModelProperty(value = "流程编排key",required = false,dataType = "String",position = 103)
    private String procDefKey;

    @ApiModelProperty(value = "流程编排名称",required = false,dataType = "String",position = 104)
    private String procDefName;

    @ApiModelProperty(value = "请求模板名称",required = false,dataType = "String",position = 105)
    private String name;

    @ApiModelProperty(value = "标签",required = false,dataType = "String",position = 107)
    private String tags;

    private Integer status;

    private String version;

    @ApiModelProperty(value = "角色类型(0.管理角色 1.使用角色)",
            required = false,dataType = "Integer",position = 108)
    private Integer roleType;

    @ApiModelProperty(value = "角色名",required = false,dataType = "String",position = 109)
    private String roleName;
    @ApiModelProperty(value = "角色集描述",required = false,dataType = "String",position = 110)
    private String displayName;




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


    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }



    public Integer getStatus() {
        return status;
    }

    public void setStatus(Integer status) {
        this.status = status;
    }

    public String getTags() {
        return tags;
    }

    public void setTags(String tags) {
        this.tags = tags;
    }


    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version;
    }

    public Integer getRoleType() {
        return roleType;
    }

    public void setRoleType(Integer roleType) {
        this.roleType = roleType;
    }

    public String getRoleName() {
        return roleName;
    }

    public void setRoleName(String roleName) {
        this.roleName = roleName;
    }

    public String getDisplayName() {
        return displayName;
    }

    public void setDisplayName(String displayName) {
        this.displayName = displayName;
    }
}
