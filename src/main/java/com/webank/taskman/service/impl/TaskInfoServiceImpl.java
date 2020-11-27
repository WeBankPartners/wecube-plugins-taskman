package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.mapper.TaskInfoMapper;
import com.webank.taskman.service.TaskInfoService;
import org.springframework.stereotype.Service;


@Service
public class TaskInfoServiceImpl extends ServiceImpl<TaskInfoMapper, TaskInfo> implements TaskInfoService {

}
