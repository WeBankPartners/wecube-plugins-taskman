package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;

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

    public static long getSerialVersionUID() {
        return serialVersionUID;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getRecordId() {
        return recordId;
    }

    public void setRecordId(String recordId) {
        this.recordId = recordId;
    }

    public String getFormId() {
        return formId;
    }

    public void setFormId(String formId) {
        this.formId = formId;
    }

    public String getItemTempId() {
        return itemTempId;
    }

    public void setItemTempId(String itemTempId) {
        this.itemTempId = itemTempId;
    }

    public String getIsCurrency() {
        return isCurrency;
    }

    public void setIsCurrency(String isCurrency) {
        this.isCurrency = isCurrency;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getValue() {
        return value;
    }

    public void setValue(String value) {
        this.value = value;
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
