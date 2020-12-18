package com.webank.taskman.support.core.dto;

import io.swagger.annotations.ApiModel;

import java.util.List;

@ApiModel
public class ProcessDataPreviewDto {

    private String processSessionId;


    private List<GraphNodeDto> entityTreeNodes;

    public List<GraphNodeDto> getEntityTreeNodes() {
        return entityTreeNodes;
    }

    public void setEntityTreeNodes(List<GraphNodeDto> entityTreeNodes) {
        this.entityTreeNodes = entityTreeNodes;
    }

    @ApiModel
    public static class GraphNodeDto {
        /**
         * packageName : wecmdb
         * entityName : deploy_package
         * dataId : 0045_0000000100
         * displayName : aDEMO_CORE_APP_mix01.zip
         * id : wecmdb:deploy_package:0045_0000000100
         * previousIds : []
         * succeedingIds : ["wecmdb:deploy_package:0045_0000000117"]
         */

        private String packageName;
        private String entityName;
        private String dataId;
        private String displayName;
        private String id;
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
    }

    public String getProcessSessionId() {
        return processSessionId;
    }

    public void setProcessSessionId(String processSessionId) {
        this.processSessionId = processSessionId;
    }
}
