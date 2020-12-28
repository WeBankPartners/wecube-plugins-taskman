package com.webank.taskman.dto;



import java.util.Date;
import java.util.StringJoiner;

public class RequestInfoDTO {

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
    private String description;

    public String getId() {
        return id;
    }

    public RequestInfoDTO setId(String id) {
        this.id = id;
        return this;
    }

    public String getRequestTempId() {
        return requestTempId;
    }

    public RequestInfoDTO setRequestTempId(String requestTempId) {
        this.requestTempId = requestTempId;
        return this;
    }

    public String getProcInstKey() {
        return procInstKey;
    }

    public RequestInfoDTO setProcInstKey(String procInstKey) {
        this.procInstKey = procInstKey;
        return this;
    }

    public String getRootEntity() {
        return rootEntity;
    }

    public RequestInfoDTO setRootEntity(String rootEntity) {
        this.rootEntity = rootEntity;
        return this;
    }

    public String getName() {
        return name;
    }

    public RequestInfoDTO setName(String name) {
        this.name = name;
        return this;
    }

    public String getReporter() {
        return reporter;
    }

    public RequestInfoDTO setReporter(String reporter) {
        this.reporter = reporter;
        return this;
    }

    public Date getReportTime() {
        return reportTime;
    }

    public RequestInfoDTO setReportTime(Date reportTime) {
        this.reportTime = reportTime;
        return this;
    }

    public String getEmergency() {
        return emergency;
    }

    public RequestInfoDTO setEmergency(String emergency) {
        this.emergency = emergency;
        return this;
    }

    public String getReportRole() {
        return reportRole;
    }

    public RequestInfoDTO setReportRole(String reportRole) {
        this.reportRole = reportRole;
        return this;
    }

    public String getAttachFileId() {
        return attachFileId;
    }

    public RequestInfoDTO setAttachFileId(String attachFileId) {
        this.attachFileId = attachFileId;
        return this;
    }

    public String getStatus() {
        return status;
    }

    public RequestInfoDTO setStatus(String status) {
        this.status = status;
        return this;
    }

    public String getDueDate() {
        return dueDate;
    }

    public RequestInfoDTO setDueDate(String dueDate) {
        this.dueDate = dueDate;
        return this;
    }

    public String getResult() {
        return result;
    }

    public RequestInfoDTO setResult(String result) {
        this.result = result;
        return this;
    }

    public String getDescription() {
        return description;
    }

    public RequestInfoDTO setDescription(String description) {
        this.description = description;
        return this;
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", RequestInfoDTO.class.getSimpleName() + "[", "]")
                .add("id='" + id + "'")
                .add("requestTempId='" + requestTempId + "'")
                .add("procInstKey='" + procInstKey + "'")
                .add("rootEntity='" + rootEntity + "'")
                .add("name='" + name + "'")
                .add("reporter='" + reporter + "'")
                .add("reportTime=" + reportTime)
                .add("emergency='" + emergency + "'")
                .add("reportRole='" + reportRole + "'")
                .add("attachFileId='" + attachFileId + "'")
                .add("status='" + status + "'")
                .add("dueDate='" + dueDate + "'")
                .add("result='" + result + "'")
                .add("description='" + description + "'")
                .toString();
    }
}