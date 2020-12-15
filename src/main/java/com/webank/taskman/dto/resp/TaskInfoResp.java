package com.webank.taskman.dto.resp;


import io.swagger.annotations.ApiModel;


import java.util.Date;

@ApiModel
public class TaskInfoResp{

    private String id;

    private String requestId;
    
    private String requestNo;

    private String parentId;

    private String taskTempId;

    private String nodeDefId;

    private String nodeName;

    private String name;
    
    private String reporter;
    
    private Date reportTime;
    
    private String emergency;

    private String reportRole;

    private String result;

    private String description;

    private String attachFileId;

    private Integer status;

    private FormInfoResq formInfoResq;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getRequestId() {
        return requestId;
    }

    public void setRequestId(String requestId) {
        this.requestId = requestId;
    }

    public String getRequestNo() {
        return requestNo;
    }

    public void setRequestNo(String requestNo) {
        this.requestNo = requestNo;
    }

    public String getParentId() {
        return parentId;
    }

    public void setParentId(String parentId) {
        this.parentId = parentId;
    }

    public String getTaskTempId() {
        return taskTempId;
    }

    public void setTaskTempId(String taskTempId) {
        this.taskTempId = taskTempId;
    }

    public String getNodeDefId() {
        return nodeDefId;
    }

    public void setNodeDefId(String nodeDefId) {
        this.nodeDefId = nodeDefId;
    }

    public String getNodeName() {
        return nodeName;
    }

    public void setNodeName(String nodeName) {
        this.nodeName = nodeName;
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

    public String getResult() {
        return result;
    }

    public void setResult(String result) {
        this.result = result;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getAttachFileId() {
        return attachFileId;
    }

    public void setAttachFileId(String attachFileId) {
        this.attachFileId = attachFileId;
    }

    public Integer getStatus() {
        return status;
    }

    public void setStatus(Integer status) {
        this.status = status;
    }

    public FormInfoResq getFormInfoResq() {
        return formInfoResq;
    }

    public void setFormInfoResq(FormInfoResq formInfoResq) {
        this.formInfoResq = formInfoResq;
    }

    @Override
    public String toString() {
        return "TaskInfoResp{" +
                "id='" + id + '\'' +
                ", requestId='" + requestId + '\'' +
                ", requestNo='" + requestNo + '\'' +
                ", parentId='" + parentId + '\'' +
                ", taskTempId='" + taskTempId + '\'' +
                ", nodeDefId='" + nodeDefId + '\'' +
                ", nodeName='" + nodeName + '\'' +
                ", name='" + name + '\'' +
                ", reporter='" + reporter + '\'' +
                ", reportTime=" + reportTime +
                ", emergency='" + emergency + '\'' +
                ", reportRole='" + reportRole + '\'' +
                ", result='" + result + '\'' +
                ", description='" + description + '\'' +
                ", attachFileId='" + attachFileId + '\'' +
                ", status=" + status +
                ", formInfoResq=" + formInfoResq +
                '}';
    }
}
