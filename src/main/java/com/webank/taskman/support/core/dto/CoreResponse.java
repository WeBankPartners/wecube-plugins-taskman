package com.webank.taskman.support.core.dto;


import com.webank.taskman.dto.OperationEventResultDto;

import java.util.List;

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

    public static class GetAllRolesResponse extends CoreResponse<List<RolesDataResponse>> {
    }
    
    public static class GetRootEntitiesResponse extends CoreResponse<Object> {
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