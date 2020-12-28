package com.webank.taskman.dto.req;


import com.webank.taskman.dto.RequestInfoDTO;
import io.swagger.annotations.ApiModelProperty;

import java.util.Date;

public class QueryRequestInfoReq extends QueryRoleRelationBaseReq {

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
    private Date reportTime;
    @ApiModelProperty(value = "",position = 10)
    private String emergency;
    @ApiModelProperty(value = "",position = 11)
    private String reportRole;
    @ApiModelProperty(value = "",position = 12)
    private String attachFileId;
    @ApiModelProperty(value = "",position = 13)
    private String status;
    @ApiModelProperty(value = "",position = 14)
    private String dueDate;
    @ApiModelProperty(value = "",position = 15)
    private String result;

    public String getId() {
        return id;
    }

    public QueryRequestInfoReq setId(String id) {
        this.id = id;
        return this;
    }

    public String getRequestTempId() {
        return requestTempId;
    }

    public QueryRequestInfoReq setRequestTempId(String requestTempId) {
        this.requestTempId = requestTempId;
        return this;
    }

    public String getRequestTempName() {
        return requestTempName;
    }

    public QueryRequestInfoReq setRequestTempName(String requestTempName) {
        this.requestTempName = requestTempName;
        return this;
    }

    public String getProcInstKey() {
        return procInstKey;
    }

    public QueryRequestInfoReq setProcInstKey(String procInstKey) {
        this.procInstKey = procInstKey;
        return this;
    }

    public String getRootEntity() {
        return rootEntity;
    }

    public QueryRequestInfoReq setRootEntity(String rootEntity) {
        this.rootEntity = rootEntity;
        return this;
    }

    public String getName() {
        return name;
    }

    public QueryRequestInfoReq setName(String name) {
        this.name = name;
        return this;
    }

    public String getDescription() {
        return description;
    }

    public QueryRequestInfoReq setDescription(String description) {
        this.description = description;
        return this;
    }

    public String getReporter() {
        return reporter;
    }

    public QueryRequestInfoReq setReporter(String reporter) {
        this.reporter = reporter;
        return this;
    }

    public Date getReportTime() {
        return reportTime;
    }

    public QueryRequestInfoReq setReportTime(Date reportTime) {
        this.reportTime = reportTime;
        return this;
    }

    public String getEmergency() {
        return emergency;
    }

    public QueryRequestInfoReq setEmergency(String emergency) {
        this.emergency = emergency;
        return this;
    }

    public String getReportRole() {
        return reportRole;
    }

    public QueryRequestInfoReq setReportRole(String reportRole) {
        this.reportRole = reportRole;
        return this;
    }

    public String getAttachFileId() {
        return attachFileId;
    }

    public QueryRequestInfoReq setAttachFileId(String attachFileId) {
        this.attachFileId = attachFileId;
        return this;
    }

    public String getStatus() {
        return status;
    }

    public QueryRequestInfoReq setStatus(String status) {
        this.status = status;
        return this;
    }

    public String getDueDate() {
        return dueDate;
    }

    public QueryRequestInfoReq setDueDate(String dueDate) {
        this.dueDate = dueDate;
        return this;
    }

    public String getResult() {
        return result;
    }

    public QueryRequestInfoReq setResult(String result) {
        this.result = result;
        return this;
    }
}