package com.webank.taskman.support.core;

import com.webank.taskman.commons.AppProperties.ServiceTaskmanProperties;
import com.webank.taskman.support.core.dto.CoreProcessDefinitionDto;
import com.webank.taskman.support.core.dto.CoreResponse.*;
import com.webank.taskman.support.core.dto.RolesDataResponse;
import com.webank.taskman.utils.JsonUtils;
import com.webank.taskman.utils.SpringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.io.IOException;
import java.util.LinkedList;
import java.util.List;

@Service
public class CoreServiceStub {
    private static final int NOT_INCLUDE_DRAFT = 0;

    private static final String GET_ALL_ROLES = "/auth/v1/roles";
    private static final String GET_ROLES_BY_USER_NAME = "/auth/v1/users/%s/roles";
    private static final String REPORT_OPERATION_EVENTS = "/platform/v1/operation-events";
    private static final String GET_ALL_PEOCESS_KEYS = "/platform/v1/process/definitions?includeDraft=%d";

    private static final String GET_ROOT_ENTITIES = "/platf@Api(tags = {\"1、 CoreResource model\"})\n" +
            "@RestController\n" +
            "@RequestMapping(\"/v1/core-resources\")orm/v1/process/definitions/process-keys/%s/root-entities";

    @Autowired
    private CoreRestTemplate template;

    @Autowired
    private ServiceTaskmanProperties smProperties;
    //1
    public List<RolesDataResponse> getAllRoles() {
        if("dev".equals(SpringUtils.getActiveProfile())){
            return addRoles();
        }
        return template.get(asCoreUrl(GET_ALL_ROLES), GetAllRolesResponse.class);
    }
    // 1
    public List<RolesDataResponse> getRolesByUserName(String userName) {
        if("dev".equals(SpringUtils.getActiveProfile())){
            return addRoleTestData();
        }
        return template.get(asCoreUrl(GET_ROLES_BY_USER_NAME, userName), GetAllRolesResponse.class);
    }

    private String asCoreUrl(String path, Object... pathVariables) {
        if (pathVariables != null && pathVariables.length > 0) {
            path = String.format(path, pathVariables);
        }
        return smProperties.getWecubeCoreAddress() + path;
    }

    // 1
    public List<CoreProcessDefinitionDto> getAllProcessDefinitionKeys()  {
        if("dev".equals(SpringUtils.getActiveProfile())){
            return addPdfTestData();
        }
        return template.get(asCoreUrl(GET_ALL_PEOCESS_KEYS, NOT_INCLUDE_DRAFT), GetAllProcessKeysResponse.class);
    }

