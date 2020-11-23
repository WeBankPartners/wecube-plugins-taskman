package com.webank.taskman.controller;


import com.webank.taskman.mapper.TaskMapper;
import com.webank.taskman.service.TaskService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;

import org.springframework.web.bind.annotation.RestController;

/**
 * <p>
 *  前端控制器
 * </p>
 *
 * @author ${author}
 * @since 2020-11-23
 */
@RestController
@RequestMapping("/task")
public class TaskController {


    @Autowired
    TaskService taskService;

    @Autowired
    TaskMapper taskMapper;

    @GetMapping("/getAll")
    public Object getAll(){
        return taskService.list();
    }

    @GetMapping("/list")
    public Object list(){
        return taskMapper.selectList(null);
    }
}

