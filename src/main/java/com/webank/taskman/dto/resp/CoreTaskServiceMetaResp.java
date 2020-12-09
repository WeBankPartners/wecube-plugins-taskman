package com.webank.taskman.dto.resp;

import io.swagger.annotations.ApiModelProperty;

import java.util.List;

public class CoreTaskServiceMetaResp {

    @ApiModelProperty(value = "",position = 1)
    private String id;

    @ApiModelProperty(value = "",position = 2)
    private String name;

    @ApiModelProperty(value = "",position = 3)
    private String description;

    @ApiModelProperty(value = "",position = 5)
    private String targetEntitys;

    @ApiModelProperty(value = "",position = 6)
    private String inputAttrDef;

    @ApiModelProperty(value = "",position = 7)
    private String outputAttrDef;

    @ApiModelProperty(value = "",position = 8)
    private String otherAttrDef;

    @ApiModelProperty(value = "",position = 9)
    private List<CoreTaskFormItemResp> items;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
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

    public String getTargetEntitys() {
        return targetEntitys;
    }

    public void setTargetEntitys(String targetEntitys) {
        this.targetEntitys = targetEntitys;
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

    public List<CoreTaskFormItemResp> getItems() {
        return items;
    }

    public void setItems(List<CoreTaskFormItemResp> items) {
        this.items = items;
    }
}
