package com.webank.taskman;

import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.QueryRequestTemplateReq;
import com.webank.taskman.dto.req.SaveRequestTemplateReq;
import com.webank.taskman.dto.resp.RequestTemplateResp;
import com.webank.taskman.service.RequestTemplateService;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit4.SpringRunner;

@SpringBootTest
@RunWith(SpringRunner.class)
public class RequestTemplateTest {

    @Autowired
    RequestTemplateService requestTemplateService;

    @Test
    public void saveRequestTemplateTest(){
        SaveRequestTemplateReq saveRequestTemplateReq=new SaveRequestTemplateReq();
        saveRequestTemplateReq.setId("");
        saveRequestTemplateReq.setRequestTempGroup("");
        saveRequestTemplateReq.setProcDefKey("");
        saveRequestTemplateReq.setProcDefId("");
        saveRequestTemplateReq.setProcDefName("");
        saveRequestTemplateReq.setName("");
        saveRequestTemplateReq.setTags("");
        requestTemplateService.saveRequestTemplate(saveRequestTemplateReq);
    }

    //Query single data details of form module according to ID
    @Test
    public void detailRequestTemplateTest() throws Exception {
        String id="";
        requestTemplateService.detailRequestTemplate(id);
    }


    //Delete single data of form module according to ID logic
    @Test
    public void deleteRequestTemplateTest() throws Exception {
        String id="";
        requestTemplateService.deleteRequestTemplateService(id);
    }

    //Paging query form module data
    @Test
    public void selectRequestTemplateTest() throws Exception {
        int current=1;
        int limit=2;
        QueryRequestTemplateReq req=new QueryRequestTemplateReq();
        req.setId("");
        req.setRequestTempGroup("");
        req.setProcDefKey("");
        req.setProcDefId("");
        req.setProcDefName("");
        req.setName("");
        req.setTags("");
        QueryResponse<RequestTemplateResp> queryResponse = requestTemplateService.selectRequestTemplatePage(current, limit, req);
    }
}
