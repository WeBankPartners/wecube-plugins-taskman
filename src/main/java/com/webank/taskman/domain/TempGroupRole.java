package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import java.io.Serializable;

/**
 * <p>
 * 模板角色关系表 
 * </p>
 *
 * @author ${author}
 * @since 2020-11-26
 */
public class TempGroupRole implements Serializable {

    private static final long serialVersionUID = 1L;

    /**
     * 主键
     */
    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;

    /**
     * 模板组编号
     */
    private String groupId;

    /**
     * 角色编号
     */
    private String roleId;

    /**
     * 类型
     */
    private Integer type;


    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getGroupId() {
        return groupId;
    }

    public void setGroupId(String groupId) {
        this.groupId = groupId;
    }

    public String getRoleId() {
        return roleId;
    }

    public void setRoleId(String roleId) {
        this.roleId = roleId;
    }

    public Integer getType() {
        return type;
    }

    public void setType(Integer type) {
        this.type = type;
    }

    @Override
    public String toString() {
        return "TempGroupRole{" +
        "id=" + id +
        ", groupId=" + groupId +
        ", roleId=" + roleId +
        ", type=" + type +
        "}";
    }
}
