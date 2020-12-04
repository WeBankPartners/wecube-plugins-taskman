package com.webank.taskman.dto.req;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import javax.validation.constraints.NotBlank;

@ApiModel
public class SaveTaskTemplateReq {

    @ApiModelProperty(value = "任务模板ID",dataType = "String")
    private String id;

    @NotBlank(message = "流程编排id不能为空")
    @ApiModelProperty(value = "流程编排id",dataType = "String")
    private String procDefId;

    @NotBlank(message = "流程编排key不能为空")
    @ApiModelProperty(value = "流程编排key",dataType = "String")
    private String procDefKey;

    @NotBlank(message = "流程编排名称不能为空")
    @ApiModelProperty(value = "流程编排名称",dataType = "String")
    private String procDefName;

    @NotBlank(message = "流程节点不能为空")
    @ApiModelProperty(value = "流程节点",dataType = "String")
    private String procNode;

    @NotBlank(message = "名称不能为空")
    @ApiModelProperty(value = "名称",dataType = "String")
    private String name;

    @ApiModelProperty(value = "描述",dataType = "String")
    private String description;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

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

    @Override
    public String toString() {
        return "SaveTaskTemplateReq{" +
                "id='" + id + '\'' +
                ", procDefId='" + procDefId + '\'' +
                ", procDefKey='" + procDefKey + '\'' +
                ", procDefName='" + procDefName + '\'' +
                ", procNode='" + procNode + '\'' +
                ", name='" + name + '\'' +
                ", description='" + description + '\'' +
                '}';
    }
}
