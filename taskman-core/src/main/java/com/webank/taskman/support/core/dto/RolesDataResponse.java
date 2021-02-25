package com.webank.taskman.support.core.dto;

public class RolesDataResponse {

    private String id;
    private String name;
    private String displayName;

    public String getId() {
        return id;
    }

    public RolesDataResponse setId(String id) {
        this.id = id;
        return this;
    }

    public String getName() {
        return name;
    }

    public RolesDataResponse setName(String name) {
        this.name = name;
        return this;
    }

    public String getDisplayName() {
        return displayName;
    }

    public RolesDataResponse setDisplayName(String displayName) {
        this.displayName = displayName;
        return this;
    }
}
