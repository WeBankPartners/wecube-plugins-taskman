package com.webank.taskman.ProcessTesting;

import com.webank.taskman.TmallApplicationTests;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.TaskmanException;
import com.webank.taskman.dto.RequestTemplateGroupDTO;
import com.webank.taskman.dto.req.SaveRequestTemplateGropReq;
import com.webank.taskman.service.RequestTemplateGroupService;
import org.junit.Test;
import org.springframework.beans.factory.annotation.Autowired;

import java.util.Scanner;

public class TemplateGroupTest extends TmallApplicationTests {


    @Autowired
    RequestTemplateGroupService requestTemplateGroupService;

    @Test
    public void selectTemplateGroupTest(){
        Integer page =1;
        Integer pageSize=10;
        RequestTemplateGroupDTO req=new RequestTemplateGroupDTO();
        QueryResponse<RequestTemplateGroupDTO> queryResponse=new QueryResponse<>();
        queryResponse=requestTemplateGroupService.selectByParam(page, pageSize, req);
        for (RequestTemplateGroupDTO content : queryResponse.getContents()) {
            System.out.println(content.toString());
        }

    }


    @Test
    public void AddTemplateGroupTest() throws TaskmanException {
        SaveRequestTemplateGropReq req=new SaveRequestTemplateGropReq();
        req.setManageRoleId("SUPER_ADMIN");
        req.setManageRoleName("SUPER_ADMIN");
        req.setName("测试2");
        req.setDescription("测试描述");
        RequestTemplateGroupDTO response=new RequestTemplateGroupDTO();
        response=requestTemplateGroupService.saveTemplateGroupByReq(req);
        System.out.println(response.toString());

    }


}
