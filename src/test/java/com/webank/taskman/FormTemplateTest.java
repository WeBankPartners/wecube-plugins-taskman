package com.webank.taskman;

import com.webank.taskman.dto.req.SaveFormTemplateReq;
import com.webank.taskman.service.FormTemplateService;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit4.SpringRunner;

@SpringBootTest
@RunWith(SpringRunner.class)
public class FormTemplateTest {
    @Autowired
    FormTemplateService formTemplateService;

    //Insert form module data
    @Test
    public void saveFormTemplateTest(){
        SaveFormTemplateReq saveFormTemplateReq=new SaveFormTemplateReq();
        saveFormTemplateReq.setId("");
        saveFormTemplateReq.setTempId("");
        saveFormTemplateReq.setTempType("");
        saveFormTemplateReq.setName("");
        saveFormTemplateReq.setDescription("");
        saveFormTemplateReq.setStyle("");
        saveFormTemplateReq.setCreatedBy("");
        saveFormTemplateReq.setUpdatedBy("");
        formTemplateService.saveFormTemplate(saveFormTemplateReq);
    }

    //
    @Test
    public void detailFormTemplateTest() throws Exception {
        String id="";
        formTemplateService.detailFormTemplate(id);
    }


    @Test
    public void deleteFormTemplateTest() throws Exception {
        String id="";
        formTemplateService.deleteFormTemplate(id);
    }



}
