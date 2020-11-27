package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.domain.TaskTemplate;
import com.webank.taskman.mapper.TaskTemplateMapper;
import com.webank.taskman.service.TaskTemplateService;
import org.springframework.stereotype.Service;


@Service
public class TaskTemplateServiceImpl extends ServiceImpl<TaskTemplateMapper, TaskTemplate> implements TaskTemplateService {

}
