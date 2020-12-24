package com.webank.taskman.dto.resp;

import io.swagger.annotations.ApiModel;

import java.util.Date;
import java.util.List;

@ApiModel
public class RequestInfoInstanceResq {
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

    private List<TaskInfoInstanceResp> taskInfoInstanceResps;

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

    public String getRootEntity() {
        return rootEntity;
    }

    public void setRootEntity(String rootEntity) {
        this.rootEntity = rootEntity;
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

    public String getDueDate() {
        return dueDate;
    }

    public void setDueDate(String dueDate) {
        this.dueDate = dueDate;
    }

    public String getResult() {
        return result;
    }

    public void setResult(String result) {
        this.result = result;
    }

    public List<TaskInfoInstanceResp> getTaskInfoInstanceResps() {
        return taskInfoInstanceResps;
    }

    public void setTaskInfoInstanceResps(List<TaskInfoInstanceResp> taskInfoInstanceResps) {
        this.taskInfoInstanceResps = taskInfoInstanceResps;
    }

    @Override
    public String toString() {
        return "RequestInfoInstanceResq{" +
                "id='" + id + '\'' +
                ", requestTempId='" + requestTempId + '\'' +
                ", procInstKey='" + procInstKey + '\'' +
                ", rootEntity='" + rootEntity + '\'' +
                ", name='" + name + '\'' +
                ", reporter='" + reporter + '\'' +
                ", reportTime=" + reportTime +
                ", emergency='" + emergency + '\'' +
                ", reportRole='" + reportRole + '\'' +
                ", attachFileId='" + attachFileId + '\'' +
                ", status='" + status + '\'' +
                ", dueDate='" + dueDate + '\'' +
                ", result='" + result + '\'' +
                ", taskInfoInstanceResps=" + taskInfoInstanceResps +
                '}';
    }
}