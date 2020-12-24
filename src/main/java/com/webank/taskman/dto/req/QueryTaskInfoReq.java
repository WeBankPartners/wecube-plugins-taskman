package com.webank.taskman.dto.req;


import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

@ApiModel
public class QueryTaskInfoReq extends QueryRoleRelationBaseReq{

    @ApiModelProperty(value = "任务id",position = 1)
    private String id;
    @ApiModelProperty(value = "请求ID",position = 1)
    private String requestId;
    @ApiModelProperty(value = "流程节点名称",position = 6)
    private String nodeName;
    @ApiModelProperty(value = "任务名称",position = 7)
    private String name;
    @ApiModelProperty(value = "处理人",position = 9)
    private String reporter;
    @ApiModelProperty(value = "紧急程度",position = 11)
    private String emergency;
    @ApiModelProperty(value = "任务状态",position = 13)
    private Integer status;

    public String getId() {
        return id;
    }

    public QueryTaskInfoReq setId(String id) {
        this.id = id;
        return this;
    }

    public String getRequestId() {
        return requestId;
    }

    public QueryTaskInfoReq setRequestId(String requestId) {
        this.requestId = requestId;
        return this;
    }

    public String getNodeName() {
        return nodeName;
    }

    public QueryTaskInfoReq setNodeName(String nodeName) {
        this.nodeName = nodeName;
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

    public String getEmergency() {
        return emergency;
    }

    public QueryTaskInfoReq setEmergency(String emergency) {
        this.emergency = emergency;
        return this;
    }

    public Integer getStatus() {
        return status;
    }

    public QueryTaskInfoReq setStatus(Integer status) {
        this.status = status;
        return this;
    }
}
