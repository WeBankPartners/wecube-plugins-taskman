package com.webank.taskman.dto;

import java.util.ArrayList;
import java.util.List;

public class EntityValueDto {
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
