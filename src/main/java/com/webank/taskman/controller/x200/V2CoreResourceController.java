package com.webank.taskman.controller.x200;

import com.fasterxml.jackson.core.JsonParseException;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.support.core.CoreServiceStub;
import com.webank.taskman.support.core.dto.RolesDataResponse;
import com.webank.taskman.support.core.dto.WorkflowDefInfoDto;
import io.swagger.annotations.Api;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import javax.servlet.http.HttpServletRequest;
import java.io.IOException;
import java.util.List;

import static com.webank.taskman.dto.JsonResponse.okayWithData;


@Api(tags = {"1、 CoreResource inteface API"})
@RestController
@RequestMapping("/v2/core-resources")
public class V2CoreResourceController {

    @Autowired
    CoreServiceStub coreServiceStub;

    @GetMapping("/users/current-user/roles")
    public JsonResponse<List<RolesDataResponse>> getRolesByCurrentUser(HttpServletRequest httpRequest) {
        String currentUserName = AuthenticationContextHolder.getCurrentUsername();
        return okayWithData(coreServiceStub.getRolesByUserName(currentUserName));
    }

    @GetMapping("/roles")
    public JsonResponse<List<RolesDataResponse>> getAllRoles() throws JsonParseException, JsonMappingException, IOException {
        return okayWithData(coreServiceStub.getAllRoles());
    }

    @GetMapping("/workflow/process-definition-keys")
    public JsonResponse<List<WorkflowDefInfoDto>> getAllProcessDefinitionKeys() {
        return okayWithData(coreServiceStub.getAllProcessDefinitionKeys());
    }

    @GetMapping("/form-item/ci-data")
    public JsonResponse getCoreCiData() {
        return okayWithData(coreServiceStub.getAllProcessDefinitionKeys());
    }


}
