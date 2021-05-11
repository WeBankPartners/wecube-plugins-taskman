package com.webank.taskman.support.platform.dto;

import java.util.ArrayList;
import java.util.List;
import java.util.StringJoiner;

public class ProcessDataPreviewDto {

    private String processSessionId;

    private List<GraphNodeDto> entityTreeNodes = new ArrayList<>();

    public List<GraphNodeDto> getEntityTreeNodes() {
        return entityTreeNodes;
    }

    public void setEntityTreeNodes(List<GraphNodeDto> entityTreeNodes) {
        this.entityTreeNodes = entityTreeNodes;
    }

    public static class GraphNodeDto {
        /**
         * packageName : wecmdb entityName : deploy_package dataId :
         * 0045_0000000100 displayName : aDEMO_CORE_APP_mix01.zip id :
         * wecmdb:deploy_package:0045_0000000100 previousIds : [] succeedingIds
         * : ["wecmdb:deploy_package:0045_0000000117"]
         */

        private String packageName;
        private String entityName;
        private String dataId;
        private String displayName;
        private String id;
        private Object entityData;
        private List<String> previousIds;
        private List<String> succeedingIds;

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

        public String getDisplayName() {
            return displayName;
        }

        public void setDisplayName(String displayName) {
            this.displayName = displayName;
        }

        public String getId() {
            return id;
        }

        public void setId(String id) {
            this.id = id;
        }

        public List<String> getPreviousIds() {
            return previousIds;
        }

        public void setPreviousIds(List<String> previousIds) {
            this.previousIds = previousIds;
        }

        public List<String> getSucceedingIds() {
            return succeedingIds;
        }

        public void setSucceedingIds(List<String> succeedingIds) {
            this.succeedingIds = succeedingIds;
        }

        @Override
        public String toString() {
            return new StringJoiner(", ", GraphNodeDto.class.getSimpleName() + "[", "]")
                    .add("packageName='" + packageName + "'").add("entityName='" + entityName + "'")
                    .add("dataId='" + dataId + "'").add("displayName='" + displayName + "'").add("id='" + id + "'")
                    .add("previousIds=" + previousIds).add("succeedingIds=" + succeedingIds).toString();
        }

        public Object getEntityData() {
            return entityData;
        }

        public void setEntityData(Object entityData) {
            this.entityData = entityData;
        }
    }

    public String getProcessSessionId() {
        return processSessionId;
    }

    public void setProcessSessionId(String processSessionId) {
        this.processSessionId = processSessionId;
    }

    @Override
    public String toString() {
        return "ProcessDataPreviewDto{" + "processSessionId='" + processSessionId + '\'' + ", entityTreeNodes="
                + entityTreeNodes + '}';
    }

}
