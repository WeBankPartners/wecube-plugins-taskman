package com.webank.taskman;

import com.google.gson.reflect.TypeToken;
import com.webank.taskman.dto.RequestTemplateGroupDTO;
import com.webank.taskman.dto.req.SaveRequestTemplateGropReq;
import com.webank.taskman.service.RequestTemplateGroupService;
import com.webank.taskman.support.core.CoreServiceStub;
import com.webank.taskman.support.core.dto.RolesDataResponse;
import com.webank.taskman.utils.GsonUtil;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit4.SpringRunner;

import java.util.Collections;
import java.util.List;


@SpringBootTest
@RunWith(SpringRunner.class)
public class TemplateGroupTest {

    private static final Logger log = LoggerFactory.getLogger(TemplateGroupTest.class);

    @Autowired
    private RequestTemplateGroupService requestTemplateGroupService;

    @Autowired
    CoreServiceStub coreServiceStub;

    @Test
    public void createRequestTemplateGroup() throws Exception {
        RolesDataResponse role = getRole();
        String request ="{\"id\": \"%s\",\"name\": \"%s\",\"manageRoleId\": \"%s\",\"manageRoleName\": \"%s\",\"description\": \"%s\",\"version\": \"%s\"}";
        request = String.format(request,"","应用部署模板组",role.getRoleName(),role.getDescription(),"应用部署","1.0");
        SaveRequestTemplateGropReq req = GsonUtil.toObject(request,new TypeToken<SaveRequestTemplateGropReq>(){});
        try {
            RequestTemplateGroupDTO dto =  requestTemplateGroupService.saveTemplateGroupByReq(req);
        }catch (Exception e){
            log.error("创建任务失败：",e.getMessage());
        }
    }

    private RolesDataResponse getRoles(int limit){
        // get all roles
        List<RolesDataResponse> roles = coreServiceStub.authRoleAll();

        limit = limit > roles.size() ? roles.size():limit;
        Collections.shuffle(roles);
        return  roles.get(limit);
    }
    private RolesDataResponse getRole(){
        return  getRoles((int)Math.random()*10);
    }
}
