package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.webank.taskman.base.BaseEntity;
import org.springframework.util.StringUtils;

import java.io.Serializable;

public class FormInfo extends BaseEntity implements Serializable {

    private static final long serialVersionUID = 1L;


    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;
    private String recordId;
    private String formTemplateId;
    private String name;
    private Integer type;
    private String formType;

    public String getId() {
        return id;
    }

    public FormInfo setId(String id) {
        this.id = id;
        return this;
    }

    public String getRecordId() {
        return recordId;
    }

    public FormInfo setRecordId(String recordId) {
        this.recordId = recordId;
        return this;
    }

    public String getFormTemplateId() {
        return formTemplateId;
    }

    public FormInfo setFormTemplateId(String formTemplateId) {
        this.formTemplateId = formTemplateId;
        return this;
    }

    public String getName() {
        return name;
    }

    public FormInfo setName(String name) {
        this.name = name;
        return this;
    }

    public Integer getType() {
        return type;
    }

    public FormInfo setType(Integer type) {
        this.type = type;
        return this;
    }

    public String getFormType() {
        return formType;
    }

    public FormInfo setFormType(String formType) {
        this.formType = formType;
        return this;
    }

    public FormInfo() {
    }

    public FormInfo(String id, String recordId, String formTemplateId, String name, Integer type) {
        this.id = id;
        this.recordId = recordId;
        this.formTemplateId = formTemplateId;
        this.name = name;
        this.type = type;
    }

    @JsonIgnore
    public LambdaQueryWrapper<FormInfo> getLambdaQueryWrapper() {
        return new LambdaQueryWrapper<FormInfo>()
            .eq(!StringUtils.isEmpty(id), FormInfo::getId, id)
            .eq(!StringUtils.isEmpty(recordId), FormInfo::getRecordId, recordId)
            .eq(!StringUtils.isEmpty(formTemplateId), FormInfo::getFormTemplateId, formTemplateId)
            .eq(!StringUtils.isEmpty(name), FormInfo::getName, name)
            .eq(!StringUtils.isEmpty(type), FormInfo::getType, type);
    }



    @Override
    public String toString() {
        return "FormInfo{" +
                "id='" + id + '\'' +
                ", recordId='" + recordId + '\'' +
                ", formTemplateId='" + formTemplateId + '\'' +
                ", name='" + name + '\'' +
                ", type=" + type +
                '}';
    }
}
