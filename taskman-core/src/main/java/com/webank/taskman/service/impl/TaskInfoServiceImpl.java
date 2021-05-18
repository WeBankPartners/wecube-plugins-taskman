package com.webank.taskman.service.impl;

import java.util.Date;
import java.util.List;

import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.github.pagehelper.PageHelper;
import com.github.pagehelper.PageInfo;
import com.webank.taskman.base.JsonResponse;
import com.webank.taskman.base.LocalPageableQueryResult;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.GenernalStatus;
import com.webank.taskman.converter.FormInfoConverter;
import com.webank.taskman.converter.FormItemInfoConverter;
import com.webank.taskman.converter.RequestInfoConverter;
import com.webank.taskman.converter.TaskInfoConverter;
import com.webank.taskman.domain.FormInfo;
import com.webank.taskman.domain.FormItemInfo;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.dto.TaskInfoDto;
import com.webank.taskman.dto.platform.CoreCancelTaskDto;
import com.webank.taskman.dto.req.ProceedTasksReqDto;
import com.webank.taskman.dto.req.TaskInfoQueryReqDto;
import com.webank.taskman.dto.resp.FormInfoResqDto;
import com.webank.taskman.dto.resp.RequestInfoInstanceResqDto;
import com.webank.taskman.dto.resp.TaskInfoRespDto;
import com.webank.taskman.mapper.FormItemInfoMapper;
import com.webank.taskman.mapper.TaskInfoMapper;
import com.webank.taskman.service.FormInfoService;
import com.webank.taskman.service.FormItemInfoService;
import com.webank.taskman.service.RequestInfoService;
import com.webank.taskman.service.TaskInfoService;
import com.webank.taskman.support.platform.PlatformCoreServiceRestClient;
import com.webank.taskman.support.platform.dto.CommonPlatformResponseDto;
import com.webank.taskman.support.platform.dto.PlatformPluginResponseDto;
import com.webank.taskman.support.platform.dto.ProceedTaskResponseDto;
import com.webank.taskman.support.platform.dto.TaskFormValueDto;

@Service
public class TaskInfoServiceImpl extends ServiceImpl<TaskInfoMapper, TaskInfo> implements TaskInfoService {

    private static final Logger log = LoggerFactory.getLogger(TaskInfoServiceImpl.class);

    @Autowired
    private FormInfoConverter formInfoConverter;

    @Autowired
    private FormItemInfoMapper formItemInfoMapper;

    @Autowired
    private RequestInfoConverter requestInfoConverter;

    @Autowired
    private FormItemInfoConverter formItemInfoConverter;

    @Autowired
    private FormItemInfoService formItemInfoService;

    @Autowired
    private TaskInfoConverter taskInfoConverter;

    @Autowired
    private FormInfoService formInfoService;

    @Autowired
    private PlatformCoreServiceRestClient platformServiceClient;

    @Autowired
    private RequestInfoService requestInfoService;

    private ObjectMapper objectMapper = new ObjectMapper();

    @Override
    public LocalPageableQueryResult<TaskInfoDto> selectTaskInfo(Integer page, Integer pageSize,
            TaskInfoQueryReqDto req) {
        String inSql = req.getConditionSql();// req.getEqUseRole();
        LambdaQueryWrapper<TaskInfo> queryWrapper = taskInfoConverter.convertToTaskInfoByQuery(req)
                .getLambdaQueryWrapper().inSql(!StringUtils.isEmpty(inSql), TaskInfo::getId, inSql);
        PageHelper.startPage(page, pageSize);
        PageInfo<TaskInfoDto> pages = new PageInfo<>(
                taskInfoConverter.convertToDtos(getBaseMapper().selectList(queryWrapper)));
        LocalPageableQueryResult<TaskInfoDto> queryResponse = new LocalPageableQueryResult<>(pages.getTotal(),
                page.longValue(), pageSize.longValue(), pages.getList());
        return queryResponse;
    }

