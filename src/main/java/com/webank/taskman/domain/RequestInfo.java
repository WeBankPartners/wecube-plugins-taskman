package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.webank.taskman.base.BaseEntity;
import org.springframework.util.StringUtils;

import java.io.Serializable;
import java.util.Date;

public class RequestInfo extends BaseEntity implements Serializable {

    private static final long serialVersionUID = 1L;

    
    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;

    private String requestTempId;
    private String procInstKey;
    private String rootEntity;
    private String name;
    private String reporter;
    private Date reportTime;
    private String emergency;
    private String reportRole;
    private String attachFileId;
    private String status;
    private String dueDate;
    private String result;

    public RequestInfo() {
    }

    public RequestInfo(String id, String requestTempId, String procInstKey, String rootEntity, String name, String reporter, Date reportTime, String emergency, String reportRole, String attachFileId, String status, String dueDate, String result) {
        this.id = id;
        this.requestTempId = requestTempId;
        this.procInstKey = procInstKey;
        this.rootEntity = rootEntity;
        this.name = name;
        this.reporter = reporter;
        this.reportTime = reportTime;
        this.emergency = emergency;
        this.reportRole = reportRole;
        this.attachFileId = attachFileId;
        this.status = status;
        this.dueDate = dueDate;
        this.result = result;
    }

    @JsonIgnore
    public LambdaQueryWrapper getLambdaQueryWrapper() {
        return new LambdaQueryWrapper<RequestInfo>()
            .eq(!StringUtils.isEmpty(id), RequestInfo::getId, id)
            .eq(!StringUtils.isEmpty(requestTempId), RequestInfo::getRequestTempId, requestTempId)
            .eq(!StringUtils.isEmpty(procInstKey), RequestInfo::getProcInstKey, procInstKey)
            .eq(!StringUtils.isEmpty(rootEntity), RequestInfo::getRootEntity, rootEntity)
            .like(!StringUtils.isEmpty(name), RequestInfo::getName, name)
            .eq(!StringUtils.isEmpty(reporter), RequestInfo::getReporter, reporter)
            .eq(!StringUtils.isEmpty(reportTime), RequestInfo::getReportTime, reportTime)
            .eq(!StringUtils.isEmpty(emergency), RequestInfo::getEmergency, emergency)
            .eq(!StringUtils.isEmpty(reportRole), RequestInfo::getReportRole, reportRole)
            .eq(!StringUtils.isEmpty(attachFileId), RequestInfo::getAttachFileId, attachFileId)
            .eq(!StringUtils.isEmpty(status), RequestInfo::getStatus, status)
            .eq(!StringUtils.isEmpty(dueDate), RequestInfo::getDueDate, dueDate)
            .eq(!StringUtils.isEmpty(result), RequestInfo::getResult, result);
    }


    public String getId() {
        return id;
    }

    public RequestInfo setId(String id) {
        this.id = id;
        return this;
    }

    public String getRequestTempId() {
        return requestTempId;
    }

    public RequestInfo setRequestTempId(String requestTempId) {
        this.requestTempId = requestTempId;
        return this;
    }

    public String getProcInstKey() {
        return procInstKey;
    }

    public RequestInfo setProcInstKey(String procInstKey) {
        this.procInstKey = procInstKey;
        return this;
    }

    public String getRootEntity() {
        return rootEntity;
    }

    public RequestInfo setRootEntity(String rootEntity) {
        this.rootEntity = rootEntity;
        return this;
    }

    public String getName() {
        return name;
    }

    public RequestInfo setName(String name) {
        this.name = name;
        return this;
    }

    public String getReporter() {
        return reporter;
    }

    public RequestInfo setReporter(String reporter) {
        this.reporter = reporter;
        return this;
    }

    public Date getReportTime() {
        return reportTime;
    }

    public RequestInfo setReportTime(Date reportTime) {
        this.reportTime = reportTime;
        return this;
    }

    public String getEmergency() {
        return emergency;
    }

    public RequestInfo setEmergency(String emergency) {
        this.emergency = emergency;
        return this;
    }

    public String getReportRole() {
        return reportRole;
    }

    public RequestInfo setReportRole(String reportRole) {
        this.reportRole = reportRole;
        return this;
    }

    public String getAttachFileId() {
        return attachFileId;
    }

    public RequestInfo setAttachFileId(String attachFileId) {
        this.attachFileId = attachFileId;
        return this;
    }

    public String getStatus() {
        return status;
    }

    public RequestInfo setStatus(String status) {
        this.status = status;
        return this;
    }

    public String getDueDate() {
        return dueDate;
    }

    public RequestInfo setDueDate(String dueDate) {
        this.dueDate = dueDate;
        return this;
    }

    public String getResult() {
        return result;
    }

    public RequestInfo setResult(String result) {
        this.result = result;
        return this;
    }
}
