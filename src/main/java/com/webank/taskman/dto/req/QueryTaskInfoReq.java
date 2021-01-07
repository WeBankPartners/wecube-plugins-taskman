package com.webank.taskman.dto.req;


import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import java.util.Date;
import java.util.StringJoiner;

@ApiModel
public class QueryTaskInfoReq extends QueryRoleRelationBaseReq{


    @ApiModelProperty(value = "",position = 1)
    private String id;
    @ApiModelProperty(value = "",position = 4)
    private String procInstId;
    @ApiModelProperty(value = "",position = 5)
    private String nodeDefId;
    @ApiModelProperty(value = "",position = 7)
    private String name;
    @ApiModelProperty(value = "",position = 9)
    private String reporter;
    @ApiModelProperty(value = "",position = 10)
    private Date reportTime;
    @ApiModelProperty(value = "",position = 11)
    private String emergency;
    @ApiModelProperty(value = "",position = 12)
    private String status;
    @ApiModelProperty(value = "",position = 13)
    private String result;

    private String isMy;



    public String getId() {
        return id;
    }

    public QueryTaskInfoReq setId(String id) {
        this.id = id;
        return this;
    }

    public String getProcInstId() {
        return procInstId;
    }

    public QueryTaskInfoReq setProcInstId(String procInstId) {
        this.procInstId = procInstId;
        return this;
    }

    public String getNodeDefId() {
        return nodeDefId;
    }

    public QueryTaskInfoReq setNodeDefId(String nodeDefId) {
        this.nodeDefId = nodeDefId;
        return this;
    }

    public String getName() {
        return name;
    }

    public QueryTaskInfoReq setName(String name) {
        this.name = name;
        return this;
    }

    public String getReporter() {
        return reporter;
    }

    public QueryTaskInfoReq setReporter(String reporter) {
        this.reporter = reporter;
        return this;
    }

    public Date getReportTime() {
        return reportTime;
    }

    public QueryTaskInfoReq setReportTime(Date reportTime) {
        this.reportTime = reportTime;
        return this;
    }

    public String getEmergency() {
        return emergency;
    }

    public QueryTaskInfoReq setEmergency(String emergency) {
        this.emergency = emergency;
        return this;
    }

    public String getStatus() {
        return status;
    }

    public QueryTaskInfoReq setStatus(String status) {
        this.status = status;
        return this;
    }

    public String getResult() {
        return result;
    }

    public QueryTaskInfoReq setResult(String result) {
        this.result = result;
        return this;
    }

    public String getIsMy() {
        return isMy;
    }

    public QueryTaskInfoReq setIsMy(String isMy) {
        this.isMy = isMy;
        return this;
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", QueryTaskInfoReq.class.getSimpleName() + "[", "]")
                .add("id='" + id + "'")
                .add("procInstId='" + procInstId + "'")
                .add("nodeDefId='" + nodeDefId + "'")
                .add("name='" + name + "'")
                .add("reporter='" + reporter + "'")
                .add("reportTime=" + reportTime)
                .add("emergency='" + emergency + "'")
                .add("status='" + status + "'")
                .add("result='" + result + "'")
                .toString();
    }
}