    private List<RolesDataResponse> addRoleTestData(){
        List<RolesDataResponse> rolesDataResponses = new LinkedList<>();
        RolesDataResponse response = new RolesDataResponse();
        response.setRoleId("2c9280827019695c017019ac974f001c");
        response.setRoleName("SUPER_ADMIN");
        response.setDescription("SUPER_ADMIN");
        rolesDataResponses.add(response);
        return rolesDataResponses;
    }
    private List<RolesDataResponse> addRoles(){
        String json = "[{\"id\":\"2c9280827019695c017019ac974f001c\",\"name\":\"SUPER_ADMIN\",\"displayName\":\"SUPER_ADMIN\",\"email\":null},{\"id\":\"2c9280836f78a84b016f794c3a270000\",\"name\":\"CMDB_ADMIN\",\"displayName\":\"CMDB管理员\",\"email\":null},{\"id\":\"2c9280836f78a84b016f794cd6dd0001\",\"name\":\"MONITOR_ADMIN\",\"displayName\":\"监控管理员\",\"email\":null},{\"id\":\"2c9280836f78a84b016f794d6bb50002\",\"name\":\"PRD_OPS\",\"displayName\":\"生产运维\",\"email\":null},{\"id\":\"2c9280836f78a84b016f794e0d3b0003\",\"name\":\"STG_OPS\",\"displayName\":\"测试运维\",\"email\":null},{\"id\":\"2c9280836f78a84b016f794e9b170004\",\"name\":\"APP_ARC\",\"displayName\":\"应用架构师\",\"email\":null},{\"id\":\"2c9280836f78a84b016f794f20440005\",\"name\":\"IFA_ARC\",\"displayName\":\"基础架构师\",\"email\":null},{\"id\":\"2c9280836f78a84b016f794ff45e0006\",\"name\":\"APP_DEV\",\"displayName\":\"应用开发人员\",\"email\":null},{\"id\":\"2c9280836f78a84b016f795068870007\",\"name\":\"IFA_OPS\",\"displayName\":\"基础架构运维人员\",\"email\":null},{\"id\":\"8ab86ba0723a78fe01723a790ceb0000\",\"name\":\"SUB_SYSTEM\",\"displayName\":\"后台系统\",\"email\":null}]";
        List<RolesDataResponse> roles = new LinkedList<>();
        try {
            roles = JsonUtils.toObject(json, roles.getClass());
        } catch (IOException e) {
            e.printStackTrace();
        }
        return roles;
    }
    private List<CoreProcessDefinitionDto> addPdfTestData()  {
        String json = "[{\"procDefId\":\"shKXjq8Z2D0\",\"procDefKey\":\"wecube1584964874321\",\"procDefName\":\"TXCLOUD_应用系统部署全流程_首次_V1.0\",\"procDefVersion\":\"1\",\"status\":\"deployed\",\"procDefData\":null,\"rootEntity\":\"wecmdb:app_system\",\"createdTime\":\"2020-11-30 02:31:42\",\"permissionToRole\":null,\"taskNodeInfos\":[]},{\"procDefId\":\"shKXHPUZ2Fr\",\"procDefKey\":\"wecube1584956874532\",\"procDefName\":\"TXCLOUD_数据中心网络及结构初始化_V1.1\",\"procDefVersion\":\"1\",\"status\":\"deployed\",\"procDefData\":null,\"rootEntity\":\"wecmdb:data_center{data_center_type eq 'REGION'}\",\"createdTime\":\"2020-11-30 02:33:16\",\"permissionToRole\":null,\"taskNodeInfos\":[]},{\"procDefId\":\"shKXMsXZ2Gi\",\"procDefKey\":\"wecube1604319943267\",\"procDefName\":\"应用部署单元升级部署_V0.1\",\"procDefVersion\":\"1\",\"status\":\"deployed\",\"procDefData\":null,\"rootEntity\":\"wecmdb:unit\",\"createdTime\":\"2020-11-30 02:33:34\",\"permissionToRole\":null,\"taskNodeInfos\":[]},{\"procDefId\":\"shKXPkDZ2GA\",\"procDefKey\":\"wecube1584968536974\",\"procDefName\":\"默认告警处理编排_V0.2\",\"procDefVersion\":\"1\",\"status\":\"deployed\",\"procDefData\":null,\"rootEntity\":\"wecube-monitor:alarm\",\"createdTime\":\"2020-11-30 02:33:45\",\"permissionToRole\":null,\"taskNodeInfos\":[]},{\"procDefId\":\"shKXS9cZ2GO\",\"procDefKey\":\"wecube1584968688912\",\"procDefName\":\"默认告警解除编排_V0.2\",\"procDefVersion\":\"1\",\"status\":\"deployed\",\"procDefData\":null,\"rootEntity\":\"wecube-monitor:alarm\",\"createdTime\":\"2020-11-30 02:33:55\",\"permissionToRole\":null,\"taskNodeInfos\":[]}]";
        List<CoreProcessDefinitionDto> pefList = new LinkedList<>();
        try {
            pefList = JsonUtils.toObject(json, pefList.getClass());
        } catch (IOException e) {
            e.printStackTrace();
        }
        return pefList;
    }

}
