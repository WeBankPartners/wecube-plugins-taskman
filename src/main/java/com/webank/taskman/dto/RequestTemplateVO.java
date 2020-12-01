package com.webank.taskman.dto;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import java.util.Date;
@ApiModel(value = "RequestTemplate对象",description = "RequestTemplate")
public class RequestTemplateVO {

    @ApiModelProperty(value = "主键",hidden = true,required = false)
    private String id;

    @ApiModelProperty(value = "模板组编号",required = true,dataType = "String")
    private String groupId;

    @ApiModelProperty(value = "表单模板编号",required = true,dataType = "String")
    private String formTempId;

    @ApiModelProperty(value = "流程编排key",required = true,dataType = "String")
    private String procDefKey;

    @ApiModelProperty(value = "名称",required = true,dataType = "String")
    private String name;

    @ApiModelProperty(value = "创建人",required = true,dataType = "String")
    private String createdBy;

    @ApiModelProperty(value = "更新人",required = true,dataType = "String")
    private String updatedBy;

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

    public String getFormTempId() {
        return formTempId;
    }

    public void setFormTempId(String formTempId) {
        this.formTempId = formTempId;
    }

    public String getProcDefKey() {
        return procDefKey;
    }

    public void setProcDefKey(String procDefKey) {
        this.procDefKey = procDefKey;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getCreatedBy() {
        return createdBy;
    }

    public void setCreatedBy(String createdBy) {
        this.createdBy = createdBy;
    }

    public String getUpdatedBy() {
        return updatedBy;
    }

    public void setUpdatedBy(String updatedBy) {
        this.updatedBy = updatedBy;
    }

    @Override
    public String toString() {
        return "RequestTemplateVO{" +
                "id='" + id + '\'' +
                ", groupId='" + groupId + '\'' +
                ", formTempId='" + formTempId + '\'' +
                ", procDefKey='" + procDefKey + '\'' +
                ", name='" + name + '\'' +
                ", createdBy='" + createdBy + '\'' +
                ", updatedBy='" + updatedBy + '\'' +
                '}';
    }
}
