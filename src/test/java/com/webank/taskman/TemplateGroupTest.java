package com.webank.taskman;

<<<<<<<<< Temporary merge branch 1
import com.webank.taskman.dto.TemplateGroupVO;
import com.webank.taskman.service.TemplateGroupService;
=========
import com.webank.taskman.service.RequestTemplateGroupService;
import com.webank.taskman.support.core.CoreServiceStub;
import com.webank.taskman.support.core.dto.RolesDataResponse;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit4.SpringRunner;

import java.util.Collections;
import java.util.List;


@SpringBootTest
@RunWith(SpringRunner.class)
public class TemplateGroupTest {

    @Autowired
    private RequestTemplateGroupService requestTemplateGroupService;

    @Autowired
    CoreServiceStub coreServiceStub;

    @Test
    public void createRequestTemplateGroup() throws Exception {
        System.out.println(getRoles(1));
        for (int i=0;i<10;i++) {
            System.out.println(getRoles());
        }
        String request =
                "{" +
                "\"id\": \"\"," +
                "\"name\": \"\"," +
                "\"manageRoleId\": \"\"," +
                "\"manageRoleName\": \"\"," +
                "\"description\": \"\"," +
                "\"version\": \"\"" +
                "}";


    }

    private List<RolesDataResponse> getRoles(int limit){
        // get all roles
        List<RolesDataResponse> roles = coreServiceStub.authRoleAll();
        limit = limit > roles.size() ? roles.size():limit;
        Collections.shuffle(roles);
        return roles.subList(0,limit);
    }
    private List<RolesDataResponse> getRoles(){
        return  getRoles(Math.round(100));
    }
}
