package com.webank.taskman;

import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.SaveFormTemplateReq;
import com.webank.taskman.dto.resp.FormTemplateResp;
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
        saveFormTemplateReq.setTargetEntitys("");
        formTemplateService.saveFormTemplate(saveFormTemplateReq);
    }

    //Query single data details of form module according to ID
    @Test
    public void detailFormTemplateTest() throws Exception {
        SaveFormTemplateReq saveFormTemplateReq=new SaveFormTemplateReq();
        saveFormTemplateReq.setId("");
        saveFormTemplateReq.setTempId("");
        saveFormTemplateReq.setTempType("");
        saveFormTemplateReq.setName("");
        saveFormTemplateReq.setDescription("");
        saveFormTemplateReq.setStyle("");
        saveFormTemplateReq.setTargetEntitys("");
        formTemplateService.detailFormTemplate(saveFormTemplateReq);
    }


    //Delete single data of form module according to ID logic
    @Test
    public void deleteFormTemplateTest() throws Exception {
        String id="";
        formTemplateService.deleteFormTemplate(id);
    }

    //Paging query form module data
    @Test
    public void selectFormTemplateTest() throws Exception {
        int current=1;
        int limit=2;
        SaveFormTemplateReq saveFormTemplateReq=new SaveFormTemplateReq();
        saveFormTemplateReq.setId("");
        saveFormTemplateReq.setTempId("");
        saveFormTemplateReq.setTempType("");
        saveFormTemplateReq.setName("");
        saveFormTemplateReq.setDescription("");
        saveFormTemplateReq.setStyle("");
        saveFormTemplateReq.setTargetEntitys("");
        QueryResponse<FormTemplateResp> queryResponse= formTemplateService.selectFormTemplate(current,limit,saveFormTemplateReq);
        System.out.println(JsonResponse.okayWithData(queryResponse).getData());
    }

}
