package com.webank.taskman.controller.x100;

import com.fasterxml.jackson.core.JsonParseException;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.base.JsonResponse;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.service.AttachFileService;
import com.webank.taskman.support.core.CoreServiceStub;
import com.webank.taskman.support.core.dto.*;
import com.webank.taskman.support.s3.dto.DownloadAttachFileResponse;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.ApiParam;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import javax.servlet.ServletOutputStream;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.util.List;

import static com.webank.taskman.base.JsonResponse.okayWithData;


@Api(tags = {"1、 CoreResource inteface API"})
@RestController
@RequestMapping("/v1/core-resources")
public class CoreResourceController {


    private static final Logger log = LoggerFactory.getLogger(CoreResourceController.class);

    @Autowired
    CoreServiceStub coreServiceStub;

    @ApiOperationSupport(order = 1)
    @GetMapping("/users/current-user/roles")
    @ApiOperation(value = "auth-role-current-user", notes = "")
    public JsonResponse<List<RolesDataResponse>> getRolesByCurrentUser(HttpServletRequest httpRequest) {
        String currentUserName = AuthenticationContextHolder.getCurrentUsername();
        return okayWithData(coreServiceStub.getRolesByUserName(currentUserName));
    }

    @ApiOperationSupport(order = 2)
    @GetMapping("/roles")
    @ApiOperation(value = "auth-role-all", notes = "")
    public JsonResponse<List<RolesDataResponse>> getAllRoles() throws JsonParseException, JsonMappingException, IOException {
        return okayWithData(coreServiceStub.getAllRoles());
    }

    @ApiOperationSupport(order = 3)
    @GetMapping("/platform/process-definition-keys")
    @ApiOperation(value = "platform-process-all", notes = "")
    public JsonResponse<List<WorkflowDefInfoDto>> fetchLatestReleasedWorkflowDefs() {
        return okayWithData(coreServiceStub.fetchLatestReleasedWorkflowDefs());
    }

    @ApiOperationSupport(order = 4)
    @GetMapping("/platform/process-definitions-nodes/{proc-def-id}")
    @ApiOperation(value = "platform-process-nodes", notes = "")
    public JsonResponse<List<WorkflowNodeDefInfoDto>> getTaskNodes(@PathVariable("proc-def-id") String procDefId) {
        return okayWithData(coreServiceStub.fetchWorkflowTasknodeInfos(procDefId));
    }

    @ApiOperationSupport(order = 5)
    @PostMapping("/platform/crate")
    @ApiOperation(value = "platform-process-crate", notes = "")
    public JsonResponse<CoreResponse.DynamicWorkflowInstInfoDto> createNewWorkflowInstance(
            @RequestBody DynamicWorkflowInstCreationInfoDto creationInfoDto)
    {
        return okayWithData(coreServiceStub.createNewWorkflowInstance(creationInfoDto));
    }



    @ApiOperationSupport(order = 6)
    @GetMapping("/platform/models")
    @ApiOperation(value = "platform-process-models", notes = "")
    public JsonResponse allDataModels() {
        return okayWithData(coreServiceStub.allDataModels());
    }

    @ApiOperationSupport(order = 7)
    @GetMapping("/platform/models/{package-name}")
    @ApiOperation(value = "platform-process-models-package", notes = "")
    public JsonResponse allDataModels(@PathVariable("package-name") String packageName) {
        return okayWithData(coreServiceStub.getModelsByPackage(packageName));
    }

    /**/
    @ApiOperationSupport(order = 8)
    @GetMapping("/platform/{proc-def-id}/root-entity")
    @ApiOperation(value = "platform-process-root-entity", notes = "")
    public JsonResponse getProcessDefinitionRootEntitiesByProcDefKey(@PathVariable("proc-def-id") String procDefId) {

        return okayWithData(coreServiceStub.getProcessDefinitionRootEntitiesByProcDefKey(procDefId));
    }

    @ApiOperationSupport(order = 9)
    @GetMapping("/platform/models/package/{package-name}/entity/{entity-name}/attributes")
    @ApiOperation(value = "platform-process-entity-attributes", notes = "")
    public JsonResponse getAttributesByPackageEntity(
            @PathVariable("package-name") String packageName,@PathVariable("entity-name")String entity) {
        return okayWithData(coreServiceStub.getAttributesByPackageEntity(packageName,entity));
    }
    @ApiOperationSupport(order = 10)
    @GetMapping("/platform/packages/{package-name}/entities/{entity-name}/retrieve")
    @ApiOperation(value = "platform-process-entity-retrieve", notes = "")
    public JsonResponse retrieveEntity( @PathVariable("package-name") String packageName,
                @PathVariable("entity-name")String entity,
           @ApiParam(value = "filters",required = false,type = "query") @RequestParam(required = false)String filters)
    {
        return okayWithData(coreServiceStub.retrieveEntity(packageName,entity,filters));
    }
    @ApiOperationSupport(order = 11)
    @GetMapping("/platform/process/definitions/{proc-def-id}/preview/entities/{entity-data-id}")
    @ApiOperation(value = "platform-process-data-preview", notes = "")
    public JsonResponse<ProcessDataPreviewDto> getProcessDataPreview( @PathVariable("proc-def-id") String procDefId,
                                        @PathVariable("entity-data-id")String entityDataId)
    {
        return okayWithData(coreServiceStub.getProcessDataPreview(procDefId,entityDataId));
    }


    @Autowired
    AttachFileService attachFileService;

    @PostMapping("/attach-file")
    @ApiOperation(value = "S3-upload-file", notes = "")
    public JsonResponse uploadServiceRequestAttachFile(@RequestParam(value = "file") MultipartFile attachFile)
            throws Exception {
        String attachFileId  = attachFileService.uploadServiceRequestAttachFile(attachFile);

        return okayWithData(attachFileId);
    }


    @GetMapping("/{attach-id}/attach-file")
    @ApiOperation(value = "S3-download-file", notes = "")
    public void downloadServiceRequestAttachFile(@PathVariable(value = "attach-id") String serviceRequestId,
                                                 HttpServletResponse response) throws Exception {
        if (serviceRequestId == null || serviceRequestId.isEmpty())
            throw new Exception("Invalid service-request-id: " + serviceRequestId);
        try {
            ServletOutputStream out = response.getOutputStream();
            DownloadAttachFileResponse attachFileInfo = attachFileService
                    .downloadServiceRequestAttachFile(serviceRequestId);

            response.setCharacterEncoding("utf-8");
            response.setContentType("application/vnd.ms-excel;charset=UTF-8");
            response.setHeader(HttpHeaders.CONTENT_DISPOSITION,
                    "attachment;fileName=" + attachFileInfo.getAttachFileName());
            response.setHeader("Accept", MediaType.APPLICATION_OCTET_STREAM_VALUE);
            out.write(attachFileInfo.getFileByteArray());
            out.flush();
            out.close();
        } catch (Exception e) {
            String errorMessage = String.format("Failed to download attach file(service request Id:%d) due to %s ",
                    serviceRequestId, e.getMessage());
            throw new TaskmanRuntimeException("3000", errorMessage, serviceRequestId, e.getMessage());
        }
    }

}
