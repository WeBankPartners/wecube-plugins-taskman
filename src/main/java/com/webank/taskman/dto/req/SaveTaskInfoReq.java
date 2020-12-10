package com.webank.taskman.dto.req;


import io.swagger.annotations.ApiModel;
import org.apache.commons.lang3.StringUtils;

@ApiModel
public class SaveTaskInfoReq {

    private String id;

    private String taskTempId;

    private String name;

    private Integer status;

    private String useRoleName;

    private String manageRoleName;

    private Integer roleType;

    private String roleName;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getTaskTempId() {
        return taskTempId;
    }

    public void setTaskTempId(String taskTempId) {
        this.taskTempId = taskTempId;
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

    @Override
    public String toString() {
        return "SaveTaskInfoReq{" +
                "id='" + id + '\'' +
                ", taskTempId='" + taskTempId + '\'' +
                ", name='" + name + '\'' +
                ", status=" + status +
                ", useRoleName='" + useRoleName + '\'' +
                ", manageRoleName='" + manageRoleName + '\'' +
                ", roleType=" + roleType +
                ", roleName='" + roleName + '\'' +
                '}';
    }
}
