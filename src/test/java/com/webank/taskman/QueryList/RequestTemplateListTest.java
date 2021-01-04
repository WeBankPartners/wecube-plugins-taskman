package com.webank.taskman.QueryList;

import com.webank.taskman.TmallApplicationTests;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.dto.RequestTemplateDTO;
import com.webank.taskman.dto.req.QueryRequestTemplateReq;
import com.webank.taskman.service.RequestTemplateService;
import com.webank.taskman.utils.GsonUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.junit.Test;

public class RequestTemplateListTest extends TmallApplicationTests {
    private static final Logger log = LoggerFactory.getLogger(RequestTemplateListTest.class);

    @Autowired
    RequestTemplateService requestTemplateService;

    @Test
    public void RequestTemplateListTest(){
        Integer page=1;
        Integer pageSize=2;
        QueryRequestTemplateReq req=new QueryRequestTemplateReq();
        log.info("Received request parameters:{}", GsonUtil.GsonString(req) );
        QueryResponse<RequestTemplateDTO> queryResponse = requestTemplateService.selectRequestTemplatePage(page, pageSize, req);
        for (RequestTemplateDTO content : queryResponse.getContents()) {
        }
    }
}
