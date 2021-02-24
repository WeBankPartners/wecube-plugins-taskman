package com.webank.taskman.dto.resp;

import java.util.Date;
import java.util.List;
import java.util.StringJoiner;

public class TaskInfoResp {
    private String id;
    private String parentId;
    private String taskTempId;
    private String procInstId;
    private String nodeDefId;
    private String nodeName;
    private String name;
    private String description;
    private String reporter;
    private Date reportTime;
    private String emergency;
    private String result;
    private String status;
    private String version;
    private String reportRole;
    private String attachFileId;

    private List<FormItemInfoResp> formItemInfo;

    public String getId() {
        return id;
    }

    public TaskInfoResp setId(String id) {
        this.id = id;
        return this;
    }

    public String getParentId() {
        return parentId;
    }

    public TaskInfoResp setParentId(String parentId) {
        this.parentId = parentId;
        return this;
    }

    public String getTaskTempId() {
        return taskTempId;
    }

    public TaskInfoResp setTaskTempId(String taskTempId) {
        this.taskTempId = taskTempId;
        return this;
    }

    public String getProcInstId() {
        return procInstId;
    }

    public TaskInfoResp setProcInstId(String procInstId) {
        this.procInstId = procInstId;
        return this;
    }

    public String getNodeDefId() {
        return nodeDefId;
    }

    public TaskInfoResp setNodeDefId(String nodeDefId) {
        this.nodeDefId = nodeDefId;
        return this;
    }

    public String getNodeName() {
        return nodeName;
    }

    public TaskInfoResp setNodeName(String nodeName) {
        this.nodeName = nodeName;
        return this;
    }

    public String getName() {
        return name;
    }

    public TaskInfoResp setName(String name) {
        this.name = name;
        return this;
    }

    public String getDescription() {
        return description;
    }

    public TaskInfoResp setDescription(String description) {
        this.description = description;
        return this;
    }

    public String getReporter() {
        return reporter;
    }

    public TaskInfoResp setReporter(String reporter) {
        this.reporter = reporter;
        return this;
    }

    public Date getReportTime() {
        return reportTime;
    }

    public TaskInfoResp setReportTime(Date reportTime) {
        this.reportTime = reportTime;
        return this;
    }

    public String getEmergency() {
        return emergency;
    }

    public TaskInfoResp setEmergency(String emergency) {
        this.emergency = emergency;
        return this;
    }

    public String getResult() {
        return result;
    }

    public TaskInfoResp setResult(String result) {
        this.result = result;
        return this;
    }

    public String getStatus() {
        return status;
    }

    public TaskInfoResp setStatus(String status) {
        this.status = status;
        return this;
    }

    public String getVersion() {
        return version;
    }

    public TaskInfoResp setVersion(String version) {
        this.version = version;
        return this;
    }

    public String getReportRole() {
        return reportRole;
    }

    public TaskInfoResp setReportRole(String reportRole) {
        this.reportRole = reportRole;
        return this;
    }

    public String getAttachFileId() {
        return attachFileId;
    }

    public TaskInfoResp setAttachFileId(String attachFileId) {
        this.attachFileId = attachFileId;
        return this;
    }

    public List<FormItemInfoResp> getFormItemInfo() {
        return formItemInfo;
    }

    public TaskInfoResp setFormItemInfo(List<FormItemInfoResp> formItemInfo) {
        this.formItemInfo = formItemInfo;
        return this;
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", TaskInfoResp.class.getSimpleName() + "[", "]")
                .add("id='" + id + "'")
                .add("parentId='" + parentId + "'")
                .add("taskTempId='" + taskTempId + "'")
                .add("procInstId='" + procInstId + "'")
                .add("nodeDefId='" + nodeDefId + "'")
                .add("nodeName='" + nodeName + "'")
                .add("name='" + name + "'")
                .add("description='" + description + "'")
                .add("reporter='" + reporter + "'")
                .add("reportTime=" + reportTime)
                .add("emergency='" + emergency + "'")
                .add("result='" + result + "'")
                .add("status='" + status + "'")
                .add("version='" + version + "'")
                .add("reportRole='" + reportRole + "'")
                .add("attachFileId='" + attachFileId + "'")
                .add("formItemInfo=" + formItemInfo)
                .toString();
    }
}
