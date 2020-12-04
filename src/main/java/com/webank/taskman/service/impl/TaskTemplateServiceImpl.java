package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.domain.TaskTemplate;
import com.webank.taskman.dto.req.SaveTaskTemplateReq;
import com.webank.taskman.mapper.TaskTemplateMapper;
import com.webank.taskman.service.TaskTemplateService;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.BeanUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Date;


@Service
public class TaskTemplateServiceImpl extends ServiceImpl<TaskTemplateMapper, TaskTemplate> implements TaskTemplateService {

    @Autowired
    TaskTemplateMapper taskTemplateMapper;

    @Override
    public TaskTemplate addOrUpdateTaskTemplate(SaveTaskTemplateReq taskTemplateReq) {
        TaskTemplate taskTemplate = new TaskTemplate();
        BeanUtils.copyProperties(taskTemplateReq, taskTemplate);
        if (StringUtils.isEmpty(taskTemplate.getId())) {
            taskTemplate.setCreatedBy("11");
            taskTemplate.setUpdatedBy("22");
            taskTemplateMapper.insert(taskTemplate);
            return taskTemplateMapper.selectById(taskTemplate);
        }
        if (!StringUtils.isEmpty(taskTemplate.getId())) {
            taskTemplate.setUpdatedTime(new Date());
            taskTemplateMapper.updateById(taskTemplate);
            return taskTemplateMapper.selectById(taskTemplate);
        }
        return null;
    }
}
