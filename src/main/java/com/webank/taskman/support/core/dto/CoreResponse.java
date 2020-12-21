package com.webank.taskman.support.core.dto;

import java.util.List;
import java.util.Map;
import java.util.Set;

public class CoreResponse<DATATYPE> {

    private String status;
    private String message;
    private DATATYPE data;

    public static class DefaultCoreResponse extends CoreResponse<Object> {
    }
    
    public static class OperationEventResultResponse extends CoreResponse<OperationEventResultDto>{
        
    }

    public static class IntegerCoreResponse extends CoreResponse<Integer> {
    }

    public static class StringCoreResponse extends CoreResponse<String> {
    }

    public static class ListDataResponse extends CoreResponse<List<Object>> {
    }

    public static class SetDataResponse extends CoreResponse<Set<Object>> {
    }
    public static class ReviewEntitiesDTOResponse extends CoreResponse<ProcessDataPreviewDto> {
    }

    public static class ListMapDataResponse extends CoreResponse<List<Map<String,Object>>> {
    }

    public static class GetAllRolesResponse extends CoreResponse<List<RolesDataResponse>> {
    }

    public static class GetModelsAllResponse extends CoreResponse<List<PluginPackageDataModelDto>> {

    }
    public static class GetAttributesByPackageEntityResponse extends CoreResponse<List<PluginPackageAttributeDto>> {

    }

    public static class DynamicWorkflowInstInfoDto extends CoreResponse<DynamicWorkflowInstInfoDto> {
    }

    public static class CommonResponseDto extends CoreResponse<List<WorkflowDefInfoDto>> {
        @Override
        public String toString() {
            return "CommonResponseDto [getStatus()=" + getStatus() + ", getMessage()=" + getMessage()
                    + ", getData()=" + getData() + ", getClass()=" + getClass() + ", hashCode()=" + hashCode()
                    + ", toString()=" + super.toString() + "]";
        }
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
}