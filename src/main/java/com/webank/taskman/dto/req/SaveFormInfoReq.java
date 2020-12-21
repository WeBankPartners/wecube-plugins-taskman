package com.webank.taskman.dto.req;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import java.util.List;

@ApiModel
public class SaveFormInfoReq {


    @ApiModelProperty(value = "表单名称",required = true)
    private String name;

    private String description;

    @ApiModelProperty(value = "输入",position = 8)
    private String inputAttrDef;
    @ApiModelProperty(value = "输出参数",position = 9)
    private String outputAttrDef;
    @ApiModelProperty(value = "其他参数",position = 10)
    private String otherAttrDef;

    private List<SaveFormItemInfoReq> formItems;

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

    public List<SaveFormItemInfoReq> getFormItems() {
        return formItems;
    }

    public void setFormItems(List<SaveFormItemInfoReq> formItems) {
        this.formItems = formItems;
    }
}
