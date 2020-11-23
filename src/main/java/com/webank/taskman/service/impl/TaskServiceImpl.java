package com.webank.taskman.service.impl;

import com.webank.taskman.domain.Task;
import com.webank.taskman.mapper.TaskMapper;
import com.webank.taskman.service.TaskService;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.stereotype.Service;

/**
 * <p>
 *  服务实现类
 * </p>
 *
 * @author ${author}
 * @since 2020-11-23
 */
@Service
public class TaskServiceImpl extends ServiceImpl<TaskMapper, Task> implements TaskService {

}
