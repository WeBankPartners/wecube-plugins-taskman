package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.converter.FormInfoConverter;
import com.webank.taskman.converter.TaskInfoConverter;
import com.webank.taskman.domain.FormInfo;
import com.webank.taskman.domain.FormItemInfo;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.dto.PageInfo;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.SaveTaskInfoReq;
import com.webank.taskman.dto.resp.FormInfoResq;
import com.webank.taskman.dto.resp.TaskInfoResp;
import com.webank.taskman.mapper.FormInfoMapper;
import com.webank.taskman.mapper.FormItemInfoMapper;
import com.webank.taskman.mapper.TaskInfoMapper;
import com.webank.taskman.service.TaskInfoService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;


@Service
public class TaskInfoServiceImpl extends ServiceImpl<TaskInfoMapper, TaskInfo> implements TaskInfoService {

    @Autowired
    TaskInfoMapper taskInfoMapper;

    @Autowired
    TaskInfoConverter taskInfoConverter;

    @Autowired
    FormInfoMapper formInfoMapper;

    @Autowired
    FormInfoConverter formInfoConverter;

    @Autowired
    FormItemInfoMapper formItemInfoMapper;


    @Override
    public QueryResponse<TaskInfoResp> selectTaskInfoService(Integer page, Integer pageSize, SaveTaskInfoReq req) {
        IPage<TaskInfo> iPage= taskInfoMapper.selectTaskInfo(new Page<>(page, pageSize),req);
        List<TaskInfoResp> respList=taskInfoConverter.toDto(iPage.getRecords());
        for (TaskInfoResp taskInfoResp : respList) {
            FormInfo formInfo=formInfoMapper.selectOne(new QueryWrapper<FormInfo>().eq("record_id",taskInfoResp.getId()));
            FormInfoResq formInfoResq=formInfoConverter.toDto(formInfo);
            if (formInfoResq!=null) {
                formInfoResq.setFormItemInfo(formItemInfoMapper.selectFormItemInfo(taskInfoResp.getId()));
                taskInfoResp.setFormInfoResq(formInfoResq);
            }
        }

        QueryResponse<TaskInfoResp> queryResponse=new QueryResponse<>();
        queryResponse.setPageInfo(new PageInfo(iPage.getTotal(),iPage.getCurrent(),iPage.getSize()));
        queryResponse.setContents(respList);
        return queryResponse;
    }
}
