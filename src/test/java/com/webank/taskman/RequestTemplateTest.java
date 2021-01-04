package com.webank.taskman;

import com.github.pagehelper.PageHelper;
import com.github.pagehelper.PageInfo;
import com.webank.taskman.dto.RequestTemplateDTO;
import com.webank.taskman.dto.req.QueryRequestTemplateReq;
import com.webank.taskman.mapper.RequestTemplateMapper;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit4.SpringRunner;

import java.util.List;

@SpringBootTest
@RunWith(SpringRunner.class)
public class RequestTemplateTest {

    @Autowired
    RequestTemplateMapper requestTemplateMapper;


    //Paging query form module data
    @Test
    public void selectRequestTemplateTest() throws Exception {
        PageHelper.startPage(1,10);
        List<RequestTemplateDTO> list = requestTemplateMapper.selectDTOListByParam(new QueryRequestTemplateReq());
        PageInfo page = new PageInfo(list);
    }
}
