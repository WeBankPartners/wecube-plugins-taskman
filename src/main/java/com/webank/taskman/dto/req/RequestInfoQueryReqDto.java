package com.webank.taskman.dto.req;

public class RequestInfoQueryReqDto extends RoleRelationBaseQueryReqDto {

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

    public RequestInfoQueryReqDto setId(String id) {
        this.id = id;
        return this;
    }

    public String getRequestTempId() {
        return requestTempId;
    }

    public RequestInfoQueryReqDto setRequestTempId(String requestTempId) {
        this.requestTempId = requestTempId;
        return this;
    }

    public String getRequestTempName() {
        return requestTempName;
    }

    public RequestInfoQueryReqDto setRequestTempName(String requestTempName) {
        this.requestTempName = requestTempName;
        return this;
    }

    public String getName() {
        return name;
    }

    public RequestInfoQueryReqDto setName(String name) {
        this.name = name;
        return this;
    }

    public String getDescription() {
        return description;
    }

    public RequestInfoQueryReqDto setDescription(String description) {
        this.description = description;
        return this;
    }

    public String getReporter() {
        return reporter;
    }

    public RequestInfoQueryReqDto setReporter(String reporter) {
        this.reporter = reporter;
        return this;
    }

    public String getReportTimeBegin() {
        return reportTimeBegin;
    }

    public RequestInfoQueryReqDto setReportTimeBegin(String reportTimeBegin) {
        this.reportTimeBegin = reportTimeBegin;
        return this;
    }

    public String getReportTimeEnd() {
        return reportTimeEnd;
    }

    public RequestInfoQueryReqDto setReportTimeEnd(String reportTimeEnd) {
        this.reportTimeEnd = reportTimeEnd;
        return this;
    }

    public String getEmergency() {
        return emergency;
    }

    public RequestInfoQueryReqDto setEmergency(String emergency) {
        this.emergency = emergency;
        return this;
    }

    public String getStatus() {
        return status;
    }

    public RequestInfoQueryReqDto setStatus(String status) {
        this.status = status;
        return this;
    }
}