package com.webank.taskman.dto;

import io.swagger.annotations.ApiModelProperty;

import java.util.ArrayList;
import java.util.List;

public class CreateTaskDto {
    @ApiModelProperty(value = "",position = 1,hidden = true)
    private String id;
    @ApiModelProperty(value = "",position = 1)
    private String requestTempId;
    @ApiModelProperty(value = "",position = 2)
    private String emergency;
    @ApiModelProperty(value = "",position = 3)
    private String name;
    @ApiModelProperty(value = "",position = 4)
    private String description;
    @ApiModelProperty(value = "",position = 5)
    private String rootEntity;
    @ApiModelProperty(value = "",position = 6)
    private List<EntityValueDto> entitys = new ArrayList<>();


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

    public List<EntityValueDto> getEntitys() {
        return entitys;
    }

    public CreateTaskDto setEntitys(List<EntityValueDto> entitys) {
        this.entitys = entitys;
        return this;
    }

    public static class EntityValueDto {
        @ApiModelProperty(value = "",position = 1)
        private String nodeId;
        @ApiModelProperty(value = "",position = 2)
        private String nodeDefId;
        @ApiModelProperty(value = "",position = 3)
        private String dataId;
        @ApiModelProperty(value = "",position = 4)
        private String packageName;
        @ApiModelProperty(value = "",position = 5)
        private String entityName;
        @ApiModelProperty(value = "",position = 6)
        private List<String> previousOids = new ArrayList<>();
        @ApiModelProperty(value = "",position = 7)
        private List<String> succeedingOids = new ArrayList<>();
        @ApiModelProperty(value = "",position = 6)
        private List<EntityAttrValueDto> attrValues = new ArrayList<>();

        public String getNodeId() {
            return nodeId;
        }

        public EntityValueDto setNodeId(String nodeId) {
            this.nodeId = nodeId;
            return this;
        }

        public String getNodeDefId() {
            return nodeDefId;
        }

        public EntityValueDto setNodeDefId(String nodeDefId) {
            this.nodeDefId = nodeDefId;
            return this;
        }

        public String getDataId() {
            return dataId;
        }

        public EntityValueDto setDataId(String dataId) {
            this.dataId = dataId;
            return this;
        }

        public String getPackageName() {
            return packageName;
        }

        public EntityValueDto setPackageName(String packageName) {
            this.packageName = packageName;
            return this;
        }

        public String getEntityName() {
            return entityName;
        }

        public EntityValueDto setEntityName(String entityName) {
            this.entityName = entityName;
            return this;
        }

        public List<String> getPreviousOids() {
            return previousOids;
        }

        public EntityValueDto setPreviousOids(List<String> previousOids) {
            this.previousOids = previousOids;
            return this;
        }

        public List<String> getSucceedingOids() {
            return succeedingOids;
        }

        public EntityValueDto setSucceedingOids(List<String> succeedingOids) {
            this.succeedingOids = succeedingOids;
            return this;
        }

        public List<EntityAttrValueDto> getAttrValues() {
            return attrValues;
        }

        public EntityValueDto setAttrValues(List<EntityAttrValueDto> attrValues) {
            this.attrValues = attrValues;
            return this;
        }
    }

    public static class EntityAttrValueDto {

        @ApiModelProperty(value = "",position = 1)
        private String itemTempId;
        @ApiModelProperty(value = "",position = 2)
        private String attrDefId;
        @ApiModelProperty(value = "",position = 3)
        private String dataType;
        @ApiModelProperty(value = "",position = 4)
        private String name;
        @ApiModelProperty(value = "",position = 5)
        private Object dataValue;

        public String getItemTempId() {
            return itemTempId;
        }

        public EntityAttrValueDto setItemTempId(String itemTempId) {
            this.itemTempId = itemTempId;
            return this;
        }

        public String getAttrDefId() {
            return attrDefId;
        }

        public EntityAttrValueDto setAttrDefId(String attrDefId) {
            this.attrDefId = attrDefId;
            return this;
        }

        public String getDataType() {
            return dataType;
        }

        public EntityAttrValueDto setDataType(String dataType) {
            this.dataType = dataType;
            return this;
        }

        public String getName() {
            return name;
        }

        public EntityAttrValueDto setName(String name) {
            this.name = name;
            return this;
        }

        public Object getDataValue() {
            return dataValue;
        }

        public EntityAttrValueDto setDataValue(Object dataValue) {
            this.dataValue = dataValue;
            return this;
        }
    }


}
