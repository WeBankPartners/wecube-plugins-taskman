package com.webank.taskman;

import com.webank.taskman.dto.TemplateGroupVO;
import com.webank.taskman.service.RequestTemplateGroupService;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit4.SpringRunner;


@SpringBootTest
@RunWith(SpringRunner.class)
public class TemplateGroupTest {
    @Autowired
    private RequestTemplateGroupService requestTemplateGroupService;

    @Test
    public void testV1Group() throws Exception {
//        TemplateGroupVO vo=new TemplateGroupVO();
//        vo.setCreatedBy("11");
//        vo.setUpdatedBy("22");
//        vo.setManageRole("11");
//        vo.setName("11");
//        vo.setVersion("11");
//        templateGroupService.createTemplateGroupService(vo);
        TemplateGroupVO vo=new TemplateGroupVO();
        vo.setId("1332131518222503937");
        vo.setCreatedBy("问我");
        vo.setName("大答");

        requestTemplateGroupService.updateTemplateGroupService(vo);
    }
}
