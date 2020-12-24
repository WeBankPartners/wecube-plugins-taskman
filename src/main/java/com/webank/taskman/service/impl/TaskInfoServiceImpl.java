package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.github.pagehelper.PageHelper;
import com.github.pagehelper.PageInfo;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.converter.*;
import com.webank.taskman.domain.FormInfo;
import com.webank.taskman.domain.FormItemInfo;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.dto.CheckTaskDTO;
import com.webank.taskman.dto.TaskInfoDTO;
import com.webank.taskman.dto.req.QueryTaskInfoReq;
import com.webank.taskman.dto.req.SaveTaskInfoReq;
import com.webank.taskman.dto.req.SynthesisTaskInfoReq;
import com.webank.taskman.dto.resp.FormInfoResq;
import com.webank.taskman.dto.resp.SaveTaskInfoResp;
import com.webank.taskman.dto.resp.SynthesisTaskInfoFormTask;
import com.webank.taskman.dto.resp.SynthesisTaskInfoResp;
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
    public QueryResponse<TaskInfoDTO> selectTaskInfo(Integer page, Integer pageSize, QueryTaskInfoReq req) {
        String currentUserRolesToString = AuthenticationContextHolder.getCurrentUserRolesToString();
        req.setUseRoleName(currentUserRolesToString);
        PageHelper.startPage(page,pageSize);

        PageInfo<TaskInfo> pages = new PageInfo(taskInfoMapper.selectTaskInfo(req));
        QueryResponse<TaskInfoDTO> queryResponse = new QueryResponse(pages.getTotal(),page.longValue(),pageSize.longValue(),pages.getList());
        return queryResponse;
    }

    @Override
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
    }

    @Override
    public QueryResponse<SynthesisTaskInfoResp> selectSynthesisTaskInfoService(Integer page, Integer pageSize, SynthesisTaskInfoReq req) {
        String currentUserRolesToString = AuthenticationContextHolder.getCurrentUserRolesToString();
        req.setRoleName(currentUserRolesToString);
        IPage<TaskInfo> iPage = taskInfoMapper.selectSynthesisRequestInfo(new Page<TaskInfo>(page, pageSize),req);
        List<SynthesisTaskInfoResp> srt=synthesisTaskInfoRespConverter.toDto(iPage.getRecords());

        QueryResponse<SynthesisTaskInfoResp> queryResponse = new QueryResponse<>();
        queryResponse.setPageInfo(new com.webank.taskman.base.PageInfo(iPage.getTotal(),iPage.getCurrent(),iPage.getSize()));
        queryResponse.setContents(srt);

        return queryResponse;
    }

    @Override
    public SynthesisTaskInfoFormTask selectSynthesisTaskInfoFormService(String id) throws Exception{
        FormInfo formInfo=formInfoMapper.selectOne(new FormInfo().setRecordId(id).getLambdaQueryWrapper());
        if (StringUtils.isEmpty(id)){
            throw new Exception("The request details do not exist");
        }
        if(StringUtils.isEmpty((CharSequence) formInfo)){
            throw new Exception("Task information cannot be empty");
        }
        List<FormItemInfo> formItemInfos=formItemInfoMapper.selectList(new FormItemInfo().setFormId(formInfo.getId()).getLambdaQueryWrapper());
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
