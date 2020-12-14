package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.webank.taskman.base.BaseEntity;

import java.io.Serializable;
import java.util.Date;

public class FormItemInfo  implements Serializable {

    private static final long serialVersionUID = 1L;

    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;

    private String formId;

    private String itemTempId;

    private String name;

    private String value;


    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
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
                ", name='" + name + '\'' +
                ", value='" + value + '\'' +
                '}';
    }
}
