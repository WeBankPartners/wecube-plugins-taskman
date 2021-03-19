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

    public static class EntityValueDto {
        private String nodeId;
        private String nodeDefId;
        private String oid;
        private String dataId;
        private String packageName;
        private String entityName;
        private List<String> previousOids = new ArrayList<>();
        private List<String> succeedingOids = new ArrayList<>();
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

        public String getOid() {
            return oid;
        }

        public void setOid(String oid) {
            this.oid = oid;
        }

        @Override
        public String toString() {
            StringBuilder builder = new StringBuilder();
            builder.append("EntityValueDto [nodeId=");
            builder.append(nodeId);
            builder.append(", nodeDefId=");
            builder.append(nodeDefId);
            builder.append(", oid=");
            builder.append(oid);
            builder.append(", dataId=");
            builder.append(dataId);
            builder.append(", packageName=");
            builder.append(packageName);
            builder.append(", entityName=");
            builder.append(entityName);
            builder.append(", previousOids=");
            builder.append(previousOids);
            builder.append(", succeedingOids=");
            builder.append(succeedingOids);
            builder.append(", attrValues=");
            builder.append(attrValues);
            builder.append("]");
            return builder.toString();
        }
        
        

    }

    public static class EntityAttrValueDto {

        private String itemTempId;
        private String attrDefId;
        private String dataType;
        private String name;
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

        @Override
        public String toString() {
            StringBuilder builder = new StringBuilder();
            builder.append("EntityAttrValueDto [itemTempId=");
            builder.append(itemTempId);
            builder.append(", attrDefId=");
            builder.append(attrDefId);
            builder.append(", dataType=");
            builder.append(dataType);
            builder.append(", name=");
            builder.append(name);
            builder.append(", dataValue=");
            builder.append(dataValue);
            builder.append("]");
            return builder.toString();
        }
        
        
    }

}
