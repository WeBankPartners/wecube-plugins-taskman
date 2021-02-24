package com.webank.taskman.dto.resp;


import java.util.Date;
import java.util.List;
import java.util.StringJoiner;

public class RequestInfoResq  {

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
    private List<FormItemInfoResp> formItemInfos;

    public String getId() {
        return id;
    }

    public RequestInfoResq setId(String id) {
        this.id = id;
        return this;
    }

    public String getRequestTempId() {
        return requestTempId;
    }

    public RequestInfoResq setRequestTempId(String requestTempId) {
        this.requestTempId = requestTempId;
        return this;
    }

    public String getRequestTempName() {
        return requestTempName;
    }

    public RequestInfoResq setRequestTempName(String requestTempName) {
        this.requestTempName = requestTempName;
        return this;
    }

    public String getProcInstId() {
        return procInstId;
    }

    public RequestInfoResq setProcInstId(String procInstId) {
        this.procInstId = procInstId;
        return this;
    }

    public String getRootEntity() {
        return rootEntity;
    }

    public RequestInfoResq setRootEntity(String rootEntity) {
        this.rootEntity = rootEntity;
        return this;
    }

    public String getName() {
        return name;
    }

    public RequestInfoResq setName(String name) {
        this.name = name;
        return this;
    }

    public String getDescription() {
        return description;
    }

    public RequestInfoResq setDescription(String description) {
        this.description = description;
        return this;
    }

    public String getReporter() {
        return reporter;
    }

    public RequestInfoResq setReporter(String reporter) {
        this.reporter = reporter;
        return this;
    }

    public Date getReportTime() {
        return reportTime;
    }

    public RequestInfoResq setReportTime(Date reportTime) {
        this.reportTime = reportTime;
        return this;
    }

    public String getEmergency() {
        return emergency;
    }

    public RequestInfoResq setEmergency(String emergency) {
        this.emergency = emergency;
        return this;
    }

    public String getReportRole() {
        return reportRole;
    }

    public RequestInfoResq setReportRole(String reportRole) {
        this.reportRole = reportRole;
        return this;
    }

    public String getAttachFileId() {
        return attachFileId;
    }

    public RequestInfoResq setAttachFileId(String attachFileId) {
        this.attachFileId = attachFileId;
        return this;
    }

    public String getStatus() {
        return status;
    }

    public RequestInfoResq setStatus(String status) {
        this.status = status;
        return this;
    }

    public String getDueDate() {
        return dueDate;
    }

    public RequestInfoResq setDueDate(String dueDate) {
        this.dueDate = dueDate;
        return this;
    }

    public String getResult() {
        return result;
    }

    public RequestInfoResq setResult(String result) {
        this.result = result;
        return this;
    }

    public List<FormItemInfoResp> getFormItemInfos() {
        return formItemInfos;
    }

    public RequestInfoResq setFormItemInfos(List<FormItemInfoResp> formItemInfos) {
        this.formItemInfos = formItemInfos;
        return this;
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", RequestInfoResq.class.getSimpleName() + "[", "]")
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