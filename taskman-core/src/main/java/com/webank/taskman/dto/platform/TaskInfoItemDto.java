package com.webank.taskman.dto.platform;

import java.util.List;

public class TaskInfoItemDto {
    private String itemId;
    private String dataId;
    private String key;
    private List<String> val;

    public String getItemId() {
        return itemId;
    }

    public TaskInfoItemDto setItemId(String itemId) {
        this.itemId = itemId;
        return this;
    }

    public String getDataId() {
        return dataId;
    }

    public TaskInfoItemDto setDataId(String dataId) {
        this.dataId = dataId;
        return this;
    }

    public String getKey() {
        return key;
    }

    public TaskInfoItemDto setKey(String key) {
        this.key = key;
        return this;
    }

    public List<String> getVal() {
        return val;
    }

    public TaskInfoItemDto setVal(List<String> val) {
        this.val = val;
        return this;
    }

}
