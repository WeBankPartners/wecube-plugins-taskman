package com.webank.taskman.support.core;

import com.webank.taskman.commons.AppProperties;
import com.webank.taskman.commons.AppProperties.ServiceTaskmanProperties;
import org.springframework.beans.factory.annotation.Autowired;

public class newCoreServiceStub {


    private static final String GET_PROC_DEF_LIST = "/platform/v1//release/process/definitions";
    private static final String GET_PROC_NODE_LIST = "/release/process/definitions/{proc-def-id}/tasknodes";
    private static final String CREATE_NEW_WORKFLOW_INSTANCE = "/release/process/instances";

    @Autowired
    private CoreRestTemplate template;

    @Autowired
    private ServiceTaskmanProperties smProperties;


}
