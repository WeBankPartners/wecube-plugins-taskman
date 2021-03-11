package com.webank.taskman.support.core.dto;

import javax.annotation.Nonnull;
import javax.annotation.Nullable;

import net.logstash.logback.encoder.org.apache.commons.lang3.StringUtils;

import java.util.ArrayList;
import java.util.List;

public class DynamicEntityValueDto {

    @Nullable
    private String entityDefId;// Entity definition id from platform.
    @Nonnull
    private String packageName;
    @Nonnull
    private String entityName;
    @Nullable
    private String dataId;// Existing data id,such as guid in cmdb.
    @Nonnull
    private String oid;// Equals to dataId once dataId presents,or a temporary
                       // assigned.

    private List<String> previousOids = new ArrayList<>();
    private List<String> succeedingOids = new ArrayList<>();

    private List<DynamicEntityAttrValueDto> attrValues = new ArrayList<>();

    public DynamicEntityValueDto() {
    }

    public DynamicEntityValueDto(@Nullable String entityDefId, @Nonnull String packageName, @Nonnull String entityName,
            @Nullable String dataId, @Nonnull String oid) {
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

    public void addPreviousOid(String previousOid) {
        if (StringUtils.isBlank(previousOid)) {
            return;
        }

        if (this.previousOids == null) {
            this.previousOids = new ArrayList<>();
        }

        if (!this.previousOids.contains(previousOid)) {
            this.previousOids.add(previousOid);
        }
    }

    public void addSucceedingOid(String succeedingOid) {
        if (StringUtils.isBlank(succeedingOid)) {
            return;
        }

        if (this.succeedingOids == null) {
            this.succeedingOids = new ArrayList<>();
        }

        if (!this.succeedingOids.contains(succeedingOid)) {
            this.succeedingOids.add(succeedingOid);
        }
    }

    public void addAttrValue(DynamicEntityAttrValueDto attrValue) {
        if (attrValue == null) {
            return;
        }

        if (this.attrValues == null) {
            this.attrValues = new ArrayList<>();
        }

        DynamicEntityAttrValueDto existAttrValueDto = null;
        for (DynamicEntityAttrValueDto a : this.attrValues) {
            if (a.getAttrName().equals(attrValue.getAttrName())) {
                existAttrValueDto = a;
                break;
            }
        }

        if (existAttrValueDto == null) {
            this.attrValues.add(attrValue);
        }
    }

    @Override
    public String toString() {
        return "DynamicEntityValueDto{" + "entityDefId='" + entityDefId + '\'' + ", packageName='" + packageName + '\''
                + ", entityName='" + entityName + '\'' + ", dataId='" + dataId + '\'' + ", oid='" + oid + '\''
                + ", previousOids=" + previousOids + ", succeedingOids=" + succeedingOids + ", attrValues=" + attrValues
                + '}';
    }

    public static class DynamicEntityAttrValueDto {
        private String attrDefId;
        private String attrName;
        private String dataType;
        private Object dataValue;

        public String getAttrDefId() {
            return attrDefId;
        }

        public void setAttrDefId(String attrDefId) {
            this.attrDefId = attrDefId;
        }

        public String getAttrName() {
            return attrName;
        }

        public void setAttrName(String attrName) {
            this.attrName = attrName;
        }

        public String getDataType() {
            return dataType;
        }

        public void setDataType(String dataType) {
            this.dataType = dataType;
        }

        public Object getDataValue() {
            return dataValue;
        }

        public void setDataValue(Object dataValue) {
            this.dataValue = dataValue;
        }

        @Override
        public String toString() {
            return "DynamicEntityAttrValueDto{" + "attrDefId='" + attrDefId + '\'' + ", attrName='" + attrName + '\''
                    + ", dataType='" + dataType + '\'' + ", dataValue=" + dataValue + '}';
        }
    }
}
