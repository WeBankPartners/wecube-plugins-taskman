package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.github.pagehelper.PageHelper;
import com.github.pagehelper.PageInfo;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.StatusCodeEnum;
import com.webank.taskman.constant.StatusEnum;
import com.webank.taskman.converter.*;
import com.webank.taskman.domain.*;
import com.webank.taskman.dto.CoreCancelTaskDTO;
import com.webank.taskman.dto.CoreCreateTaskDTO;
import com.webank.taskman.dto.CoreCreateTaskDTO.TaskInfoReq.FormItemBean;
import com.webank.taskman.dto.TaskInfoDTO;
import com.webank.taskman.dto.req.QueryTaskInfoReq;
import com.webank.taskman.dto.req.*;
import com.webank.taskman.dto.resp.*;
import com.webank.taskman.mapper.*;
import com.webank.taskman.service.FormInfoService;
import com.webank.taskman.service.FormItemInfoService;
import com.webank.taskman.service.TaskInfoService;
import com.webank.taskman.support.core.CommonResponseDto;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import java.util.stream.Collectors;


@Service
public class TaskInfoServiceImpl extends ServiceImpl<TaskInfoMapper, TaskInfo> implements TaskInfoService {

    @Autowired
    TaskInfoMapper taskInfoMapper;

    @Autowired
    FormInfoMapper formInfoMapper;


    @Autowired
    FormInfoConverter formInfoConverter;

    @Autowired
    FormItemInfoMapper formItemInfoMapper;



    @Autowired
    RequestInfoMapper requestInfoMapper;

    @Autowired
    RequestInfoConverter requestInfoConverter;


    @Autowired
    FormItemTemplateMapper formItemTemplateMapper;

    @Autowired
    FormItemInfoConverter formItemInfoConverter;

    @Autowired
    FormItemInfoService formItemInfoService;

    @Autowired
    TaskInfoConverter taskInfoConverter;

    @Autowired
    FormInfoService formInfoService;

    @Override
    public QueryResponse<TaskInfoDTO> selectTaskInfo(Integer page, Integer pageSize, QueryTaskInfoReq req) {
        LambdaQueryWrapper<TaskInfo> queryWrapper = taskInfoConverter.toEntityByQuery(req)
                .getLambdaQueryWrapper().inSql(TaskInfo::getId,req.getEqUseRole());
        PageHelper.startPage(page,pageSize);
        PageInfo<TaskInfo> pages = new PageInfo(taskInfoMapper.selectList(queryWrapper));
        QueryResponse<TaskInfoDTO> queryResponse = new QueryResponse(pages.getTotal(),page.longValue(),pageSize.longValue(),pages.getList());
        return queryResponse;
    }

    @Override
    public TaskInfoResp selectSynthesisTaskInfoFormService(String id){
        TaskInfoResp resp =taskInfoConverter.toResp(taskInfoMapper.selectOne(new TaskInfo().setId(id).getLambdaQueryWrapper()));
        resp.setFormItemInfo(returnDetail(id));
        return resp;
    }

