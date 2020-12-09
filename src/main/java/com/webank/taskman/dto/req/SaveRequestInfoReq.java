package com.webank.taskman.dto.req;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import org.apache.commons.lang3.StringUtils;

@ApiModel(value = "AddRequestInfoReq",description = "add RequestInfo req")
public class SaveRequestInfoReq {
    @ApiModelProperty(value = "主键",required = false,dataType = "String",position = 100)
    private String id;

    @ApiModelProperty(value = "请求模板id",required = false,dataType = "String",position = 101)
    private String requestTempId;

    @ApiModelProperty(value = "流程编排key",required = false,dataType = "String",position = 102)
    private String procInstKey;

    @ApiModelProperty(value = "请求信息名称",required = false,dataType = "String",position = 103)
    private String name;

    @ApiModelProperty(value = "发布状态",required = false,dataType = "Integer",position = 104)
    private Integer status;

    @ApiModelProperty(value = "使用角色",required = false,position = 109)
    private String useRoleName;

    @ApiModelProperty(value = "管理角色",required = false,position = 110)
    private String manageRoleName;

    @ApiModelProperty(hidden = true)
    private Integer roleType;

    @ApiModelProperty(hidden = true)
    private String roleName;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getRequestTempId() {
        return requestTempId;
    }

    public void setRequestTempId(String requestTempId) {
        this.requestTempId = requestTempId;
    }

    public String getProcInstKey() {
        return procInstKey;
    }

    public void setProcInstKey(String procInstKey) {
        this.procInstKey = procInstKey;
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
    public String getUseRoleName() {
        return useRoleName;
    }

    public void setUseRoleName(String useRoleName) {
        this.useRoleName = useRoleName;
        if(!StringUtils.isEmpty(useRoleName)){
            this.roleType = null != this.roleType ? 2:1;
            this.roleName = this.useRoleName;
        }
    }

    public String getManageRoleName() {
        return manageRoleName;
    }

    public void setManageRoleName(String manageRoleName) {
        this.manageRoleName = manageRoleName;
        if(!StringUtils.isEmpty(manageRoleName)){
            this.roleType = null != this.roleType ? 2:0;
            this.roleName = this.manageRoleName;
        }
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
}
