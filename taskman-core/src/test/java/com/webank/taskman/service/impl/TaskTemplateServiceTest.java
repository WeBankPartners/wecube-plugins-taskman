package com.webank.taskman.service.impl;

import static org.junit.Assert.fail;

import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit4.SpringRunner;

import com.webank.taskman.service.TaskTemplateService;

@RunWith(SpringRunner.class)
@SpringBootTest
public class TaskTemplateServiceTest {
    
    @Autowired
    TaskTemplateService taskTemplateService;

    @Test
    public void testSaveTaskTemplateByReq() {
        fail("Not yet implemented");
    }

}
