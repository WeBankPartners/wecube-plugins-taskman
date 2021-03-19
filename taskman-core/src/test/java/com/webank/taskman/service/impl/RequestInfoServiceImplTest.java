package com.webank.taskman.service.impl;

import java.util.ArrayList;
import java.util.List;

import org.junit.Assert;
import org.junit.Ignore;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit4.SpringRunner;

import com.webank.taskman.dto.CreateTaskDto;
import com.webank.taskman.dto.CreateTaskDto.EntityValueDto;
import com.webank.taskman.dto.resp.RequestInfoResqDto;
import com.webank.taskman.service.RequestInfoService;

@RunWith(SpringRunner.class)
@SpringBootTest
public class RequestInfoServiceImplTest {
    
    @Autowired
    RequestInfoService requestInfoService;

    @Ignore
    @Test
    public void testCreateNewRequestInfo() {
        
        CreateTaskDto reqDto = mockCreateTaskDto();
        
        
        RequestInfoResqDto respDto = requestInfoService.createNewRequestInfo(reqDto);
        
        
        Assert.assertNotNull(respDto);
    }
    
    
    private CreateTaskDto mockCreateTaskDto(){
        CreateTaskDto reqDto = new CreateTaskDto();
        reqDto.setName("req-test-001");
        reqDto.setRequestTempId("1339506042731552770");
        reqDto.setEmergency("Normal");
        reqDto.setRootEntity("0001_11112222");
        
        List<EntityValueDto> entities = new ArrayList<>();
        reqDto.setEntities(entities);
        
        return reqDto;
    }

}
