package com.webank.taskman.dto.req;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import org.apache.commons.lang3.StringUtils;

import java.util.List;

@ApiModel(value = "AddRequestInfoReq",description = "add RequestInfo req")
public class SaveRequestInfoReq {

    @ApiModelProperty(value = "主键",required = false,dataType = "String",position = 100)
    private String id;

    @ApiModelProperty(value = "请求模板id",required = true,dataType = "String",position = 101)
    private String requestTempId;

    @ApiModelProperty(value = "目标对象",required = true,dataType = "String",position = 102)
    private String rootEntity;

    @ApiModelProperty(value = "请求信息名称",required = true,dataType = "String",position = 103)
    private String name;

    @ApiModelProperty(value = "描述",required = false,dataType = "String",position = 104)
    private String description;

    @ApiModelProperty(value = "标签",required = false,dataType = "String",position = 105)
    private String tags;

    @ApiModelProperty(value = "使用角色(多个,分割)",required = false,position = 106)
    private String useRoleName;

    @ApiModelProperty(value = "发布状态(0.未发布 1.已发布 2.已完成)",required = false,position = 107)
    private String status;

    private List<SaveFormItemInfoReq> formItems;

    public String getId() {
        return id;
    }

    public SaveRequestInfoReq setId(String id) {
        this.id = id;
        return this;
    }

    public String getRequestTempId() {
        return requestTempId;
    }

    public void setRequestTempId(String requestTempId) {
        this.requestTempId = requestTempId;
    }

    public String getRootEntity() {
        return rootEntity;
    }

    public void setRootEntity(String rootEntity) {
        this.rootEntity = rootEntity;
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

    public String getUseRoleName() {
        return useRoleName;
    }

    public void setUseRoleName(String useRoleName) {
        this.useRoleName = useRoleName;
    }

    public List<SaveFormItemInfoReq> getFormItems() {
        return formItems;
    }

    public void setFormItems(List<SaveFormItemInfoReq> formItems) {
        this.formItems = formItems;
    }
}
