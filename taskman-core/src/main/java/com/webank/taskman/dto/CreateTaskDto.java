package com.webank.taskman.dto;

import java.util.ArrayList;
import java.util.List;

public class CreateTaskDto {
    private String id;
    private String requestTempId;
    private String emergency;
    private String name;
    private String description;
    private String rootEntityPackage;
    private String rootEntityName;
    private String rootEntity;
    private List<EntityValueDto> entities = new ArrayList<>();

    public String getId() {
        return id;
    }

    public CreateTaskDto setId(String id) {
        this.id = id;
        return this;
    }

    public String getRequestTempId() {
        return requestTempId;
    }

    public CreateTaskDto setRequestTempId(String requestTempId) {
        this.requestTempId = requestTempId;
        return this;
    }

    public String getEmergency() {
        return emergency;
    }

    public CreateTaskDto setEmergency(String emergency) {
        this.emergency = emergency;
        return this;
    }

    public String getName() {
        return name;
    }

    public CreateTaskDto setName(String name) {
        this.name = name;
        return this;
    }

    public String getDescription() {
        return description;
    }

    public CreateTaskDto setDescription(String description) {
        this.description = description;
        return this;
    }

    public String getRootEntity() {
        return rootEntity;
    }

    public CreateTaskDto setRootEntity(String rootEntity) {
        this.rootEntity = rootEntity;
        return this;
    }

    public List<EntityValueDto> getEntities() {
        return entities;
    }

    public CreateTaskDto setEntities(List<EntityValueDto> entities) {
        this.entities = entities;
        return this;
    }

    public String getRootEntityPackage() {
        return rootEntityPackage;
    }

    public void setRootEntityPackage(String rootEntityPackage) {
        this.rootEntityPackage = rootEntityPackage;
    }

    public String getRootEntityName() {
        return rootEntityName;
    }

    public void setRootEntityName(String rootEntityName) {
        this.rootEntityName = rootEntityName;
    }

}
