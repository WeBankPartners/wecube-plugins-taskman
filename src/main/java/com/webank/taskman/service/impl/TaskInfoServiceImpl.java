package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.github.pagehelper.PageHelper;
import com.github.pagehelper.PageInfo;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.StatusEnum;
import com.webank.taskman.converter.*;
import com.webank.taskman.domain.*;
import com.webank.taskman.dto.CheckTaskDTO;
import com.webank.taskman.dto.CoreCreateTaskDTO;
import com.webank.taskman.dto.CoreCreateTaskDTO.TaskInfoReq.FormItemBean;
import com.webank.taskman.dto.TaskInfoDTO;
import com.webank.taskman.dto.req.QueryTaskInfoReq;
import com.webank.taskman.dto.req.*;
import com.webank.taskman.dto.resp.*;
import com.webank.taskman.mapper.*;
import com.webank.taskman.service.FormItemInfoService;
import com.webank.taskman.service.TaskInfoService;
import com.webank.taskman.support.core.CommonResponseDto;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.ArrayList;
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
    SynthesisTaskInfoFormTaskConverter synthesisTaskInfoFormTaskConverter;

    @Autowired
    RequestInfoMapper requestInfoMapper;

    @Autowired
    RequestInfoInstanceConverter requestInfoInstanceConverter;

    @Autowired
    TaskInfoInstanceConverter taskInfoInstanceConverter;

    @Autowired
    TaskInfoGetConverter taskInfoGetConverter;

    @Autowired
    FormItemInfoRespConverter formItemInfoRespConverter;

    @Autowired
    FormItemTemplateMapper formItemTemplateMapper;

    @Autowired
    FormItemTemplateRespConverter formItemTemplateRespConverter;

    @Autowired
    FormItemInfoConverter formItemInfoConverter;

    @Autowired
    RequestFormInfoConverter requestFormInfoConverter;

    @Autowired
    TaskFormInfoConverter taskFormInfoConverter;

    @Override
    public QueryResponse<TaskInfoDTO> selectTaskInfo(Integer page, Integer pageSize, QueryTaskInfoReq req) {
        String currentUserRolesToString = AuthenticationContextHolder.getCurrentUserRolesToString();
        req.setSourceTableFix("tt");
        req.setUseRoleName(currentUserRolesToString);
        PageHelper.startPage(page,pageSize);

        PageInfo<TaskInfo> pages = new PageInfo(taskInfoMapper.selectTaskInfo(req));
        QueryResponse<TaskInfoDTO> queryResponse = new QueryResponse(pages.getTotal(),page.longValue(),pageSize.longValue(),pages.getList());
        return queryResponse;
    }

    /*@Override
    public SaveTaskInfoResp saveTaskInfo(SaveTaskInfoReq saveTaskInfoReq) {
        String currentUsername = AuthenticationContextHolder.getCurrentUsername();
        TaskInfo taskInfo = taskInfoConverter.svTOInfo(saveTaskInfoReq);
        taskInfo.setUpdatedBy(currentUsername);
        if (StringUtils.isEmpty(taskInfo.getId())) {
            taskInfo.setCreatedBy(currentUsername);
            taskInfoMapper.insert(taskInfo);
        }
        String taskInfoId = taskInfo.getId();
        FormInfoResq formInfoResq = checkTheTask(taskInfoId).getFormInfoResq();
        FormInfo formInfo = formInfoConverter.saveReqToEntity(saveTaskInfoReq.getFormInfo());
        List<FormItemInfo> formItemInfos = formItemInfoConverter.toEntity(saveTaskInfoReq.getFormInfo().getFormItems());

        List<FormItemTemplate> formItemTemplateList = new ArrayList<>();
        String msg = "success";
        formItemInfos.stream().forEach(f -> {
            QueryWrapper<FormItemTemplate> queryWrapper = new QueryWrapper<>();
            queryWrapper.eq("id", f.getItemTempId());
            formItemTemplateList.add(formItemTemplateMapper.selectOne(new FormItemTemplate(f.getItemTempId()).getLambdaQueryWrapper() ));
        });

        String Regular = null;
        for (FormItemTemplate formItemTemplate : formItemTemplateList) {
            for (FormItemInfo itemInfo : formItemInfos) {
                if (formItemTemplate.getId().equals(itemInfo.getItemTempId())) {
                    if (0 == formItemTemplate.getRequired()) {
                        if (StringUtils.isEmpty(itemInfo.getValue())) {
                            msg = itemInfo.getName() + "必须填写";
                            throw new TaskmanRuntimeException(msg);
                        }
                    }
                    Regular = formItemTemplate.getRegular();
                    boolean isMatch = Pattern.matches(Regular, itemInfo.getValue());
                    if (false == isMatch) {
                        msg = itemInfo.getName() + "不符合规则";
                        throw new TaskmanRuntimeException(msg);
                    }
                }
            }
        }

        if (formInfoResq == null) {
            formInfo.setRecordId(taskInfoId);
            formInfo.setCreatedBy(currentUsername);
            formInfo.setUpdatedBy(currentUsername);
            formInfoMapper.insert(formInfo);
            formItemInfos.stream().forEach(f -> {
                f.setFormId(formInfo.getId());
                formItemInfoMapper.insert(f);
            });
        }
        if (formInfoResq != null) {
            formInfo.setUpdatedBy(currentUsername);
            formInfoMapper.updateById(formInfo);
            formItemInfos.stream().forEach(f -> {
                QueryWrapper<FormItemInfo> queryWrapper = new QueryWrapper<>();
                queryWrapper.eq("form_id", f.getFormId())
                        .eq("name", f.getName());
                formItemInfoMapper.update(f, queryWrapper);
            });
        }
        SaveTaskInfoResp saveTaskInfoResp=new SaveTaskInfoResp();
        saveTaskInfoResp.setId(taskInfoId);
        return saveTaskInfoResp;
    }*/

    @Override
    public SynthesisTaskInfoFormTask selectSynthesisTaskInfoFormService(String id) throws Exception{
        SynthesisTaskInfoFormTask srt=synthesisTaskInfoFormTaskConverter.toDto(taskInfoMapper.selectOne(new TaskInfo().setId(id).getLambdaQueryWrapper()));
        srt.setFormItemInfo(returnDetail(id));
        return srt;
    }

    @Override
    public String ProcessingTasksService(ProcessingTasksReq ptr) throws Exception {
        TaskInfo taskInfo=taskInfoMapper.selectById(ptr.getRecordId());
        String currentUsername = AuthenticationContextHolder.getCurrentUsername();
        if (!currentUsername.equals(taskInfo.getReporter())){
            return "Failed to process. Please claim";
        }
        if (!"already_received".equals(taskInfo.getStatus())){
            return "Processing failed. The current task is not claimed";
        }
        FormInfo formInfo=formInfoConverter.ProcessingTasks(ptr);
        formInfo.setCreatedBy(currentUsername);
        formInfo.setUpdatedBy(currentUsername);
        formInfo.setType(1);
        formInfoMapper.insert(formInfo);
        for (FormItemInfoReq formItemInfo : ptr.getFormItemInfoList()) {
            FormItemInfo formItemInfo1=formItemInfoConverter.processTask(formItemInfo);
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
        taskInfo.setResult(ptr.getResult());
        taskInfo.setUpdatedBy(currentUsername);
        taskInfoMapper.updateById(taskInfo);

        return "processing successful";
    }

    @Override
    public RequestInfoInstanceResq selectTaskInfoInstanceService(String taskId, String requestId) {
        RequestInfo requestInfo = requestInfoMapper.selectOne(new QueryWrapper<RequestInfo>().lambda().eq(RequestInfo::getId, requestId));
        RequestInfoInstanceResq requestInfoInstanceResq = requestInfoInstanceConverter.toDto(requestInfo);

        FormInfo formInfo=formInfoMapper.selectOne(new FormInfo().setRecordId(requestInfo.getId()).getLambdaQueryWrapper());
        if (null==formInfo||"".equals(formInfo)){
            throw new TaskmanRuntimeException("The request details do not exist");
        }
        List<FormItemInfo> formItemInfos=formItemInfoMapper.selectList(new FormItemInfo().setFormId(formInfo.getId()).getLambdaQueryWrapper());
        requestInfoInstanceResq.setRequestFormResq(requestFormInfoConverter.toDto(formInfo));
        requestInfoInstanceResq.getRequestFormResq().setFormItemInfo(formItemInfos);

        List<TaskInfo> taskInfos = taskInfoMapper.selectList(new QueryWrapper<TaskInfo>().lambda().eq(TaskInfo::getRequestId, requestId).orderBy(true,true,TaskInfo::getUpdatedTime));

        List<TaskInfoInstanceResp> taskInfoInstanceResps = new ArrayList<>();
        for (TaskInfo taskInfo : taskInfos) {
            if (!(taskInfo.getId().equals(taskId))) {
                TaskInfoInstanceResp resp =taskInfoInstanceConverter.toDto(taskInfo);
                formInfo=formInfoMapper.selectOne(new FormInfo().setRecordId(taskInfo.getId()).getLambdaQueryWrapper());
                if (null==formInfo||"".equals(formInfo)){
                    throw new TaskmanRuntimeException("The request details do not exist");
                }
                formItemInfos=formItemInfoMapper.selectList(new FormItemInfo().setFormId(formInfo.getId()).getLambdaQueryWrapper());
                resp.setTaskFormResq(taskFormInfoConverter.toDto(formInfo));
                resp.getTaskFormResq().setFormItemInfo(formItemInfos);

                taskInfoInstanceResps.add(resp);
            }
        }
        requestInfoInstanceResq.setTaskInfoInstanceResps(taskInfoInstanceResps);

        return requestInfoInstanceResq;
    }

    @Override
    public TaskInfoGetResp getTheTaskInfoService(String id) {
        TaskInfo taskInfo = taskInfoMapper.selectOne(new QueryWrapper<TaskInfo>().lambda().eq(TaskInfo::getId, id));
        TaskInfoGetResp taskInfoGetResp=new TaskInfoGetResp();
        if(taskInfo.getStatus().equals(StatusEnum.UNCLAIMED.toString())){
            taskInfo.setStatus(StatusEnum.ALREADY_RECEIVED.toString());
            taskInfo.setReporter(AuthenticationContextHolder.getCurrentUsername());
            taskInfoMapper.updateById(taskInfo);
            taskInfoGetResp= taskInfoGetConverter.toDto(taskInfo);
        }
        return taskInfoGetResp;
    }

    @Autowired
    FormItemInfoService formItemInfoService;

    @Override
    @Transactional
    public CommonResponseDto createTask(CoreCreateTaskDTO req) throws TaskmanRuntimeException{
        if(null == req.getInputs() || req.getInputs().size()==0){
            throw new TaskmanRuntimeException(" inputs is null");
        }
        req.getInputs().stream().forEach(task->{
            TaskInfo taskInfo = new TaskInfo();
            taskInfo.setProcInstKey(task.getProcInstKey());
            taskInfo.setNodeDefId(task.getTaskNodeId());
            TaskInfo isExists =  taskInfoMapper.selectOne(taskInfo.getLambdaQueryWrapper());
            if(null != isExists){
                throw new TaskmanRuntimeException(String.format(
                        "Task is exists! procInstKey:%s,TaskNodeId:%s",
                        task.getProcInstKey(),task.getTaskNodeId()));
            }
            taskInfo.setProcDefId(task.getProcDefId());
            taskInfo.setProcDefName(task.getProcDefName());
            taskInfo.setNodeName(task.getTaskName());
            taskInfo.setReporter(task.getReporter());
            taskInfo.setCallbackParameter(task.getCallbackParameter());
            taskInfo.setCallbackUrl(task.getCallbackUrl());
            taskInfo.setDescription(task.getTaskDescription());
            taskInfo.setOverTime(task.getOverTime());
            taskInfo.setReportRole(task.getRoleName());
            taskInfo.setCurrenUserName(taskInfo,taskInfo.getId());
            saveOrUpdate(taskInfo);
            List<FormItemBean> items = task.getFormItems();
            if(null != items && items.size() > 0){
                items.stream().forEach(item->{
                    FormItemInfo formItemInfo = new FormItemInfo();
                    formItemInfo.setRecordId(taskInfo.getId());
                    formItemInfo.setItemTempId(item.getItemId());
                    formItemInfo.setName(item.getKey());
                    formItemInfo.setValue(item.getVal().stream().collect(Collectors.joining(",")));
                    formItemInfoService.save(formItemInfo);
                });
            }
        });
        return  CommonResponseDto.okay();
    }

    public CheckTaskDTO checkTheTask(String taskId) {
        FormInfo formInfo = formInfoMapper.selectOne(new QueryWrapper<FormInfo>().eq("record_id", taskId));
        FormInfoResq formInfoResq = formInfoConverter.toDto(formInfo);
        CheckTaskDTO checkTaskDTO = new CheckTaskDTO();
        if (formInfoResq == null) {
            checkTaskDTO.setFormInfoResq(null);
            return checkTaskDTO;
        }
        formInfoResq.setFormItemInfo(formItemInfoMapper.selectFormItemInfo(taskId));
        checkTaskDTO.setFormInfoResq(formInfoResq);
        return checkTaskDTO;
    }

    public  List<FormItemInfoResp> returnDetail(String id){
        FormInfo formInfo=formInfoMapper.selectOne(new FormInfo().setRecordId(id).getLambdaQueryWrapper());
        if (null==formInfo||"".equals(formInfo)){
            throw new TaskmanRuntimeException("The request details do not exist");
        }
        List<FormItemInfo> formItemInfos=formItemInfoMapper.selectList(new FormItemInfo().setFormId(formInfo.getId()).getLambdaQueryWrapper());
        List<FormItemInfoResp> formItemInfoResps = formItemInfoRespConverter.toDto(formItemInfos);
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
