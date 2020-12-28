package com.webank.taskman.dto.resp;


import com.fasterxml.jackson.annotation.JsonFormat;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import java.util.Date;
import java.util.List;
import java.util.StringJoiner;

@ApiModel
public class RequestInfoResq  {

    @ApiModelProperty(value = "",position = 1)
    private String id;
    @ApiModelProperty(value = "",position = 2)
    private String requestTempId;
    @ApiModelProperty(value = "",position = 3)
    private String requestTempName;
    @ApiModelProperty(value = "",position = 4)
    private String procInstKey;
    @ApiModelProperty(value = "",position = 5)
    private String rootEntity;
    @ApiModelProperty(value = "",position = 6)
    private String name;
    @ApiModelProperty(value = "",position = 7)
    private String description;
    @ApiModelProperty(value = "",position = 8)
    private String reporter;

    @ApiModelProperty(value = "",position = 9)
    @JsonFormat(pattern="yyyy-MM-dd HH:mm:ss",timezone="GMT+8")
    private Date reportTime;
    @ApiModelProperty(value = "",position = 10)
    private String emergency;
    @ApiModelProperty(value = "",position = 11)
    private String reportRole;
    @ApiModelProperty(value = "",position = 13)
    private String status;
    @ApiModelProperty(value = "",position = 14)
    private String dueDate;
    @ApiModelProperty(value = "",position = 15)
    private String result;
    @ApiModelProperty(value = "",position = 16)
    private String attachFileId;
    @ApiModelProperty(value = "",position = 17)
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

    public String getProcInstKey() {
        return procInstKey;
    }

    public RequestInfoResq setProcInstKey(String procInstKey) {
        this.procInstKey = procInstKey;
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
                .toString();
    }
}