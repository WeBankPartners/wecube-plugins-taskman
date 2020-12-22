package com.webank.taskman.support.core.dto;

import javax.annotation.Nonnull;
import javax.annotation.Nullable;
import java.util.ArrayList;
import java.util.List;

public class DynamicEntityValueDto {

    @Nullable
    private String entityDefId;//Entity definition id from platform.
    @Nonnull
    private String packageName;
    @Nonnull
    private String entityName;
    @Nullable
    private String dataId;//Existing data id,such as guid in cmdb.
    @Nonnull
    private String oid;//Equals to dataId once dataId presents,or a temporary assigned.

    private List<String> previousOids = new ArrayList<>();
    private List<String> succeedingOids = new ArrayList<>();
    
    private List<DynamicEntityAttrValueDto> attrValues = new ArrayList<>();

    public DynamicEntityValueDto() {
    }

    public DynamicEntityValueDto(@Nullable String entityDefId, @Nonnull String packageName, @Nonnull String entityName, @Nullable String dataId, @Nonnull String oid) {
        this.entityDefId = entityDefId;
        this.packageName = packageName;
        this.entityName = entityName;
        this.dataId = dataId;
        this.oid = oid;
    }

    public String getEntityDefId() {
        return entityDefId;
    }

    public void setEntityDefId(String entityDefId) {
        this.entityDefId = entityDefId;
    }

    public String getPackageName() {
        return packageName;
    }

    public void setPackageName(String packageName) {
        this.packageName = packageName;
    }

    public String getEntityName() {
        return entityName;
    }

    public void setEntityName(String entityName) {
        this.entityName = entityName;
    }

    public String getDataId() {
        return dataId;
    }

    public void setDataId(String dataId) {
        this.dataId = dataId;
    }

    public String getOid() {
        return oid;
    }

    public void setOid(String oid) {
        this.oid = oid;
    }

    public List<String> getPreviousOids() {
        return previousOids;
    }

    public void setPreviousOids(List<String> previousOids) {
        this.previousOids = previousOids;
    }

    public List<String> getSucceedingOids() {
        return succeedingOids;
    }

    public void setSucceedingOids(List<String> succeedingOids) {
        this.succeedingOids = succeedingOids;
    }

    public List<DynamicEntityAttrValueDto> getAttrValues() {
        return attrValues;
    }

    public void setAttrValues(List<DynamicEntityAttrValueDto> attrValues) {
        this.attrValues = attrValues;
    }

    @Override
    public String toString() {
        return "DynamicEntityValueDto{" +
                "entityDefId='" + entityDefId + '\'' +
                ", packageName='" + packageName + '\'' +
                ", entityName='" + entityName + '\'' +
                ", dataId='" + dataId + '\'' +
                ", oid='" + oid + '\'' +
                ", previousOids=" + previousOids +
                ", succeedingOids=" + succeedingOids +
                ", attrValues=" + attrValues +
                '}';
    }
}
