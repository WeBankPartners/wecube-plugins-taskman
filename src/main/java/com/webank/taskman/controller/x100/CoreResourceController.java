package com.webank.taskman.controller.x100;

import com.fasterxml.jackson.core.JsonParseException;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.support.core.CoreServiceStub;
import com.webank.taskman.support.core.dto.*;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiOperation;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import javax.servlet.http.HttpServletRequest;
import java.io.IOException;
import java.util.List;

import static com.webank.taskman.dto.JsonResponse.okayWithData;


@Api(tags = {"1、 CoreResource inteface API"})
@RestController
@RequestMapping("/v1/core-resources")
public class CoreResourceController {

    @Autowired
    CoreServiceStub coreServiceStub;

    @ApiOperationSupport(order = 1)
    @GetMapping("/users/current-user/roles")
    @ApiOperation(value = "core-resources-role-current-user", notes = "")
    public JsonResponse<List<RolesDataResponse>> getRolesByCurrentUser(HttpServletRequest httpRequest) {
        String currentUserName = AuthenticationContextHolder.getCurrentUsername();
        return okayWithData(coreServiceStub.getRolesByUserName(currentUserName));
    }

    @ApiOperationSupport(order = 2)
    @GetMapping("/roles")
    @ApiOperation(value = "core-resources-role-all", notes = "")
    public JsonResponse<List<RolesDataResponse>> getAllRoles() throws JsonParseException, JsonMappingException, IOException {
        return okayWithData(coreServiceStub.getAllRoles());
    }

    @ApiOperationSupport(order = 3)
    @GetMapping("/workflow/process-definition-keys")
    @ApiOperation(value = "workflow-process-all", notes = "")
    public JsonResponse<List<WorkflowDefInfoDto>> fetchLatestReleasedWorkflowDefs() {
        return okayWithData(coreServiceStub.fetchLatestReleasedWorkflowDefs());
    }

    //    @GetMapping("/workflow/process/definitions/{proc-def-id}/tasknodes")
    @ApiOperationSupport(order = 4)
    @GetMapping("/workflow/process-definitions-nodes/{proc-def-id}")
    @ApiOperation(value = "workflow-process-nodes", notes = "")
    public JsonResponse<List<WorkflowNodeDefInfoDto>> getTaskNodes(@PathVariable("proc-def-id") String procDefId) {
        return okayWithData(coreServiceStub.fetchWorkflowTasknodeInfos(procDefId));
    }

    @ApiOperationSupport(order = 5)
    @PostMapping("/workflow/process/crate")
    @ApiOperation(value = "workflow-process-crate", notes = "")
    public JsonResponse<CoreResponse.DynamicWorkflowInstInfoDto> createNewWorkflowInstance(
            @RequestBody DynamicWorkflowInstCreationInfoDto creationInfoDto)
    {
        return okayWithData(coreServiceStub.createNewWorkflowInstance(creationInfoDto));
    }

    /*@ApiOperationSupport(order = 6)
    @GetMapping("/form-item/ci-data")
    @ApiOperation(value = "workflow-process-crate", notes = "")
    public JsonResponse getCoreCiData() {

        return okayWithData(null);
    }*/


}
