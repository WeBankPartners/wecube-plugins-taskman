package com.webank.taskman.dto.req;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

@ApiModel(value = "AddRequestTemplateReq",description = "add RequestTemplate req")
public class AddRequestTemplateReq {
    @ApiModelProperty(value = "模板组编号",required = true,dataType = "String")
    private String groupId;

    @ApiModelProperty(value = "表单模板编号",required = true,dataType = "String")
    private String formTempId;

    @ApiModelProperty(value = "名称",required = true,dataType = "String")
    private String name;

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

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }
}
