package com.webank.taskman.service.impl;

import com.webank.taskman.domain.TaskInfo;
import com.webank.taskman.mapper.TaskInfoMapper;
import com.webank.taskman.service.TaskInfoService;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.stereotype.Service;

/**
 * <p>
 * 任务记录表  服务实现类
 * </p>
 *
 * @author ${author}
 * @since 2020-11-26
 */
@Service
public class TaskInfoServiceImpl extends ServiceImpl<TaskInfoMapper, TaskInfo> implements TaskInfoService {

}
