package com.webank.taskman;

import com.webank.taskman.domain.TemplateGroup;
import com.webank.taskman.dto.TemplateGroupVO;
import com.webank.taskman.service.TemplateGroupService;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit4.SpringRunner;


@SpringBootTest
@RunWith(SpringRunner.class)
public class TemplateGroupTest {
    @Autowired
    private TemplateGroupService templateGroupService;

    @Test
    public void testV1Group() throws Exception {
        TemplateGroupVO vo=new TemplateGroupVO();
        vo.setCreatedBy("11");
        vo.setUpdatedBy("22");
        templateGroupService.createTemplateGroupService(vo);
    }
}
