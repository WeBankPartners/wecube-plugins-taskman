package com.webank.taskman.dto;

import java.util.Date;
import java.util.StringJoiner;

import com.fasterxml.jackson.annotation.JsonInclude;

@JsonInclude(JsonInclude.Include.NON_NULL)
public class TaskInfoDto {

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
    private String reportRole;
    private String attachFileId;

    public TaskInfoDto() {
    }

    public TaskInfoDto(String id, String reporter, Date reportTime, String status) {
        this.id = id;
        this.reporter = reporter;
        this.reportTime = reportTime;
        this.status = status;
    }

    public String getId() {
        return id;
    }

    public TaskInfoDto setId(String id) {
        this.id = id;
        return this;
    }

    public String getTaskTempId() {
        return taskTempId;
    }

    public TaskInfoDto setTaskTempId(String taskTempId) {
        this.taskTempId = taskTempId;
        return this;
    }

    public String getParentId() {
        return parentId;
    }

    public TaskInfoDto setParentId(String parentId) {
        this.parentId = parentId;
        return this;
    }

    public String getProcInstId() {
        return procInstId;
    }

    public TaskInfoDto setProcInstId(String procInstId) {
        this.procInstId = procInstId;
        return this;
    }

    public String getNodeDefId() {
        return nodeDefId;
    }

    public TaskInfoDto setNodeDefId(String nodeDefId) {
        this.nodeDefId = nodeDefId;
        return this;
    }

    public String getNodeName() {
        return nodeName;
    }

    public TaskInfoDto setNodeName(String nodeName) {
        this.nodeName = nodeName;
        return this;
    }

    public String getName() {
        return name;
    }

    public TaskInfoDto setName(String name) {
        this.name = name;
        return this;
    }

    public String getDescription() {
        return description;
    }

    public TaskInfoDto setDescription(String description) {
        this.description = description;
        return this;
    }

    public String getReporter() {
        return reporter;
    }

    public TaskInfoDto setReporter(String reporter) {
        this.reporter = reporter;
        return this;
    }

    public Date getReportTime() {
        return reportTime;
    }

    public TaskInfoDto setReportTime(Date reportTime) {
        this.reportTime = reportTime;
        return this;
    }

    public String getEmergency() {
        return emergency;
    }

    public TaskInfoDto setEmergency(String emergency) {
        this.emergency = emergency;
        return this;
    }

    public String getResult() {
        return result;
    }

    public TaskInfoDto setResult(String result) {
        this.result = result;
        return this;
    }

    public String getStatus() {
        return status;
    }

    public TaskInfoDto setStatus(String status) {
        this.status = status;
        return this;
    }

    public String getReportRole() {
        return reportRole;
    }

    public TaskInfoDto setReportRole(String reportRole) {
        this.reportRole = reportRole;
        return this;
    }

    public String getAttachFileId() {
        return attachFileId;
    }

    public TaskInfoDto setAttachFileId(String attachFileId) {
        this.attachFileId = attachFileId;
        return this;
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", TaskInfoDto.class.getSimpleName() + "[", "]").add("id='" + id + "'")
                .add("parentId='" + parentId + "'").add("taskTempId='" + taskTempId + "'")
                .add("procInstId='" + procInstId + "'").add("nodeDefId='" + nodeDefId + "'")
                .add("nodeName='" + nodeName + "'").add("name='" + name + "'").add("description='" + description + "'")
                .add("reporter='" + reporter + "'").add("reportTime=" + reportTime).add("emergency='" + emergency + "'")
                .add("result='" + result + "'").add("status='" + status + "'").add("reportRole='" + reportRole + "'")
                .add("attachFileId='" + attachFileId + "'").toString();
    }
}
