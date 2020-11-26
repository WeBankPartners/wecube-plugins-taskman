package com.webank.taskman;

import com.webank.taskman.converter.TaskConverter;
import com.webank.taskman.converter.TaskConverterImpl;
import org.junit.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;

import java.util.List;

@SpringBootTest(classes = {TaskConverterImpl.class,})
public class TaskTest {

    @Autowired
    TaskService taskService;

    @Autowired
    TaskMapper taskMapper;

    @Autowired
    TaskConverter taskConverter;

    @Test
    public void testSelect() {
        System.out.println(("----- selectList method test ------"));
        List<Task> taskList = taskMapper.selectList(null);
        for(Task task:taskList) {
            System.out.println(task);
        }
        taskConverter.toDto(taskList).forEach(e->{
            System.out.println(e);
        });
    }
}
