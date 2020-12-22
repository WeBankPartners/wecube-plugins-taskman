package com.webank.taskman;

import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.req.QueryRequestTemplateReq;
import com.webank.taskman.dto.resp.RequestTemplateResp;
import com.webank.taskman.mapper.RequestTemplateMapper;
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

//    @Test
//    @Ignore
//    public void saveRequestTemplateTest(){
//        SaveRequestTemplateReq saveRequestTemplateReq=new SaveRequestTemplateReq();
//        saveRequestTemplateReq.setId("");
//        saveRequestTemplateReq.setRequestTempGroup("");
//        saveRequestTemplateReq.setProcDefKey("");
//        saveRequestTemplateReq.setProcDefId("");
//        saveRequestTemplateReq.setProcDefName("");
//        saveRequestTemplateReq.setName("");
//        saveRequestTemplateReq.setTags("");
//        requestTemplateService.saveRequestTemplate(saveRequestTemplateReq);
//    }
//
//    //Query single data details of form module according to ID
//    @Test
//    @Ignore
//    public void detailRequestTemplateTest() throws Exception {
//        String id="";
//        requestTemplateService.detailRequestTemplate(id);
//    }
//
//
//    //Delete single data of form module according to ID logic
//    @Test
//    @Ignore
//    public void deleteRequestTemplateTest() throws Exception {
//        String id="";
//        requestTemplateService.deleteRequestTemplateService(id);
//    }

    @Autowired
    RequestTemplateMapper requestTemplateMapper;
    //Paging query form module data
    @Test
    public void selectRequestTemplateTest() throws Exception {
        requestTemplateMapper.selectList(new RequestTemplate().getLambdaQueryWrapper());
        int current=1;
        int limit=2;
        QueryRequestTemplateReq req=new QueryRequestTemplateReq();
        QueryResponse<RequestTemplateResp> queryResponse = requestTemplateService.selectRequestTemplatePage(current, limit, req);
    }
}
