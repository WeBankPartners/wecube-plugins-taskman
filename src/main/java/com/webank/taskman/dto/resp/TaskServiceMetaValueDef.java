package com.webank.taskman.dto.resp;

import io.swagger.annotations.ApiModelProperty;

public class TaskServiceMetaValueDef {

    @ApiModelProperty(value = "",position = 1)
    private String type;
    @ApiModelProperty(value = "",position = 2)
    private String expr;

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public String getExpr() {
        return expr;
    }

    public void setExpr(String expr) {
        this.expr = expr;
    }
}
