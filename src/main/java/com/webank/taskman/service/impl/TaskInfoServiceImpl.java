package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.github.pagehelper.PageHelper;
import com.github.pagehelper.PageInfo;
import com.google.common.collect.Lists;
import com.webank.taskman.base.JsonResponse;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.StatusCodeEnum;
import com.webank.taskman.constant.StatusEnum;
import com.webank.taskman.converter.*;
import com.webank.taskman.domain.*;
import com.webank.taskman.dto.CoreCancelTaskDTO;
import com.webank.taskman.dto.CoreCreateTaskDTO;
import com.webank.taskman.dto.TaskInfoDTO;
import com.webank.taskman.dto.req.QueryTaskInfoReq;
import com.webank.taskman.dto.req.*;
import com.webank.taskman.dto.resp.*;
import com.webank.taskman.mapper.*;
import com.webank.taskman.service.FormInfoService;
import com.webank.taskman.service.FormItemInfoService;
import com.webank.taskman.service.RequestInfoService;
import com.webank.taskman.service.TaskInfoService;
import com.webank.taskman.support.core.CommonResponseDto;
import com.webank.taskman.support.core.CoreRemoteCallException;
import com.webank.taskman.support.core.PlatformCoreServiceRestClient;
import com.webank.taskman.support.core.dto.CallbackRequestDto;
import com.webank.taskman.support.core.dto.CallbackRequestDto.*;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.ArrayList;
import java.util.Date;
import java.util.List;


@Service
public class TaskInfoServiceImpl extends ServiceImpl<TaskInfoMapper, TaskInfo> implements TaskInfoService {


    @Autowired
    FormInfoConverter formInfoConverter;

    @Autowired
    FormItemInfoMapper formItemInfoMapper;


    @Autowired
    RequestInfoConverter requestInfoConverter;

    @Autowired
    FormItemInfoConverter formItemInfoConverter;

    @Autowired
    FormItemInfoService formItemInfoService;

    @Autowired
    TaskInfoConverter taskInfoConverter;

    @Autowired
    FormInfoService formInfoService;

    @Autowired
    PlatformCoreServiceRestClient coreServiceStub;

    @Override
    public QueryResponse<TaskInfoDTO> selectTaskInfo(Integer page, Integer pageSize, QueryTaskInfoReq req) {
        String inSql = req.getConditionSql();// req.getEqUseRole();
        LambdaQueryWrapper<TaskInfo> queryWrapper = taskInfoConverter.toEntityByQuery(req)
                .getLambdaQueryWrapper().inSql(!StringUtils.isEmpty(inSql), TaskInfo::getId, inSql);
        PageHelper.startPage(page, pageSize);
        PageInfo<TaskInfoDTO> pages = new PageInfo( taskInfoConverter.toDto(getBaseMapper().selectList(queryWrapper)));
        QueryResponse<TaskInfoDTO> queryResponse = new QueryResponse(pages.getTotal(), page.longValue(), pageSize.longValue(), pages.getList());
        return queryResponse;
    }

    @Override
    @Transactional
    public JsonResponse taskInfoProcessing(ProcessingTasksReq req) throws TaskmanRuntimeException {
        TaskInfo taskInfo = getBaseMapper().selectById(req.getRecordId());
        String currentUsername = AuthenticationContextHolder.getCurrentUsername();
        if (!currentUsername.equals(taskInfo.getReporter())) {
            throw new TaskmanRuntimeException("Failed to process. Please claim");
        }
        if (!"already_received".equals(taskInfo.getStatus())) {
            throw new TaskmanRuntimeException("Processing failed. The current task is not claimed");
        }
        callbackByTaskInfo(req, taskInfo);
//        List<FormItemInfo> formItemInfos = formItemInfoConverter.toEntityByReq(req.getFormItemInfoList());
//        formInfoService.saveFormInfoAndItems(formItemInfos, taskInfo.getTaskTempId(), taskInfo.getId());
        taskInfo.setStatus(StatusEnum.Processed.name());
        taskInfo.setResult(req.getResult());
        taskInfo.setUpdatedBy(currentUsername);
        taskInfo.setUpdatedTime(new Date());
        updateById(taskInfo);

        return JsonResponse.okay();
    }

    private void callbackByTaskInfo(ProcessingTasksReq req, TaskInfo taskInfo) {
        CallbackRequestDto callbackRequest = new CallbackRequestDto();
        CallbackRequestResultDataDto callbackRequestResultDataDto = new CallbackRequestResultDataDto();
        callbackRequestResultDataDto.setRequestId(taskInfo.getRequestId());
        callbackRequestResultDataDto.setOutputs(Lists.newArrayList(new CallbackRequestOutputsDto(
                CallbackRequestOutputsDto.ERROR_CODE_SUCCESSFUL, req.getResultMessage(),
                req.getResultMessage(), taskInfo.getCallbackParameter())));

        callbackRequest.setResults(callbackRequestResultDataDto);
        callbackRequest.setResultMessage(req.getResultMessage());
        callbackRequest.setResultCode(req.getResult());

        try {
            coreServiceStub.callback(taskInfo.getCallbackUrl(), callbackRequest);
        } catch (CoreRemoteCallException e) {
            String msg = String.format("Callback wecube meet error: %s", e.getMessage());
            log.error(msg);
            throw new TaskmanRuntimeException("3014", msg, e.getMessage());
        }
    }

