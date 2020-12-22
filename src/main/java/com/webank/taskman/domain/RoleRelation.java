package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.fasterxml.jackson.annotation.JsonIgnore;
import org.springframework.util.StringUtils;

import java.io.Serializable;

public class RoleRelation  implements Serializable {

    private static final long serialVersionUID = 1L;

    
    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;
    private String recordId;
    private Integer roleType;
    private String roleName;
    private String displayName;

    public RoleRelation() {
    }

    public RoleRelation( String recordId, Integer roleType, String roleName, String displayName) {
        this.recordId = recordId;
        this.roleType = roleType;
        this.roleName = roleName;
        this.displayName = displayName;
    }

    @JsonIgnore
    public LambdaQueryWrapper getLambdaQueryWrapper() {
        return new LambdaQueryWrapper<RoleRelation>()
                .eq(!StringUtils.isEmpty(id), RoleRelation::getId, id)
                .eq(!StringUtils.isEmpty(recordId), RoleRelation::getRecordId, recordId)
                .eq(!StringUtils.isEmpty(roleType), RoleRelation::getRoleType, roleType)
                .eq(!StringUtils.isEmpty(roleName), RoleRelation::getRoleName, roleName)
                .like(!StringUtils.isEmpty(displayName), RoleRelation::getDisplayName, displayName)
                ;
    }

    public String getId() {
        return id;
    }

    public RoleRelation setId(String id) {
        this.id = id;
        return this;
    }

    public String getRecordId() {
        return recordId;
    }

    public RoleRelation setRecordId(String recordId) {
        this.recordId = recordId;
        return this;
    }

    public Integer getRoleType() {
        return roleType;
    }

    public RoleRelation setRoleType(Integer roleType) {
        this.roleType = roleType;
        return this;
    }

    public String getRoleName() {
        return roleName;
    }

    public RoleRelation setRoleName(String roleName) {
        this.roleName = roleName;
        return this;
    }

    public String getDisplayName() {
        return displayName;
    }

    public RoleRelation setDisplayName(String displayName) {
        this.displayName = displayName;
        return this;
    }

    @Override
    public String toString() {
        return "RoleRelation{" +
                "id='" + id + '\'' +
                ", recordId='" + recordId + '\'' +
                ", roleType=" + roleType +
                ", roleName='" + roleName + '\'' +
                ", displayName='" + displayName + '\'' +
                '}';
    }
}
