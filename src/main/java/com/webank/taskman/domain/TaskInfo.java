package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.webank.taskman.base.BaseEntity;
import org.springframework.util.StringUtils;

import java.io.Serializable;
import java.util.Date;

public class TaskInfo extends BaseEntity implements Serializable {

    private static final long serialVersionUID = 1L;

    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;
    private String requestId;
    private String parentId;
    private String taskTempId;
    private String procInstId;
    private String nodeDefId;
    private String nodeName;
    private String callbackUrl;
    private String callbackParameter;
    private String name;
    private String reporter;
    private Date reportTime;
    private String emergency;
    private String reportRole;
    private String result;
    private String description;
    private String attachFileId;

    private String status;

    private String version;

    private String overTime;

    public TaskInfo() {
    }

    public TaskInfo(String parentId, String nodeDefId) {
        this.parentId = parentId;
        this.nodeDefId = nodeDefId;
    }

    public TaskInfo(String id, String parentId, String taskTempId, String nodeDefId, String nodeName,
            String callbackUrl, String callbackParameter, String name, String reporter, Date reportTime,
            String emergency, String reportRole, String result, String description, String attachFileId, String status,
            String version) {
        this.id = id;
        this.parentId = parentId;
        this.taskTempId = taskTempId;
        this.nodeDefId = nodeDefId;
        this.nodeName = nodeName;
        this.callbackUrl = callbackUrl;
        this.callbackParameter = callbackParameter;
        this.name = name;
        this.reporter = reporter;
        this.reportTime = reportTime;
        this.emergency = emergency;
        this.reportRole = reportRole;
        this.result = result;
        this.description = description;
        this.attachFileId = attachFileId;
        this.status = status;
        this.version = version;
    }

    @JsonIgnore
    public LambdaQueryWrapper<TaskInfo> getLambdaQueryWrapper() {
        return new LambdaQueryWrapper<TaskInfo>().eq(!StringUtils.isEmpty(id), TaskInfo::getId, id)
                .eq(!StringUtils.isEmpty(parentId), TaskInfo::getParentId, parentId)
                .eq(!StringUtils.isEmpty(taskTempId), TaskInfo::getTaskTempId, taskTempId)
                .eq(!StringUtils.isEmpty(nodeDefId), TaskInfo::getNodeDefId, nodeDefId)
                .eq(!StringUtils.isEmpty(nodeName), TaskInfo::getNodeName, nodeName)
                .eq(!StringUtils.isEmpty(callbackUrl), TaskInfo::getCallbackUrl, callbackUrl)
                .eq(!StringUtils.isEmpty(callbackParameter), TaskInfo::getCallbackParameter, callbackParameter)
                .like(!StringUtils.isEmpty(name), TaskInfo::getName, name)
                .eq(!StringUtils.isEmpty(reporter), TaskInfo::getReporter, reporter)
                .eq(!StringUtils.isEmpty(reportTime), TaskInfo::getReportTime, reportTime)
                .eq(!StringUtils.isEmpty(emergency), TaskInfo::getEmergency, emergency)
                .eq(!StringUtils.isEmpty(reportRole), TaskInfo::getReportRole, reportRole)
                .eq(!StringUtils.isEmpty(result), TaskInfo::getResult, result)
                .eq(!StringUtils.isEmpty(description), TaskInfo::getDescription, description)
                .eq(!StringUtils.isEmpty(attachFileId), TaskInfo::getAttachFileId, attachFileId)
                .eq(!StringUtils.isEmpty(status), TaskInfo::getStatus, status)
                .eq(!StringUtils.isEmpty(version), TaskInfo::getVersion, version);
    }

    public String getId() {
        return id;
    }

    public TaskInfo setId(String id) {
        this.id = id;
        return this;
    }

    public String getRequestId() {
        return requestId;
    }

    public TaskInfo setRequestId(String requestId) {
        this.requestId = requestId;
        return this;
    }

    public String getParentId() {
        return parentId;
    }

    public TaskInfo setParentId(String parentId) {
        this.parentId = parentId;
        return this;
    }

    public String getTaskTempId() {
        return taskTempId;
    }

    public TaskInfo setTaskTempId(String taskTempId) {
        this.taskTempId = taskTempId;
        return this;
    }

    public String getProcInstId() {
        return procInstId;
    }

    public TaskInfo setProcInstId(String procInstId) {
        this.procInstId = procInstId;
        return this;
    }

    public String getNodeDefId() {
        return nodeDefId;
    }

    public TaskInfo setNodeDefId(String nodeDefId) {
        this.nodeDefId = nodeDefId;
        return this;
    }

    public String getNodeName() {
        return nodeName;
    }

    public TaskInfo setNodeName(String nodeName) {
        this.nodeName = nodeName;
        return this;
    }

    public String getCallbackUrl() {
        return callbackUrl;
    }

    public TaskInfo setCallbackUrl(String callbackUrl) {
        this.callbackUrl = callbackUrl;
        return this;
    }

    public String getCallbackParameter() {
        return callbackParameter;
    }

    public TaskInfo setCallbackParameter(String callbackParameter) {
        this.callbackParameter = callbackParameter;
        return this;
    }

    public String getName() {
        return name;
    }

    public TaskInfo setName(String name) {
        this.name = name;
        return this;
    }

    public String getReporter() {
        return reporter;
    }

    public TaskInfo setReporter(String reporter) {
        this.reporter = reporter;
        return this;
    }

    public Date getReportTime() {
        return reportTime;
    }

    public TaskInfo setReportTime(Date reportTime) {
        this.reportTime = reportTime;
        return this;
    }

    public String getEmergency() {
        return emergency;
    }

    public TaskInfo setEmergency(String emergency) {
        this.emergency = emergency;
        return this;
    }

    public String getReportRole() {
        return reportRole;
    }

    public TaskInfo setReportRole(String reportRole) {
        this.reportRole = reportRole;
        return this;
    }

    public String getResult() {
        return result;
    }

    public TaskInfo setResult(String result) {
        this.result = result;
        return this;
    }

    public String getDescription() {
        return description;
    }

    public TaskInfo setDescription(String description) {
        this.description = description;
        return this;
    }

    public String getAttachFileId() {
        return attachFileId;
    }

    public TaskInfo setAttachFileId(String attachFileId) {
        this.attachFileId = attachFileId;
        return this;
    }

    public String getStatus() {
        return status;
    }

    public TaskInfo setStatus(String status) {
        this.status = status;
        return this;
    }

    public String getVersion() {
        return version;
    }

    public TaskInfo setVersion(String version) {
        this.version = version;
        return this;
    }

    public String getOverTime() {
        return overTime;
    }

    public TaskInfo setOverTime(String overTime) {
        this.overTime = overTime;
        return this;
    }

    @Override
    public String toString() {
        return "TaskInfo{" + "id='" + id + '\'' + ", parentId='" + parentId + '\'' + ", taskTempId='" + taskTempId
                + '\'' + ", nodeDefId='" + nodeDefId + '\'' + ", nodeName='" + nodeName + '\'' + ", callbackUrl='"
                + callbackUrl + '\'' + ", callbackParameter='" + callbackParameter + '\'' + ", name='" + name + '\''
                + ", reporter='" + reporter + '\'' + ", reportTime=" + reportTime + ", emergency='" + emergency + '\''
                + ", reportRole='" + reportRole + '\'' + ", result='" + result + '\'' + ", description='" + description
                + '\'' + ", attachFileId='" + attachFileId + '\'' + ", status=" + status + ", version='" + version
                + '\'' + '}';
    }
}
