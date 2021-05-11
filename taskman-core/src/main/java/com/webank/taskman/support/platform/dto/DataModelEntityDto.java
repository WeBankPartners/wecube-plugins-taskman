package com.webank.taskman.support.platform.dto;

import java.util.ArrayList;
import java.util.List;

public class DataModelEntityDto extends PluginPackageEntityDto {

    private LeafEntityList leafEntityList;

    public class LeafEntityList {

        private List<BoundInterfaceEntityDto> referenceToEntityList = new ArrayList<>();
        private List<BoundInterfaceEntityDto> referenceByEntityList = new ArrayList<>();

        public List<BoundInterfaceEntityDto> getReferenceToEntityList() {
            return referenceToEntityList;
        }

        public void setReferenceToEntityList(List<BoundInterfaceEntityDto> referenceToEntityList) {
            this.referenceToEntityList = referenceToEntityList;
        }

        public List<BoundInterfaceEntityDto> getReferenceByEntityList() {
            return referenceByEntityList;
        }

        public void setReferenceByEntityList(List<BoundInterfaceEntityDto> referenceByEntityList) {
            this.referenceByEntityList = referenceByEntityList;
        }

        public LeafEntityList(List<BoundInterfaceEntityDto> referenceToEntityList,
                List<BoundInterfaceEntityDto> referenceByEntityList) {
            super();
            this.referenceToEntityList = referenceToEntityList;
            this.referenceByEntityList = referenceByEntityList;
        }

        public LeafEntityList() {
            super();
            this.referenceToEntityList = new ArrayList<>();
            this.referenceByEntityList = new ArrayList<>();
        }
    }

    public LeafEntityList getLeafEntityList() {
        return leafEntityList;
    }

    public void setLeafEntityList(LeafEntityList leafEntityList) {
        this.leafEntityList = leafEntityList;
    }

    public DataModelEntityDto() {
        super();
        this.leafEntityList = new LeafEntityList();
    }

    public static class BoundInterfaceEntityDto {
        private String packageName;
        private String entityName;
        private String filterRule;

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

        public BoundInterfaceEntityDto() {
            super();
        }

        public BoundInterfaceEntityDto(String packageName, String entityName, String filterRule) {
            super();
            this.packageName = packageName;
            this.entityName = entityName;
            this.filterRule = filterRule;
        }

        public String getFilterRule() {
            return filterRule;
        }

        public void setFilterRule(String filterRule) {
            this.filterRule = filterRule;
        }

    }
}
