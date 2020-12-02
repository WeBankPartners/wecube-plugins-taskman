package com.webank.taskman.dto.req;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import java.util.List;

@ApiModel
public class AddFormTemplateReq {

    @ApiModelProperty(value = "模板类型(0.请求模板 1.任务模板)",required = true,dataType = "int")
    private Integer tempType;

    @ApiModelProperty(value = "",required = true)
    private String name;

    private String description;

    private String style;

    @ApiModelProperty(value = "表单项模板",required = true,dataType = "int")
    private List<AddFormItemTemplateReq> formItemTemplates;


    public Integer getTempType() {
        return tempType;
    }

    public void setTempType(Integer tempType) {
        this.tempType = tempType;
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

    public String getStyle() {
        return style;
    }

    public void setStyle(String style) {
        this.style = style;
    }

    public List<AddFormItemTemplateReq> getFormItemTemplates() {
        return formItemTemplates;
    }

    public void setFormItemTemplates(List<AddFormItemTemplateReq> formItemTemplates) {
        this.formItemTemplates = formItemTemplates;
    }
}
