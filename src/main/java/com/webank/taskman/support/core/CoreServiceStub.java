package com.webank.taskman.support.core;

import com.webank.taskman.support.core.dto.CoreProcessDefinitionDto;
import com.webank.taskman.support.core.dto.RolesDataResponse;
import org.springframework.stereotype.Service;

import java.util.LinkedList;
import java.util.List;

@Service
public class CoreServiceStub {
    private static final int NOT_INCLUDE_DRAFT = 0;

    private static final String GET_ALL_ROLES = "/auth/v1/roles";
    private static final String GET_ROLES_BY_USER_NAME = "/auth/v1/users/%s/roles";
    private static final String REPORT_OPERATION_EVENTS = "/platform/v1/operation-events";
    private static final String GET_ALL_PEOCESS_KEYS = "/platform/v1/process/definitions?includeDraft=%d";

    private static final String GET_ROOT_ENTITIES = "/platform/v1/process/definitions/process-keys/%s/root-entities";

//    @Autowired
//    private CoreRestTemplate template;


    //1
    public List<RolesDataResponse> getAllRoles() {
        List<RolesDataResponse> rolesDataResponses = new LinkedList<>();
        RolesDataResponse response = new RolesDataResponse();
        response.setRoleId("2c9280827019695c017019ac974f001c");
        response.setRoleName("SUPER_ADMIN");
        response.setDescription("SUPER_ADMIN");
        rolesDataResponses.add(response);
//        return template.get(asCoreUrl(GET_ALL_ROLES), GetAllRolesResponse.class);
        return rolesDataResponses;
    }
    // 1
    public List<RolesDataResponse> getRolesByUserName(String userName) {
        List<RolesDataResponse> rolesDataResponses = new LinkedList<>();
        RolesDataResponse response = new RolesDataResponse();
        response.setRoleId("2c9280827019695c017019ac974f001c");
        response.setRoleName("SUPER_ADMIN");
        response.setDescription("SUPER_ADMIN");
        rolesDataResponses.add(response);

//        return template.get(asCoreUrl(GET_ROLES_BY_USER_NAME, userName), GetAllRolesResponse.class);
        return rolesDataResponses;
    }

//    private String asCoreUrl(String path, Object... pathVariables) {
//        if (pathVariables != null && pathVariables.length > 0) {
//            path = String.format(path, pathVariables);
//        }
//        return smProperties.getWecubeCoreAddress() + path;
//    }

    // 1
    public List<CoreProcessDefinitionDto> getAllProcessDefinitionKeys() {
        List<CoreProcessDefinitionDto> coreProcessDefinitionDtos = new LinkedList<>();

//        return template.get(asCoreUrl(GET_ALL_PEOCESS_KEYS, NOT_INCLUDE_DRAFT), GetAllProcessKeysResponse.class);
        return coreProcessDefinitionDtos;

    }

}
