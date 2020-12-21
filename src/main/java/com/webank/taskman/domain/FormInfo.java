package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.webank.taskman.base.BaseEntity;

import java.io.Serializable;

public class FormInfo extends BaseEntity implements Serializable {

    private static final long serialVersionUID = 1L;


    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;

    private String recordId;

    private String formTemplateId;

    private String name;

    private Integer type;

    public String getId() {
        return id;
    }

    public FormInfo setId(String id) {
        this.id = id;
        return  this;
    }

    public String getRecordId() {
        return recordId;
    }

    public FormInfo setRecordId(String recordId) {
        this.recordId = recordId;
        return  this;
    }

    public String getFormTemplateId() {
        return formTemplateId;
    }

    public FormInfo setFormTemplateId(String formTemplateId) {
        this.formTemplateId = formTemplateId;
        return  this;
    }

    public String getName() {
        return name;
    }

    public FormInfo setName(String name) {
        this.name = name;
        return  this;
    }

    public Integer getType() {
        return type;
    }

    public FormInfo setType(Integer type) {
        this.type = type;
        return  this;
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
