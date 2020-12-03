package com.webank.taskman.dto.resp;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

@ApiModel
public class TaskTemplateResp {

    @ApiModelProperty(value = "",required = true)
    private String procDefId;
    @ApiModelProperty(value = "",required = true)
    private String procDefKey;
    @ApiModelProperty(value = "",required = true)
    private String procDefName;
    @ApiModelProperty(value = "",required = true)
    private String procNode;
    @ApiModelProperty(value = "",required = true)
    private String name;

    private String description;


    public String getProcDefId() {
        return procDefId;
    }

    public void setProcDefId(String procDefId) {
        this.procDefId = procDefId;
    }

    public String getProcDefKey() {
        return procDefKey;
    }

    public void setProcDefKey(String procDefKey) {
        this.procDefKey = procDefKey;
    }

    public String getProcDefName() {
        return procDefName;
    }

    public void setProcDefName(String procDefName) {
        this.procDefName = procDefName;
    }

    public String getProcNode() {
        return procNode;
    }

    public void setProcNode(String procNode) {
        this.procNode = procNode;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

}