    @Override
    public String ProcessingTasksService(ProcessingTasksReq req) throws TaskmanRuntimeException {
        TaskInfo taskInfo=taskInfoMapper.selectById(req.getRecordId());
        String currentUsername = AuthenticationContextHolder.getCurrentUsername();
        if (!currentUsername.equals(taskInfo.getReporter())){
            return "Failed to process. Please claim";
        }
        if (!"already_received".equals(taskInfo.getStatus())){
            return "Processing failed. The current task is not claimed";
        }
        FormInfo formInfo = new FormInfo();
        formInfo.setRecordId(taskInfo.getId());
        formInfo.setUpdatedBy(currentUsername);
        formInfo.setType(StatusEnum.ENABLE.ordinal());
        formInfoMapper.insert(formInfo);
        for (SaveFormItemInfoReq formItemInfo : req.getFormItemInfoList()) {
            FormItemInfo formItemInfo1=formItemInfoConverter.toEntityByReq(formItemInfo);
            formItemInfo1.setFormId(formInfo.getId());
            formItemInfo1.setRecordId(taskInfo.getId());
            formItemInfoMapper.insert(formItemInfo1);
        }

//        List<FormItemInfo> list=formItemInfoMapper.selectList(
//                new QueryWrapper<FormItemInfo>()
//                        .eq("form_id",formInfo.getId()));
//        FormInfo formInfo1=formInfoMapper.selectById(formInfo.getId());
//        ProcessingTasksResp processingTasksResp=formInfoConverter.processingTasksResp(formInfo1);
//        processingTasksResp.setFormItemInfoList(list);

        taskInfo.setStatus("Processed");
        taskInfo.setResult(req.getResult());
        taskInfo.setUpdatedBy(currentUsername);
        taskInfoMapper.updateById(taskInfo);

        return "processing successful";
    }

    @Override
    public CommonResponseDto cancelTask(CoreCancelTaskDTO req) {
        TaskInfo taskInfo = taskInfoMapper.selectOne(
                new TaskInfo().setProcInstId(req.getProcInstId()).setNodeDefId(req.getTaskNodeId()).getLambdaQueryWrapper());
        if(null == taskInfo){
            throw new TaskmanRuntimeException(StatusCodeEnum.NOT_FOUND_RECORD);
        }
        taskInfo.setCurrenUserName(taskInfo,taskInfo.getId());
        taskInfo.setUpdatedTime(new Date());
        taskInfo.setStatus(StatusEnum.SUSPENSION.toString());
        updateById(taskInfo);
        return CommonResponseDto.okay();
    }

    @Override
    public TaskInfoResp taskInfoDetail(String id) {
        TaskInfoResp taskInfoResp = taskInfoConverter.toResp(getById(id));
        List<FormItemInfo> formItemInfo = formItemInfoService.list(new FormItemInfo().setRecordId(id).getLambdaQueryWrapper());
        taskInfoResp.setFormItemInfo(formItemInfoConverter.toDto(formItemInfo));
        return taskInfoResp;
    }

    @Override
        public RequestInfoInstanceResq selectTaskInfoInstanceService(String procInstId,String taskId) {

        RequestInfo requestInfo = requestInfoMapper.selectOne(new RequestInfo().setProcInstId(procInstId).getLambdaQueryWrapper());
        RequestInfoInstanceResq requestInfoInstanceResq = requestInfoConverter.toInstanceResp(requestInfo);

        FormInfo formInfo = formInfoMapper.selectOne(new FormInfo().setRecordId(requestInfo.getId()).getLambdaQueryWrapper());
        if (null == formInfo){
            throw new TaskmanRuntimeException("The request details do not exist");
        }
        List<FormItemInfo> formItemInfos = formItemInfoMapper.selectList(new FormItemInfo().setRecordId(requestInfo.getId()).getLambdaQueryWrapper());
        requestInfoInstanceResq.setRequestFormResq(formInfoConverter.toRequestFormResq(formInfo));
        requestInfoInstanceResq.getRequestFormResq().setFormItemInfo(formItemInfos);

        List<TaskInfo> taskInfos = taskInfoMapper.selectList( new QueryWrapper<TaskInfo>().lambda().eq(TaskInfo::getProcInstId, procInstId).orderByAsc(TaskInfo::getUpdatedTime));

        List<TaskInfoInstanceResp> taskInfoInstanceResps = new ArrayList<>();
        for (TaskInfo taskInfo : taskInfos) {
            if (!(taskInfo.getId().equals(taskId))) {
                TaskInfoInstanceResp resp = taskInfoConverter.toInstanceResp(taskInfo);
                formInfo=formInfoMapper.selectOne(new FormInfo().setRecordId(taskInfo.getId()).getLambdaQueryWrapper());
                if (null==formInfo||"".equals(formInfo)){
                    throw new TaskmanRuntimeException("The request details do not exist");
                }
                formItemInfos=formItemInfoMapper.selectList(new FormItemInfo().setFormId(formInfo.getId()).getLambdaQueryWrapper());
                resp.setTaskFormResq(formInfoConverter.toTaskFormResq(formInfo));
                resp.getTaskFormResq().setFormItemInfo(formItemInfos);

                taskInfoInstanceResps.add(resp);
            }
        }
        requestInfoInstanceResq.setTaskInfoInstanceResps(taskInfoInstanceResps);

        return requestInfoInstanceResq;
    }

