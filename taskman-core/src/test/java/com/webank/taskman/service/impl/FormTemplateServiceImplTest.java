package com.webank.taskman.service.impl;

import org.junit.Ignore;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit4.SpringRunner;

import com.webank.taskman.dto.req.FormTemplateSaveReqDto;
import com.webank.taskman.dto.resp.FormTemplateQueryResultDto;
import com.webank.taskman.service.FormTemplateService;


@Ignore
@RunWith(SpringRunner.class)
@SpringBootTest
public class FormTemplateServiceImplTest {
    
    @Autowired
    FormTemplateService formTemplateService;

    @Test
    public void testDetailFormTemplate() {
        
        FormTemplateSaveReqDto req = new FormTemplateSaveReqDto();
        req.setTempId("1339560042995204098");
        req.setTempType("0");
        
        FormTemplateQueryResultDto respDto = formTemplateService.detailFormTemplate(req);
        
        System.out.println(respDto);
    }

}
