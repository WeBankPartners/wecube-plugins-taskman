package com.webank.taskman.dto.resp;


import java.util.Date;
import java.util.List;
import java.util.StringJoiner;

public class RequestInfoResqDto  {

    private String id;
    private String requestTempId;
    private String requestTempName;
    private String procInstId;
    private String rootEntity;
    private String name;
    private String description;
    private String reporter;

    private Date reportTime;
    private String emergency;
    private String reportRole;
    private String status;
    private String dueDate;
    private String result;
    private String attachFileId;
    private List<FormItemInfoRespDto> formItemInfos;

    public String getId() {
        return id;
    }

    public RequestInfoResqDto setId(String id) {
        this.id = id;
        return this;
    }

    public String getRequestTempId() {
        return requestTempId;
    }

    public RequestInfoResqDto setRequestTempId(String requestTempId) {
        this.requestTempId = requestTempId;
        return this;
    }

    public String getRequestTempName() {
        return requestTempName;
    }

    public RequestInfoResqDto setRequestTempName(String requestTempName) {
        this.requestTempName = requestTempName;
        return this;
    }

    public String getProcInstId() {
        return procInstId;
    }

    public RequestInfoResqDto setProcInstId(String procInstId) {
        this.procInstId = procInstId;
        return this;
    }

    public String getRootEntity() {
        return rootEntity;
    }

    public RequestInfoResqDto setRootEntity(String rootEntity) {
        this.rootEntity = rootEntity;
        return this;
    }

    public String getName() {
        return name;
    }

    public RequestInfoResqDto setName(String name) {
        this.name = name;
        return this;
    }

    public String getDescription() {
        return description;
    }

    public RequestInfoResqDto setDescription(String description) {
        this.description = description;
        return this;
    }

    public String getReporter() {
        return reporter;
    }

    public RequestInfoResqDto setReporter(String reporter) {
        this.reporter = reporter;
        return this;
    }

    public Date getReportTime() {
        return reportTime;
    }

    public RequestInfoResqDto setReportTime(Date reportTime) {
        this.reportTime = reportTime;
        return this;
    }

    public String getEmergency() {
        return emergency;
    }

    public RequestInfoResqDto setEmergency(String emergency) {
        this.emergency = emergency;
        return this;
    }

    public String getReportRole() {
        return reportRole;
    }

    public RequestInfoResqDto setReportRole(String reportRole) {
        this.reportRole = reportRole;
        return this;
    }

    public String getAttachFileId() {
        return attachFileId;
    }

    public RequestInfoResqDto setAttachFileId(String attachFileId) {
        this.attachFileId = attachFileId;
        return this;
    }

    public String getStatus() {
        return status;
    }

    public RequestInfoResqDto setStatus(String status) {
        this.status = status;
        return this;
    }

    public String getDueDate() {
        return dueDate;
    }

    public RequestInfoResqDto setDueDate(String dueDate) {
        this.dueDate = dueDate;
        return this;
    }

    public String getResult() {
        return result;
    }

    public RequestInfoResqDto setResult(String result) {
        this.result = result;
        return this;
    }

    public List<FormItemInfoRespDto> getFormItemInfos() {
        return formItemInfos;
    }

    public RequestInfoResqDto setFormItemInfos(List<FormItemInfoRespDto> formItemInfos) {
        this.formItemInfos = formItemInfos;
        return this;
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", RequestInfoResqDto.class.getSimpleName() + "[", "]")
                .add("id='" + id + "'")
                .add("requestTempId='" + requestTempId + "'")
                .add("requestTempName='" + requestTempName + "'")
                .add("procInstId='" + procInstId + "'")
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
                .toString();
    }
}