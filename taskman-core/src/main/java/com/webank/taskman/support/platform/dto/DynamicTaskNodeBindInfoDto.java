package com.webank.taskman.support.platform.dto;

import java.util.ArrayList;
import java.util.List;
import java.util.StringJoiner;

public class DynamicTaskNodeBindInfoDto {
    private String nodeId;//Physical node id from BPMN.
    private String nodeDefId;//Node record id from platform database.
    
    private List<DynamicEntityValueDto> boundEntityValues = new ArrayList<>();


    public DynamicTaskNodeBindInfoDto() {
    }

    public DynamicTaskNodeBindInfoDto(String nodeId, String nodeDefId) {
        this.nodeId = nodeId;
        this.nodeDefId = nodeDefId;
    }

    public String getNodeId() {
        return nodeId;
    }

    public void setNodeId(String nodeId) {
        this.nodeId = nodeId;
    }

    public String getNodeDefId() {
        return nodeDefId;
    }

    public void setNodeDefId(String nodeDefId) {
        this.nodeDefId = nodeDefId;
    }

    public List<DynamicEntityValueDto> getBoundEntityValues() {
        return boundEntityValues;
    }

    public void setBoundEntityValues(List<DynamicEntityValueDto> boundEntityValues) {
        this.boundEntityValues = boundEntityValues;
    }
    
    public void addBoundEntityValue(DynamicEntityValueDto boundEntityValue){
        if(boundEntityValue == null){
            return;
        }
        
        if(this.boundEntityValues == null){
            this.boundEntityValues = new ArrayList<>();
        }
        
        this.boundEntityValues.add(boundEntityValue);
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", DynamicTaskNodeBindInfoDto.class.getSimpleName() + "[", "]")
                .add("nodeId='" + nodeId + "'")
                .add("nodeDefId='" + nodeDefId + "'")
                .add("boundEntityValues=" + boundEntityValues)
                .toString();
    }
}
