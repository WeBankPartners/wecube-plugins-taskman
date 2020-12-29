package com.webank.taskman.dto.req;


import com.webank.taskman.dto.RequestInfoDTO;
import io.swagger.annotations.ApiModelProperty;
import org.springframework.util.StringUtils;

import java.util.Date;

public class QueryRequestInfoReq extends QueryRoleRelationBaseReq {

    @ApiModelProperty(value = "",position = 1)
    private String id;
    @ApiModelProperty(value = "",position = 2)
    private String requestTempId;
    @ApiModelProperty(value = "",position = 3)
    private String requestTempName;
    @ApiModelProperty(value = "",position = 6)
    private String name;
    @ApiModelProperty(value = "",position = 7)
    private String description;
    @ApiModelProperty(value = "",position = 8)
    private String reporter;
    @ApiModelProperty(value = "",position = 9)
    private String reportTime;
    @ApiModelProperty(value = "",position = 12)
    private String emergency;
    @ApiModelProperty(value = "",position = 15)
    private String status;

    @ApiModelProperty(value = "",hidden = true)
    private String reportTimeBegin;
    @ApiModelProperty(value = "",hidden = true)
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

    public String getReportTime() {
        return reportTime;
    }

    public QueryRequestInfoReq setReportTime(String reportTime) {
        this.reportTime = reportTime;
        if(!StringUtils.isEmpty(reportTime)){
            String [] reportTimes = reportTime.split(",");
            this.reportTimeBegin = reportTimes[0];
            this.reportTimeEnd = reportTimes.length > 1 ? reportTimes[1] :"";
        }
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