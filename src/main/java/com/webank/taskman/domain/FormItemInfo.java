package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.fasterxml.jackson.annotation.JsonIgnore;
import org.springframework.util.StringUtils;

import java.io.Serializable;

public class FormItemInfo  implements Serializable {

    private static final long serialVersionUID = 1L;

    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;

    private String recordId;

    private String formId;

    private String itemTempId;

    private String isCurrency;

    private String name;

    private String value;

    public String getId() {
        return id;
    }

    public FormItemInfo setId(String id) {
        this.id = id;
        return this;
    }

    public String getRecordId() {
        return recordId;
    }

    public FormItemInfo setRecordId(String recordId) {
        this.recordId = recordId;
        return this;
    }

    public String getFormId() {
        return formId;
    }

    public FormItemInfo setFormId(String formId) {
        this.formId = formId;
        return this;
    }

    public String getItemTempId() {
        return itemTempId;
    }

    public FormItemInfo setItemTempId(String itemTempId) {
        this.itemTempId = itemTempId;
        return this;
    }

    public String getIsCurrency() {
        return isCurrency;
    }

    public FormItemInfo setIsCurrency(String isCurrency) {
        this.isCurrency = isCurrency;
        return this;
    }

    public String getName() {
        return name;
    }

    public FormItemInfo setName(String name) {
        this.name = name;
        return this;
    }

    public String getValue() {
        return value;
    }

    public FormItemInfo setValue(String value) {
        this.value = value;
        return this;
    }

    @JsonIgnore
    public LambdaQueryWrapper getLambdaQueryWrapper() {
        return new LambdaQueryWrapper<FormItemInfo>()
                .eq(!StringUtils.isEmpty(id), FormItemInfo::getId, id)
                .eq(!StringUtils.isEmpty(recordId), FormItemInfo::getRecordId, recordId)
                .eq(!StringUtils.isEmpty(formId), FormItemInfo::getFormId, formId)
                .eq(!StringUtils.isEmpty(itemTempId), FormItemInfo::getItemTempId, itemTempId)
                .eq(!StringUtils.isEmpty(isCurrency), FormItemInfo::getIsCurrency, isCurrency)
                .eq(!StringUtils.isEmpty(name), FormItemInfo::getName, name)
                .eq(!StringUtils.isEmpty(value), FormItemInfo::getValue, value);
    }

    @Override
    public String toString() {
        return "FormItemInfo{" +
                "id='" + id + '\'' +
                ", formId='" + formId + '\'' +
                ", itemTempId='" + itemTempId + '\'' +
                ", isCurrency='" + isCurrency + '\'' +
                ", name='" + name + '\'' +
                ", value='" + value + '\'' +
                '}';
    }
}
