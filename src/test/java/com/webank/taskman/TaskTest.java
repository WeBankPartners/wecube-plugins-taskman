package com.webank.taskman;

import com.webank.taskman.domain.Task;
import com.webank.taskman.mapper.TaskMapper;
import com.webank.taskman.service.TaskService;
import org.junit.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;

import java.util.List;

public class TaskTest extends TmallApplicationTests{

    @Autowired
    TaskService taskService;

    @Autowired
    TaskMapper taskMapper;

    @Test
    public void testSelect() {
        System.out.println(("----- selectAll method test ------"));
        List<Task> taskList = taskMapper.selectList(null);
        for(Task task:taskList) {
            System.out.println(task);
        }
    }
}
