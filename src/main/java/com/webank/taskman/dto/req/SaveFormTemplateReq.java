package com.webank.taskman.dto.req;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import javax.validation.constraints.NotBlank;

@ApiModel
public class SaveFormTemplateReq {

    private String id;

    @NotBlank(message = "模板id不能为空")
    private String tempId;

    @NotBlank(message = "模板类型不能为空")
    @ApiModelProperty(value = "模板类型(0.请求模板 1.任务模板)",required = true,dataType = "int")
    private String tempType;

    @NotBlank(message = "名称不能为空")
    @ApiModelProperty(value = "",required = true)
    private String name;

    private String description;

    private String style;

    @NotBlank(message = "目标对象级不能为空")
    private String targetEntitys;


    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getTempId() {
        return tempId;
    }

    public void setTempId(String tempId) {
        this.tempId = tempId;
    }

    public String getTempType() {
        return tempType;
    }

    public void setTempType(String tempType) {
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

    public String getTargetEntitys() {
        return targetEntitys;
    }

    public void setTargetEntitys(String targetEntitys) {
        this.targetEntitys = targetEntitys;
    }

    public void setStyle(String style) {
        this.style = style;
    }
}
