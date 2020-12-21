package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.converter.*;
import com.webank.taskman.domain.FormInfo;
import com.webank.taskman.domain.FormItemInfo;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.dto.CheckTaskDTO;
import com.webank.taskman.base.PageInfo;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.dto.req.SaveTaskInfoAndFormInfoReq;
import com.webank.taskman.dto.req.SelectTaskInfoReq;
import com.webank.taskman.dto.req.SynthesisTaskInfoReq;
import com.webank.taskman.dto.resp.*;
import com.webank.taskman.mapper.FormInfoMapper;
import com.webank.taskman.mapper.FormItemInfoMapper;
import com.webank.taskman.mapper.FormItemTemplateMapper;
import com.webank.taskman.mapper.TaskInfoMapper;
import com.webank.taskman.service.TaskInfoService;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.regex.Pattern;


@Service
public class TaskInfoServiceImpl extends ServiceImpl<TaskInfoMapper, TaskInfo> implements TaskInfoService {

    @Autowired
    TaskInfoMapper taskInfoMapper;

    @Autowired
    TaskInfoConverter taskInfoConverter;

    @Autowired
    FormInfoMapper formInfoMapper;

    @Autowired
    FormItemTemplateMapper formItemTemplateMapper;

    @Autowired
    FormInfoConverter formInfoConverter;

    @Autowired
    FormItemInfoMapper formItemInfoMapper;

    @Autowired
    FormItemInfoConverter formItemInfoConverter;

    @Autowired
    SynthesisTaskInfoRespConverter synthesisTaskInfoRespConverter;

    @Autowired
    SynthesisTaskInfoFormTaskConverter synthesisTaskInfoFormTaskConverter;

    @Override
    public QueryResponse<TaskInfoResp> selectTaskInfoService(Integer page, Integer pageSize, SelectTaskInfoReq req) {
        String currentUserRolesToString = AuthenticationContextHolder.getCurrentUserRolesToString();
        IPage<TaskInfo> iPage = taskInfoMapper.selectTaskInfo(new Page<>(page, pageSize), currentUserRolesToString);
        List<TaskInfoResp> respList = taskInfoConverter.toDto(iPage.getRecords());
        for (TaskInfoResp taskInfoResp : respList) {
            FormInfo formInfo = formInfoMapper.selectOne(new QueryWrapper<FormInfo>().eq("record_id", taskInfoResp.getId()));
            FormInfoResq formInfoResq = formInfoConverter.toDto(formInfo);
            if (formInfoResq != null) {
                formInfoResq.setFormItemInfo(formItemInfoMapper.selectFormItemInfo(taskInfoResp.getId()));
                taskInfoResp.setFormInfoResq(formInfoResq);
            }
        }

        QueryResponse<TaskInfoResp> queryResponse = new QueryResponse<>();
        queryResponse.setPageInfo(new PageInfo(iPage.getTotal(), iPage.getCurrent(), iPage.getSize()));
        queryResponse.setContents(respList);
        return queryResponse;
    }

    @Override
    public SaveTaskInfoResp saveTaskInfo(SaveTaskInfoAndFormInfoReq saveTaskInfoAndFormInfoReq) {
        String currentUsername = AuthenticationContextHolder.getCurrentUsername();
        TaskInfo taskInfo = taskInfoConverter.svTOInfo(saveTaskInfoAndFormInfoReq);
        taskInfo.setUpdatedBy(currentUsername);
        if (StringUtils.isEmpty(taskInfo.getId())) {
            taskInfo.setCreatedBy(currentUsername);
            taskInfoMapper.insert(taskInfo);
        }
        String taskInfoId = taskInfo.getId();
        FormInfoResq formInfoResq = checkTheTask(taskInfoId).getFormInfoResq();
        FormInfo formInfo = formInfoConverter.svToFormInfo(saveTaskInfoAndFormInfoReq.getSaveFormInfoAndFormItemInfoReq());
        List<FormItemInfo> formItemInfos = formItemInfoConverter.toEntity(saveTaskInfoAndFormInfoReq.getSaveFormInfoAndFormItemInfoReq().getSaveFormItemInfoReqs());

        List<FormItemTemplate> formItemTemplateList = new ArrayList<>();
        String msg = "success";
        formItemInfos.stream().forEach(f -> {
            QueryWrapper<FormItemTemplate> queryWrapper = new QueryWrapper<>();
            queryWrapper.eq("id", f.getItemTempId());
            formItemTemplateList.add(formItemTemplateMapper.selectOne(queryWrapper));
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
    }

    @Override
    public QueryResponse<SynthesisTaskInfoResp> selectSynthesisTaskInfoService(Integer page, Integer pageSize, SynthesisTaskInfoReq req) {
        String currentUserRolesToString = AuthenticationContextHolder.getCurrentUserRolesToString();
        req.setRoleName(currentUserRolesToString);
        List<Map<String,Object>> list = new ArrayList<>();
        IPage<TaskInfo> iPage = taskInfoMapper.selectSynthesisRequestInfo(new Page<TaskInfo>(page, pageSize),req);
        List<SynthesisTaskInfoResp> srt=synthesisTaskInfoRespConverter.toDto(iPage.getRecords());

        QueryResponse<SynthesisTaskInfoResp> queryResponse = new QueryResponse<>();
        queryResponse.setPageInfo(new PageInfo(iPage.getTotal(),iPage.getCurrent(),iPage.getSize()));
        queryResponse.setContents(srt);

        return queryResponse;
    }

    @Override
    public SynthesisTaskInfoFormTask selectSynthesisTaskInfoFormService(String id) throws Exception{
        FormInfo formInfo=formInfoMapper.selectOne(new QueryWrapper<FormInfo>().eq("record_id",id));
        if (StringUtils.isEmpty(id)){
            throw new Exception("The request details do not exist");
        }
        if(StringUtils.isEmpty((CharSequence) formInfo)){
            throw new Exception("Task information cannot be empty");
        }
        List<FormItemInfo> formItemInfos=formItemInfoMapper.selectList(new QueryWrapper<FormItemInfo>().eq("form_id",formInfo.getId()));
        SynthesisTaskInfoFormTask srt=synthesisTaskInfoFormTaskConverter.toDto(formInfo);
        srt.setFormItemInfo(formItemInfos);

        return srt;
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
}
