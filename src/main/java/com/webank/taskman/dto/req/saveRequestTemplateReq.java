package com.webank.taskman.dto.req;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

@ApiModel(value = "AddRequestTemplateReq",description = "add RequestTemplate req")
public class saveRequestTemplateReq {

    @ApiModelProperty(value = "主键",required = false,dataType = "String")
    private String id;

    @ApiModelProperty(value = "模板组编号",required = true,dataType = "String")
    private String groupId;

    @ApiModelProperty(value = "流程编排Id",required = true,dataType = "String")
    private String procDefId;
    @ApiModelProperty(value = "流程编排key",required = true,dataType = "String")
    private String procDefKey;
    @ApiModelProperty(value = "流程编排名称",required = true,dataType = "String")
    private String procDefName;

    @ApiModelProperty(value = "请求模板名称",required = true,dataType = "String")
    private String name;
    @ApiModelProperty(value = "描述",required = true,dataType = "String")
    private String description;

    @ApiModelProperty(value = "标签",required = true,dataType = "String")
    private String tags;


    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getGroupId() {
        return groupId;
    }

    public void setGroupId(String groupId) {
        this.groupId = groupId;
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

    public String getTags() {
        return tags;
    }

    public void setTags(String tags) {
        this.tags = tags;
    }

}
