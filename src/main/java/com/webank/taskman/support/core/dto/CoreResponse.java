package com.webank.taskman.support.core.dto;

import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.StringJoiner;

public class CoreResponse<DATATYPE> {

    private String status;
    private String message;
    private DATATYPE data;

    public static class DefaultCoreResponse extends CoreResponse<Object> {
    }
    public static class LinkedHashMapResponse extends CoreResponse<LinkedHashMap> {
    }

    public static class ListDataResponse extends CoreResponse<List<Object>> {
    }
    public static class ListLinkedHashMapResponse extends CoreResponse<List<LinkedHashMap>> {
    }
    public static class ListMapDataResponse extends CoreResponse<List<Map<String,Object>>> {
    }

    public static class ListRolesDataResponse extends CoreResponse<List<RolesDataResponse>> {
    }


    public static class ListWorkflowDefInfoResponse extends CoreResponse<List<WorkflowDefInfoDto>> {
    }

    public static class ListWorkflowNodeDefInfoResponse extends CoreResponse<List<WorkflowNodeDefInfoDto>> {
    }

    public static class ListPluginPackageDataModelResponse extends CoreResponse<List<PluginPackageDataModelDto>> {

    }
    public static class ListPluginPackageAttributeResponse extends CoreResponse<List<PluginPackageAttributeDto>> {

    }


    public static class ProcessDataPreviewResponse extends CoreResponse<LinkedHashMap> {

    }


    public static class ListTaskNodeDefObjectBindInfoResponse extends CoreResponse<List<TaskNodeDefObjectBindInfoDto>> {
    }

    public static class DynamicWorkflowInstInfoResponse extends CoreResponse<DynamicWorkflowInstInfoDto> {
    }



    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public DATATYPE getData() {
        return data;

    }

    @SuppressWarnings("unchecked")
    public void setData(Object data) {
        this.data = (DATATYPE) data;
    }

    @Override
    public String toString() {
        return "CoreResponse{" +
                "status='" + status + '\'' +
                ", message='" + message + '\'' +
                ", data=" + data +
                '}';
    }
}