package com.webank.taskman.ProcessTesting;

import com.webank.taskman.TmallApplicationTests;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.dto.RequestTemplateDTO;
import com.webank.taskman.dto.req.QueryRequestTemplateReq;
import com.webank.taskman.dto.req.SaveRequestTemplateReq;
import com.webank.taskman.service.RequestTemplateService;
import org.junit.Test;
import org.springframework.beans.factory.annotation.Autowired;

public class TemplateTest extends TmallApplicationTests {
    @Autowired
    RequestTemplateService requestTemplateService;

    @Test
    public void selectTemplateTest(){
        Integer page =1;
        Integer pageSize=10;
        QueryRequestTemplateReq req=new QueryRequestTemplateReq();
        req.setName("");
        req.setManageRoleName("");
        req.setUseRoleName("");
        req.setTags("");
        QueryResponse<RequestTemplateDTO> queryResponse = requestTemplateService.selectRequestTemplatePage(page, pageSize, req);
        for (RequestTemplateDTO content : queryResponse.getContents()) {
        }
    }

    @Test
    public void AddTemplateTest(){
        SaveRequestTemplateReq req=new SaveRequestTemplateReq();
        RequestTemplateDTO requestTemplateDTO = requestTemplateService.saveRequestTemplate(req);
    }
}
