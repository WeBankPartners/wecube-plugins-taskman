package com.webank.taskman.QueryList;

import com.webank.taskman.TmallApplicationTests;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.dto.RequestTemplateGroupDTO;
import com.webank.taskman.service.RequestTemplateGroupService;
import com.webank.taskman.utils.GsonUtil;
import org.junit.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;


public class RequestTemplateGroupListTest extends TmallApplicationTests {

    private static final Logger log = LoggerFactory.getLogger(RequestTemplateGroupListTest.class);

    @Autowired
    RequestTemplateGroupService requestTemplateGroupService;

    @Test
    public void RequestTemplateGroupListTest(){
        Integer page=1;
        Integer pageSize=2;
        RequestTemplateGroupDTO req=new RequestTemplateGroupDTO();
        log.info("Received request parameters:{}", GsonUtil.GsonString(req) );
        QueryResponse<RequestTemplateGroupDTO> queryResponse=new QueryResponse<>();
        queryResponse=requestTemplateGroupService.selectRequestTemplateGroupPage(page, pageSize, req);
        for (RequestTemplateGroupDTO content : queryResponse.getContents()) {
        }
    }

}