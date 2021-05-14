package com.webank.taskman.service;

import java.util.Collections;
import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.constant.RecordDeleteFlag;
import com.webank.taskman.domain.FormTemplate;
import com.webank.taskman.domain.RequestInfo;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.domain.TaskTemplate;
import com.webank.taskman.mapper.FormTemplateMapper;
import com.webank.taskman.mapper.RequestInfoMapper;
import com.webank.taskman.mapper.RequestTemplateMapper;
import com.webank.taskman.mapper.TaskTemplateMapper;
import com.webank.taskman.support.platform.dto.CreateTaskRequestInputDto;
import com.webank.taskman.support.platform.dto.PlatformPluginRequestDto;
import com.webank.taskman.support.platform.dto.TaskFormMetaDto;

@Service
public class TaskmanPluginServiceManagementService {

    private static final Logger log = LoggerFactory.getLogger(TaskmanPluginServiceManagementService.class);

    @Autowired
    private RequestInfoService requestInfoService;

    @Autowired
    private RequestInfoMapper requestInfoMapper;
    @Autowired
    private RequestTemplateMapper requestTemplateMapper;
    @Autowired
    private TaskTemplateMapper taskTemplateMapper;
    @Autowired
    private FormTemplateMapper formTemplateMapper;
    @Autowired
    private TaskInfoService taskInfoService;

    private ObjectMapper objectMapper = new ObjectMapper();

    /**
     * 
     * @param requestDto
     */
    public void createTask(PlatformPluginRequestDto requestDto) {
        // TODO
        log.info("About to create task with request:{}", requestDto);
        TaskInfo taskInfo = new TaskInfo();
        taskInfo.setNodeRequestId(requestDto.getRequestId());
        taskInfo.setDueDate(requestDto.getDueDate());
        taskInfo.setReportTime(new Date());
        if (requestDto.getAllowedOptions() != null && !requestDto.getAllowedOptions().isEmpty()) {
            String allowedOptions = convertObjectToJson(requestDto.getAllowedOptions());
            taskInfo.setAllowedOptions(allowedOptions);
        }else{
            String allowedOptions = "[\"deny\",\"approval\"]";
            taskInfo.setAllowedOptions(allowedOptions);
        }
        
        taskInfo.setUpdatedBy(AuthenticationContextHolder.getCurrentUsername());
        taskInfo.setUpdatedTime(new Date());
        taskInfo.setCreatedBy(AuthenticationContextHolder.getCurrentUsername());
        taskInfo.setCreatedTime(new Date());
        
        //TODO
        CreateTaskRequestInputDto createTaskRequestInputDto = tryPickoutPlatformRequestObject(requestDto);
        
        
        taskInfoService.save(taskInfo);
    }

    /**
     * 
     * @param procInstId
     * @param nodeDefId
     * @return
     */
    public TaskFormMetaDto fetchTaskCreationMeta(String procInstId, String nodeDefId) {
        log.info("try to fetch task creation meta for process intance:{} and node:{}", procInstId, nodeDefId);
        RequestInfo requestInfoCriteria = new RequestInfo();
        requestInfoCriteria.setProcInstId(procInstId);

        LambdaQueryWrapper<RequestInfo> requestInfoQueryWrapper = requestInfoCriteria.getLambdaQueryWrapper();
        List<RequestInfo> requestInfos = requestInfoMapper.selectList(requestInfoQueryWrapper);

        if (requestInfos == null || requestInfos.isEmpty()) {
            log.info("Cannot find request infomation with process instance id :{}", procInstId);
            return null;
        }
        RequestInfo requestInfo = requestInfos.get(0);
        String requestTemplateId = requestInfo.getRequestTempId();

        RequestTemplate requestTemplate = requestTemplateMapper.selectById(requestTemplateId);
        String procDefId = requestTemplate.getProcDefId();

        TaskTemplate taskTemplateCriteria = new TaskTemplate();
        taskTemplateCriteria.setProcDefId(procDefId);
        taskTemplateCriteria.setNodeDefId(nodeDefId);
        taskTemplateCriteria.setRequestTemplateId(requestTemplate.getId());
        taskTemplateCriteria.setDelFlag(RecordDeleteFlag.NotDeleted.ordinal());

        LambdaQueryWrapper<TaskTemplate> taskTemplateQueryWrapper = taskTemplateCriteria.getLambdaQueryWrapper();
        List<TaskTemplate> taskTemplates = taskTemplateMapper.selectList(taskTemplateQueryWrapper);

        if (taskTemplates == null || taskTemplates.isEmpty()) {
            return null;
        }

        TaskTemplate taskTemplate = taskTemplates.get(0);

        FormTemplate formTemplateCriteria = new FormTemplate();
        formTemplateCriteria.setDelFlag(RecordDeleteFlag.NotDeleted.ordinal());
        formTemplateCriteria.setTempId(taskTemplate.getId());
        // TODO
        
       
        TaskFormMetaDto taskFormMetaDto = buildTaskFormMetaDto();
        return taskFormMetaDto;
    }
    
    private CreateTaskRequestInputDto tryPickoutPlatformRequestObject(PlatformPluginRequestDto requestDto){
        List<CreateTaskRequestInputDto> inputs = requestDto.getInputs();
        if(inputs == null || inputs.isEmpty()){
            return null;
        }
        
       return inputs.get(0);
    }

    private TaskFormMetaDto buildTaskFormMetaDto() {
        // TODO
        return null;
    }

    private <T> T convertJsonToObject(String json, Class<T> valueType) {
        try {
            T t = objectMapper.readValue(json, valueType);
            return t;
        } catch (Exception e) {
            log.error("Failed to convert json to object.", e);
            return null;
        }
    }

    private String convertObjectToJson(Object obj) {
        try {
            String json = objectMapper.writeValueAsString(obj);
            return json;
        } catch (JsonProcessingException e) {
            log.error("Failed to convert object to json.", e);
            return null;
        }
    }

}
