package com.webank.taskman.dto;

public class TemplateGroupReq {
    private String id;
    private String name;
    private String manageRole;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getManageRole() {
        return manageRole;
    }

    public void setManageRole(String manageRole) {
        this.manageRole = manageRole;
    }


    @Override
    public String toString() {
        return "TemplateGroupReq{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", manageRole='" + manageRole + '\'' +
                '}';
    }
}
