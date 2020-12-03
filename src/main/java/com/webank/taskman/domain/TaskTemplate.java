package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import java.util.Date;
import com.baomidou.mybatisplus.annotation.TableId;
import java.io.Serializable;

/**
 * <p>
 * 任务模板信息表 
 * </p>
 *
 * @author ${author}
 * @since 2020-11-26
 */
public class TaskTemplate implements Serializable {

    private static final long serialVersionUID = 1L;

    /**
     * 主键
     */
    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;

    /**
     * 表单模板编号
     */
    private String formTempId;

    /**
     * 请求模板编号
     */
    private String requestTempId;

    /**
     * 所属角色
     */
    private String dealRole;

    /**
     * 流程节点模板
     */
    private String procNodeTemp;

    /**
     * 名称
     */
    private String name;

    /**
     * 创建人
     */
    private String createdBy;

    /**
     * 创建时间
     */
    private Date createdTime;

    /**
     * 更新人
     */
    private String updatedBy;

    /**
     * 更新时间
     */
    private Date updatedTime;

    /**
     * 是否删除
     */
    private Integer delFlag;


    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getFormTempId() {
        return formTempId;
    }

    public void setFormTempId(String formTempId) {
        this.formTempId = formTempId;
    }

    public String getRequestTempId() {
        return requestTempId;
    }

    public void setRequestTempId(String requestTempId) {
        this.requestTempId = requestTempId;
    }

    public String getDealRole() {
        return dealRole;
    }

    public void setDealRole(String dealRole) {
        this.dealRole = dealRole;
    }

    public String getProcNodeTemp() {
        return procNodeTemp;
    }

    public void setProcNodeTemp(String procNodeTemp) {
        this.procNodeTemp = procNodeTemp;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getCreatedBy() {
        return createdBy;
    }

    public void setCreatedBy(String createdBy) {
        this.createdBy = createdBy;
    }

    public Date getCreatedTime() {
        return createdTime;
    }

    public void setCreatedTime(Date createdTime) {
        this.createdTime = createdTime;
    }

    public String getUpdatedBy() {
        return updatedBy;
    }

    public void setUpdatedBy(String updatedBy) {
        this.updatedBy = updatedBy;
    }

    public Date getUpdatedTime() {
        return updatedTime;
    }

    public void setUpdatedTime(Date updatedTime) {
        this.updatedTime = updatedTime;
    }

    public Integer getDelFlag() {
        return delFlag;
    }

    public void setDelFlag(Integer delFlag) {
        this.delFlag = delFlag;
    }

    @Override
    public String toString() {
        return "TaskTemplate{" +
        "id=" + id +
        ", formTempId=" + formTempId +
        ", requestTempId=" + requestTempId +
        ", dealRole=" + dealRole +
        ", procNodeTemp=" + procNodeTemp +
        ", name=" + name +
        ", createdBy=" + createdBy +
        ", createdTime=" + createdTime +
        ", updatedBy=" + updatedBy +
        ", updatedTime=" + updatedTime +
        ", delFlag=" + delFlag +
        "}";
    }
}