    @Override
    public JsonResponse proceedUserTask(ProceedTasksReqDto req) {
        TaskInfo taskInfo = getBaseMapper().selectById(req.getRecordId());
        String currentUsername = AuthenticationContextHolder.getCurrentUsername();
        if (!currentUsername.equals(taskInfo.getReporter())) {
            throw new TaskmanRuntimeException("Failed to process. Please claim");
        }
        if (!"already_received".equals(taskInfo.getStatus())) {
            throw new TaskmanRuntimeException("Processing failed. The current task is not claimed");
        }
        
        //TODO
        asyncProceedUserTask();
//        callbackByTaskInfo(req, taskInfo);
        // List<FormItemInfo> formItemInfos =
        // formItemInfoConverter.toEntityByReq(req.getFormItemInfoList());
        // formInfoService.saveFormInfoAndItems(formItemInfos,
        // taskInfo.getTaskTempId(), taskInfo.getId());
        taskInfo.setStatus(GenernalStatus.Processed.name());
        taskInfo.setResult(req.getResult());
        taskInfo.setUpdatedBy(currentUsername);
        taskInfo.setUpdatedTime(new Date());
        updateById(taskInfo);

        return JsonResponse.okay();
    }

    private void asyncProceedUserTask() {
        // TODO

        String callbackUrl = null;
        PlatformPluginResponseDto platformPluginResponseDto = null;
        PlatformPluginResponseDto.ResultData resultData = null;

        // TODO
        ProceedTaskResponseDto proceedTaskRespDto = new ProceedTaskResponseDto();


        TaskFormValueDto taskFormValueDto = buildTaskFormValueDto();
        String taskFormOutput = convertTaskFormValueDtoToJson(taskFormValueDto);
    }

    private TaskFormValueDto buildTaskFormValueDto() {
        // TODO
        
        return null;
    }

    private String convertTaskFormValueDtoToJson(TaskFormValueDto taskFormValueDto) {
        try {
            String json = objectMapper.writeValueAsString(taskFormValueDto);
            return json;
        } catch (JsonProcessingException e) {
            log.error("Failed to write object to json.", e);
            throw new TaskmanRuntimeException("Failed to write object to json.");
        }
    }

    private void callbackByTaskInfo(ProceedTasksReqDto req, TaskInfo taskInfo) {
        // CallbackRequestDto callbackRequest = new CallbackRequestDto();
        // CallbackRequestResultDataDto callbackRequestResultDataDto = new
        // CallbackRequestResultDataDto();
        // callbackRequestResultDataDto.setRequestId(taskInfo.getRequestId());
        // callbackRequestResultDataDto.setOutputs(
        // Lists.newArrayList(new
        // CallbackRequestOutputsDto(CallbackRequestOutputsDto.ERROR_CODE_SUCCESSFUL,
        // req.getResultMessage(), req.getResultMessage(),
        // taskInfo.getCallbackParameter())));
        //
        // callbackRequest.setResults(callbackRequestResultDataDto);
        // callbackRequest.setResultMessage(req.getResultMessage());
        // callbackRequest.setResultCode(req.getResult());
        //
        // try {
        // platformServiceClient.callback(taskInfo.getCallbackUrl(),
        // callbackRequest);
        // } catch (CoreRemoteCallException e) {
        // String msg = String.format("Callback wecube meet error: %s",
        // e.getMessage());
        // log.error(msg);
        // throw new TaskmanRuntimeException("3014", msg, e.getMessage());
        // }
    }

    @Override
    public CommonPlatformResponseDto cancelTask(CoreCancelTaskDto req) {
        TaskInfo taskInfo = getBaseMapper().selectOne(new TaskInfo().setProcInstId(req.getProcInstId())
                .setNodeDefId(req.getTaskNodeId()).getLambdaQueryWrapper());
        if (null == taskInfo) {
            throw new TaskmanRuntimeException("NOT_FOUND_RECORD");
        }
        // taskInfo.setCurrenUserName(taskInfo, taskInfo.getId());
        taskInfo.setCreatedBy(AuthenticationContextHolder.getCurrentUsername());
        taskInfo.setUpdatedBy(AuthenticationContextHolder.getCurrentUsername());
        taskInfo.setUpdatedTime(new Date());
        taskInfo.setStatus(GenernalStatus.SUSPENSION.toString());
        getBaseMapper().updateById(taskInfo);
        return CommonPlatformResponseDto.okay();
    }

    @Override
    public TaskInfoRespDto taskInfoDetail(String id) {
        // TaskInfo taskInfo = getBaseMapper().selectOne(new
        // TaskInfo().setId(id).getLambdaQueryWrapper());
        // TaskInfoRespDto taskInfoResp = taskInfoConverter.toResp(taskInfo);
        // taskInfoResp.setFormItemInfo(formItemInfoService.returnDetail(id));
        // return taskInfoResp;

        return null;
    }

