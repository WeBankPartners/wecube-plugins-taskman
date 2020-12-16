package com.webank.taskman.dto.resp;

import com.webank.taskman.domain.FormItemTemplate;
import io.swagger.annotations.ApiModelProperty;

import java.util.List;

public class FormTemplateResp {

    @ApiModelProperty(value = "",position = 1)
    private String id;

    @ApiModelProperty(value = "",position = 2)
    private String name;

    @ApiModelProperty(value = "",position = 3)
    private String description;

    @ApiModelProperty(value = "",position = 4)
    private String style;

    @ApiModelProperty(value = "",position = 5)
    private String targetEntitys;

    @ApiModelProperty(value = "",position = 6)
    private String inputAttrDef;

    @ApiModelProperty(value = "",position = 7)
    private String outputAttrDef;

    @ApiModelProperty(value = "",position = 8)
    private String otherAttrDef;

    @ApiModelProperty(value = "",position = 9)
    private List<FormItemTemplate> items;

    public List<FormItemTemplate> getItems() {
        return items;
    }

    public FormTemplateResp setItems(List<FormItemTemplate> items) {
        this.items = items;
        return  this;
    }

    public String getId() {
        return id;
    }

    public FormTemplateResp setId(String id) {
        this.id = id;
        return  this;
    }

    public String getName() {
        return name;
    }

    public FormTemplateResp setName(String name) {
        this.name = name;
        return  this;
    }

    public String getDescription() {
        return description;
    }

    public FormTemplateResp setDescription(String description) {
        this.description = description;
        return  this;
    }

    public String getStyle() {
        return style;
    }

    public FormTemplateResp setStyle(String style) {
        this.style = style;
        return  this;
    }

    public String getTargetEntitys() {
        return targetEntitys;
    }

    public FormTemplateResp setTargetEntitys(String targetEntitys) {
        this.targetEntitys = targetEntitys;
        return  this;
    }

    public String getInputAttrDef() {
        return inputAttrDef;
    }

    public void setInputAttrDef(String inputAttrDef) {
        this.inputAttrDef = inputAttrDef;
    }

    public String getOutputAttrDef() {
        return outputAttrDef;
    }

    public void setOutputAttrDef(String outputAttrDef) {
        this.outputAttrDef = outputAttrDef;
    }

    public String getOtherAttrDef() {
        return otherAttrDef;
    }

    public void setOtherAttrDef(String otherAttrDef) {
        this.otherAttrDef = otherAttrDef;
    }

    @Override
    public String toString() {
        return "FormTemplateResp{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", description='" + description + '\'' +
                ", targetEntitys='" + targetEntitys + '\'' +
                ", style='" + style + '\'' +
                '}';
    }
}
