package com.webank.taskman.dto.req;

import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.RoleDTO;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import javax.validation.constraints.NotBlank;
import javax.validation.constraints.NotEmpty;
import java.util.LinkedList;
import java.util.List;

@ApiModel
public class SaveFormTemplateReq {

    @ApiModelProperty(value = "",position = 1)
    private String id;

    @NotBlank(message = "模板id不能为空")
    @ApiModelProperty(value = "",position = 2)
    private String tempId;

    @ApiModelProperty(value = "模板类型(0.请求模板 1.任务模板)",required = true,dataType = "int",position = 3)
    private Integer tempType=0;

    @NotBlank(message = "名称不能为空")
    @ApiModelProperty(value = "",required = true,position = 4)
    private String name;

    @ApiModelProperty(value = "",position = 5)
    private String description;

    @ApiModelProperty(value = "",position = 6)
    private String style;

    @ApiModelProperty(value = "目标对象",position = 7)
    private String targetEntitys;

    @ApiModelProperty(value = "输入",position = 8)
    private String inputAttrDef;
    @ApiModelProperty(value = "输出参数",position = 9)
    private String outputAttrDef;
    @ApiModelProperty(value = "其他参数",position = 10)
    private String otherAttrDef;

    @ApiModelProperty(value = "使用角色集",position = 12)
    private List<RoleDTO> useRole;

    @ApiModelProperty(value = "表单项",position = 13)
    private List<SaveFormItemTemplateReq> formItems = new LinkedList<>();


    public SaveFormTemplateReq() {
    }

    public SaveFormTemplateReq(String tempId, Integer tempType) {
        this.tempId = tempId;
        this.tempType = tempType;
    }

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

    public List<RoleDTO> getUseRole() {
        return useRole;
    }

    public void setUseRole(List<RoleDTO> useRole) {
        this.useRole = useRole;
    }

    public List<SaveFormItemTemplateReq> getFormItems() {
        return formItems;
    }

    public void setFormItems(List<SaveFormItemTemplateReq> formItems) {
        this.formItems = formItems;
    }
}