    @Override
    public CommonResponseDto cancelTask(CoreCancelTaskDTO req) {
        TaskInfo taskInfo = getBaseMapper().selectOne(
                new TaskInfo().setProcInstId(req.getProcInstId()).setNodeDefId(req.getTaskNodeId()).getLambdaQueryWrapper());
        if (null == taskInfo) {
            throw new TaskmanRuntimeException(StatusCodeEnum.NOT_FOUND_RECORD);
        }
        taskInfo.setCurrenUserName(taskInfo, taskInfo.getId());
        taskInfo.setUpdatedTime(new Date());
        taskInfo.setStatus(StatusEnum.SUSPENSION.toString());
        getBaseMapper().updateById(taskInfo);
        return CommonResponseDto.okay();
    }

    @Override
    public TaskInfoResp taskInfoDetail(String id) {
        TaskInfo taskInfo = getBaseMapper().selectOne(new TaskInfo().setId(id).getLambdaQueryWrapper());
        TaskInfoResp taskInfoResp = taskInfoConverter.toResp(taskInfo);
        taskInfoResp.setFormItemInfo(formItemInfoService.returnDetail(id));
        return taskInfoResp;
    }

    @Autowired
    RequestInfoService requestInfoService;

    @Override
    public RequestInfoInstanceResq selectTaskInfoInstanceService(String requestId, String taskId) {
        RequestInfo requestInfo = requestInfoService.getOne(new RequestInfo().setId(requestId).getLambdaQueryWrapper());
        RequestInfoInstanceResq requestInfoInstanceResq = requestInfoConverter.toInstanceResp(requestInfo);

        requestInfoInstanceResq.setRequestFormResq(getFormInfoResq(requestId));
        List<TaskInfo> taskInfos = getBaseMapper().selectList(new QueryWrapper<TaskInfo>().lambda().eq(TaskInfo::getRequestId,requestId).orderByAsc(TaskInfo::getUpdatedTime));
        List<TaskInfoInstanceResp> taskInfoInstanceResps = new ArrayList<>();
        for (TaskInfo taskInfo : taskInfos) {
            if (!(taskInfo.getId().equals(taskId))) {
                TaskInfoInstanceResp resp = taskInfoConverter.toInstanceResp(taskInfo);
                resp.setTaskFormResq(getFormInfoResq(taskId));
                taskInfoInstanceResps.add(resp);
            }
        }
        requestInfoInstanceResq.setTaskInfoInstanceResps(taskInfoInstanceResps);

        return requestInfoInstanceResq;
    }

    private FormInfoResq getFormInfoResq(String recordId) {
        FormInfoResq formInfoResq = formInfoConverter.toDto(formInfoService.getOne(new FormInfo().setRecordId(recordId).getLambdaQueryWrapper()));
        if (null != formInfoResq) {
            List<FormItemInfo> formItemInfos = formItemInfoMapper.selectList(new FormItemInfo().setRecordId(recordId).getLambdaQueryWrapper());
            formInfoResq.setFormItemInfo(formItemInfoConverter.toDto(formItemInfos));
        }
        return formInfoResq;
    }

    @Override
    public TaskInfoDTO taskInfoReceive(String id) {
        TaskInfo taskInfo = getBaseMapper().selectById(id);

        if (taskInfo.getStatus().equals(StatusEnum.UNCLAIMED.toString())) {
            taskInfo.setStatus(StatusEnum.ALREADY_RECEIVED.toString());
            taskInfo.setReporter(AuthenticationContextHolder.getCurrentUsername());
            taskInfo.setReportTime(new Date());
            taskInfo.setUpdatedTime(new Date());
            getBaseMapper().updateById(taskInfo);
        }
        return taskInfoConverter.toDto(taskInfo);
    }

    @Override
    @Transactional
    public CommonResponseDto createTask(CoreCreateTaskDTO req) throws TaskmanRuntimeException {
        if (null == req.getInputs() || req.getInputs().size() == 0) {
            throw new TaskmanRuntimeException(" inputs is null");
        }
        req.getInputs().stream().forEach(task -> {
            TaskInfo taskInfo = taskInfoConverter.toEntityByReq(task);
            /*int isExists = taskInfoMapper.selectCount(
                    new TaskInfo(task.getProcInstId(),task.getTaskNodeId()).getLambdaQueryWrapper());
            if(0 < isExists){
                throw new TaskmanRuntimeException(String.format(
                        "Task is exists! procInstId:%s,TaskNodeId:%s",
                        task.getProcInstId(),task.getTaskNodeId()));
            }*/
            taskInfo.setRequestId(req.getRequestId());
            taskInfo.setCurrenUserName(taskInfo, taskInfo.getId());
            save(taskInfo);
            List<FormItemInfo> items = formItemInfoConverter.toEntityByBean(task.getFormItems());
            formInfoService.saveFormInfoAndItems(items, taskInfo.getTaskTempId(), taskInfo.getId());
        });
        return CommonResponseDto.okay();
    }

}