    @Override
    public RequestInfoInstanceResqDto selectTaskInfoInstanceService(String requestId, String taskId) {
        // RequestInfo requestInfo = requestInfoService.getOne(new
        // RequestInfo().setId(requestId).getLambdaQueryWrapper());
        // RequestInfoInstanceResqDto requestInfoInstanceResq =
        // requestInfoConverter.toInstanceResp(requestInfo);
        //
        // requestInfoInstanceResq.setRequestFormResq(getFormInfoResq(requestId));
        // List<TaskInfo> taskInfos = getBaseMapper().selectList(new
        // QueryWrapper<TaskInfo>().lambda()
        // .eq(TaskInfo::getRequestId,
        // requestId).orderByAsc(TaskInfo::getUpdatedTime));
        // List<TaskInfoInstanceRespDto> taskInfoInstanceResps = new
        // ArrayList<>();
        // for (TaskInfo taskInfo : taskInfos) {
        // if (!(taskInfo.getId().equals(taskId))) {
        // TaskInfoInstanceRespDto resp =
        // taskInfoConverter.toInstanceResp(taskInfo);
        // resp.setTaskFormResq(getFormInfoResq(taskId));
        // taskInfoInstanceResps.add(resp);
        // }
        // }
        // requestInfoInstanceResq.setTaskInfoInstanceResps(taskInfoInstanceResps);
        //
        // return requestInfoInstanceResq;

        return null;
    }

    private FormInfoResqDto getFormInfoResq(String recordId) {
        FormInfoResqDto formInfoResq = formInfoConverter
                .convertToDto(formInfoService.getOne(new FormInfo().setRecordId(recordId).getLambdaQueryWrapper()));
        if (null != formInfoResq) {
            List<FormItemInfo> formItemInfos = formItemInfoMapper
                    .selectList(new FormItemInfo().setRecordId(recordId).getLambdaQueryWrapper());
            formInfoResq.setFormItemInfo(formItemInfoConverter.convertToDtos(formItemInfos));
        }
        return formInfoResq;
    }

    @Override
    public TaskInfoDto taskInfoReceive(String id) {
        TaskInfo taskInfo = getBaseMapper().selectById(id);

        if (taskInfo.getStatus().equals(GenernalStatus.UNCLAIMED.toString())) {
            taskInfo.setStatus(GenernalStatus.ALREADY_RECEIVED.toString());
            taskInfo.setReporter(AuthenticationContextHolder.getCurrentUsername());
            taskInfo.setReportTime(new Date());
            taskInfo.setUpdatedTime(new Date());
            getBaseMapper().updateById(taskInfo);
        }
        return taskInfoConverter.convertToDto(taskInfo);
    }

    // @Override
    // @Transactional
    // public CommonResponseDto createTask(PlatformTaskCreationReqDto req) {
    // if ( req.getInputs() == null || req.getInputs().isEmpty()) {
    // throw new TaskmanRuntimeException(" inputs is null");
    // }
    // req.getInputs().stream().forEach(task -> {
    // TaskInfo taskInfo = taskInfoConverter.toEntityByReq(task);
    // /*
    // * int isExists = taskInfoMapper.selectCount( new
    // * TaskInfo(task.getProcInstId(),task.getTaskNodeId()).
    // * getLambdaQueryWrapper()); if(0 < isExists){ throw new
    // * TaskmanRuntimeException(String.format(
    // * "Task is exists! procInstId:%s,TaskNodeId:%s",
    // * task.getProcInstId(),task.getTaskNodeId())); }
    // */
    // taskInfo.setRequestId(req.getRequestId());
    // // taskInfo.setCurrenUserName(taskInfo, taskInfo.getId());
    // taskInfo.setCreatedBy(AuthenticationContextHolder.getCurrentUsername());
    // taskInfo.setUpdatedBy(AuthenticationContextHolder.getCurrentUsername());
    // save(taskInfo);
    // List<FormItemInfo> items =
    // formItemInfoConverter.toEntityByBean(task.getFormItems());
    // formInfoService.saveFormInfoAndFormItems(items, taskInfo.getTaskTempId(),
    // taskInfo.getId());
    // });
    // return CommonResponseDto.okay();
    // }

}
