package com.webank.taskman;

import org.junit.Test;
import org.springframework.beans.factory.annotation.Autowired;

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