    @Override
    public TaskInfoDTO taskInfoReceive(String id) {
        TaskInfo taskInfo = taskInfoMapper.selectById(id);

        if(taskInfo.getStatus().equals(StatusEnum.UNCLAIMED.toString())){
            taskInfo.setStatus(StatusEnum.ALREADY_RECEIVED.toString());
            taskInfo.setReporter(AuthenticationContextHolder.getCurrentUsername());
            taskInfo.setUpdatedTime(new Date());
            taskInfoMapper.updateById(taskInfo);
        }
        return taskInfoConverter.toDto(taskInfo);
    }

    @Override
    @Transactional
    public CommonResponseDto createTask(CoreCreateTaskDTO req) throws TaskmanRuntimeException{
        if(null == req.getInputs() || req.getInputs().size()==0){
            throw new TaskmanRuntimeException(" inputs is null");
        }
        req.getInputs().stream().forEach(task->{
            TaskInfo taskInfo = taskInfoConverter.toEntityByReq(task);
            TaskInfo isExists = taskInfoMapper.selectOne(
                    new TaskInfo(task.getProcInstId(),task.getTaskNodeId()).getLambdaQueryWrapper());
            if(null != isExists){
                throw new TaskmanRuntimeException(String.format(
                        "Task is exists! procInstId:%s,TaskNodeId:%s",
                        task.getProcInstId(),task.getTaskNodeId()));
            }
            taskInfo.setCurrenUserName(taskInfo,taskInfo.getId());
            save(taskInfo);
            List<FormItemInfo> items = formItemInfoConverter.toEntityByBeans(task.getFormItems());
            formInfoService.saveFormInfoAndItems(items,taskInfo.getTaskTempId(),taskInfo.getId());
        });
        return  CommonResponseDto.okay();
    }


    public  List<FormItemInfoResp> returnDetail(String id){
        FormInfo formInfo=formInfoMapper.selectOne(new FormInfo().setRecordId(id).getLambdaQueryWrapper());
        if (null==formInfo||"".equals(formInfo)){
            throw new TaskmanRuntimeException("The request details do not exist");
        }
        List<FormItemInfo> formItemInfos=formItemInfoMapper.selectList(new FormItemInfo().setFormId(formInfo.getId()).getLambdaQueryWrapper());
        List<FormItemInfoResp> formItemInfoResps = formItemInfoConverter.toDto(formItemInfos);
        for (FormItemInfoResp formItemInfoResp : formItemInfoResps) {
            FormItemTemplate formItemTemplate = formItemTemplateMapper.selectOne(new QueryWrapper<FormItemTemplate>().lambda().
                    eq(FormItemTemplate::getId, formItemInfoResp.getItemTempId()));
            formItemInfoResp.setElementType(formItemTemplate.getElementType());
            formItemInfoResp.setTitle(formItemTemplate.getTitle());
            formItemInfoResp.setWidth(formItemTemplate.getWidth());
            formItemInfoResp.setIsEdit(formItemTemplate.getIsEdit());
            formItemInfoResp.setIsView(formItemTemplate.getIsView());
            formItemInfoResp.setSort(formItemTemplate.getSort());
            formItemInfoResp.setName(formItemTemplate.getName());
            formItemInfoResp.setDataOptions(formItemTemplate.getDataOptions());
        }
        return formItemInfoResps;
    }
}
