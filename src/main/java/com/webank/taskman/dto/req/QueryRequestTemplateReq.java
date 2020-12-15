package com.webank.taskman.dto.req;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.webank.taskman.dto.RoleDTO;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import org.apache.commons.lang3.StringUtils;

import javax.validation.constraints.NotBlank;
import java.util.List;

@ApiModel(value = "AddRequestTemplateReq",description = "add RequestTemplate req")
public class QueryRequestTemplateReq extends QueryRoleRelationBaseReq{

    @ApiModelProperty(value = "主键",required = false,dataType = "String",position = 100)
    private String id;

    @ApiModelProperty(value = "模板组编号",required = false,dataType = "String",position = 101)
    private String requestTempGroup;

    @ApiModelProperty(value = "流程编排Id",required = false,dataType = "String",position = 102)
    private String procDefId;

    @ApiModelProperty(value = "流程编排key",required = false,dataType = "String",position = 103)
    private String procDefKey;

    @ApiModelProperty(value = "流程编排名称",required = false,dataType = "String",position = 104)
    private String procDefName;

    @ApiModelProperty(value = "请求模板名称",required = false,dataType = "String",position = 105)
    private String name;

    @ApiModelProperty(value = "标签",required = false,dataType = "String",position = 107)
    private String tags;

    @ApiModelProperty(value = "发布状态",required = false,dataType = "Integer",position = 108)
    private Integer status;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getRequestTempGroup() {
        return requestTempGroup;
    }

    public void setRequestTempGroup(String requestTempGroup) {
        this.requestTempGroup = requestTempGroup;
    }

    public String getProcDefId() {
        return procDefId;
    }

    public void setProcDefId(String procDefId) {
        this.procDefId = procDefId;
    }

    public String getProcDefKey() {
        return procDefKey;
    }

    public void setProcDefKey(String procDefKey) {
        this.procDefKey = procDefKey;
    }

    public String getProcDefName() {
        return procDefName;
    }

    public void setProcDefName(String procDefName) {
        this.procDefName = procDefName;
    }


    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }



    public Integer getStatus() {
        return status;
    }

    public void setStatus(Integer status) {
        this.status = status;
    }

    public String getTags() {
        return tags;
    }

    public void setTags(String tags) {
        this.tags = tags;
    }



}
