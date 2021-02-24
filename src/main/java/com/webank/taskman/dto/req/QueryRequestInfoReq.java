package com.webank.taskman.dto.req;

public class QueryRequestInfoReq extends QueryRoleRelationBaseReq {

    private String id;
    private String requestTempId;
    private String requestTempName;
    private String name;
    private String description;
    private String reporter;
    private String emergency;
    private String status;

    private String reportTimeBegin;
    private String reportTimeEnd;

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

    public String getReportTimeBegin() {
        return reportTimeBegin;
    }

    public QueryRequestInfoReq setReportTimeBegin(String reportTimeBegin) {
        this.reportTimeBegin = reportTimeBegin;
        return this;
    }

    public String getReportTimeEnd() {
        return reportTimeEnd;
    }

    public QueryRequestInfoReq setReportTimeEnd(String reportTimeEnd) {
        this.reportTimeEnd = reportTimeEnd;
        return this;
    }

    public String getEmergency() {
        return emergency;
    }

    public QueryRequestInfoReq setEmergency(String emergency) {
        this.emergency = emergency;
        return this;
    }

    public String getStatus() {
        return status;
    }

    public QueryRequestInfoReq setStatus(String status) {
        this.status = status;
        return this;
    }
}