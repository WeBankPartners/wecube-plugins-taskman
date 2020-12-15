package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.webank.taskman.base.BaseEntity;

import java.io.Serializable;
import java.util.Date;

public class RequestInfo extends BaseEntity implements Serializable {

    private static final long serialVersionUID = 1L;

    
    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;

    private String requestTempId;
    private String procInstKey;
    private String name;
    private String reporter;
    private Date reportTime;
    private String emergency;
    private String reportRole;
    private String attachFileId;
    private String status;
    private String result;


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

    public String getReporter() {
        return reporter;
    }

    public void setReporter(String reporter) {
        this.reporter = reporter;
    }

    public Date getReportTime() {
        return reportTime;
    }

    public void setReportTime(Date reportTime) {
        this.reportTime = reportTime;
    }

    public String getEmergency() {
        return emergency;
    }

    public void setEmergency(String emergency) {
        this.emergency = emergency;
    }

    public String getReportRole() {
        return reportRole;
    }

    public void setReportRole(String reportRole) {
        this.reportRole = reportRole;
    }

    public String getAttachFileId() {
        return attachFileId;
    }

    public void setAttachFileId(String attachFileId) {
        this.attachFileId = attachFileId;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public String getResult() {
        return result;
    }

    public void setResult(String result) {
        this.result = result;
    }

    @Override
    public String toString() {
        return "RequestInfo{" +
        "id=" + id +
        ", requestTempId=" + requestTempId +
        ", procInstKey=" + procInstKey +
        ", name=" + name +
        ", reporter=" + reporter +
        ", reportTime=" + reportTime +
        ", emergency=" + emergency +
        ", reportRole=" + reportRole +
        ", attachFileId=" + attachFileId +
        ", status=" + status +
        ", result=" + result +
        ", createdBy=" + getCreatedBy() +
        ", createdTime=" + getCreatedTime() +
        ", updatedBy=" + getUpdatedBy() +
        ", updatedTime=" + getUpdatedTime() +
        ", delFlag=" + getDelFlag() +
        "}";
    }
}
