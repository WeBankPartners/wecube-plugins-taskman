package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.webank.taskman.base.BaseEntity;

import java.io.Serializable;

public class RequestTemplateRole extends BaseEntity implements Serializable {

    private static final long serialVersionUID = 1L;

    
    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;

    
    private String requestTemplateId;

    
    private String roleId;

    
    private Integer roleType;


    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getRequestTemplateId() {
        return requestTemplateId;
    }

    public void setRequestTemplateId(String requestTemplateId) {
        this.requestTemplateId = requestTemplateId;
    }

    public String getRoleId() {
        return roleId;
    }

    public void setRoleId(String roleId) {
        this.roleId = roleId;
    }

    public Integer getRoleType() {
        return roleType;
    }

    public void setRoleType(Integer roleType) {
        this.roleType = roleType;
    }


    public RequestTemplateRole() {
    }

    public RequestTemplateRole(String requestTemplateId, String roleId, Integer roleType) {
        this.requestTemplateId = requestTemplateId;
        this.roleId = roleId;
        this.roleType = roleType;
    }

    @Override
    public String toString() {
        return "RequestTemplaeRole{" +
        "id=" + id +
        ", requestTemplateId=" + requestTemplateId +
        ", roleId=" + roleId +
        ", roleType=" + roleType +
        "}";
    }
}
