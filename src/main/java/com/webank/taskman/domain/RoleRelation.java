package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;

import java.io.Serializable;

public class RoleRelation  implements Serializable {

    private static final long serialVersionUID = 1L;

    
    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;

    private String recordTable;

    private String recordId;

    private Integer roleType;

    private String roleName;

    private String displayName;

    public RoleRelation() {
    }

    public RoleRelation(String recordTable, String recordId, Integer roleType, String roleName, String displayName) {
        this.recordTable = recordTable;
        this.recordId = recordId;
        this.roleType = roleType;
        this.roleName = roleName;
        this.displayName = displayName;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getRecordTable() {
        return recordTable;
    }

    public void setRecordTable(String recordTable) {
        this.recordTable = recordTable;
    }

    public String getRecordId() {
        return recordId;
    }

    public void setRecordId(String recordId) {
        this.recordId = recordId;
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

    @Override
    public String toString() {
        return "RoleRelation{" +
                "id='" + id + '\'' +
                ", recordTable='" + recordTable + '\'' +
                ", recordId='" + recordId + '\'' +
                ", roleType=" + roleType +
                ", roleName='" + roleName + '\'' +
                ", displayName='" + displayName + '\'' +
                '}';
    }
}
