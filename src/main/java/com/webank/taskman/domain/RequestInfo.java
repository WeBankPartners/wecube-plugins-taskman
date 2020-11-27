package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableLogic;

import java.io.Serializable;
import java.util.Date;

public class RequestInfo implements Serializable {

    private static final long serialVersionUID = 1L;

    
    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;

    
    private String formId;

    
    private String requestTempId;

    
    private String proc;

    
    private String name;

    
    private Integer status;

    
    private String createdBy;

    
    private Date createdTime;

    
    private String updatedBy;

    
    private Date updatedTime;


    @TableLogic
    private Integer delFlag;


    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getFormId() {
        return formId;
    }

    public void setFormId(String formId) {
        this.formId = formId;
    }

    public String getRequestTempId() {
        return requestTempId;
    }

    public void setRequestTempId(String requestTempId) {
        this.requestTempId = requestTempId;
    }

    public String getProc() {
        return proc;
    }

    public void setProc(String proc) {
        this.proc = proc;
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
        return "RequestInfo{" +
        "id=" + id +
        ", formId=" + formId +
        ", requestTempId=" + requestTempId +
        ", proc=" + proc +
        ", name=" + name +
        ", status=" + status +
        ", createdBy=" + createdBy +
        ", createdTime=" + createdTime +
        ", updatedBy=" + updatedBy +
        ", updatedTime=" + updatedTime +
        ", delFlag=" + delFlag +
        "}";
    }
}
