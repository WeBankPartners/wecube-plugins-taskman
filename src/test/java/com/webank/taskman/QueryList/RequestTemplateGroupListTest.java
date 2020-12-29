package com.webank.taskman.QueryList;

import com.webank.taskman.TmallApplicationTests;
import com.webank.taskman.base.JsonResponse;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.controller.x100.TaskmanRequestController;
import com.webank.taskman.dto.RequestTemplateGroupDTO;
import com.webank.taskman.service.RequestTemplateGroupService;
import com.webank.taskman.utils.GsonUtil;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit4.SpringRunner;


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
        queryResponse=requestTemplateGroupService.selectByParam(page, pageSize, req);
        for (RequestTemplateGroupDTO content : queryResponse.getContents()) {
            System.out.println(content.toString());
        }
    }

}
