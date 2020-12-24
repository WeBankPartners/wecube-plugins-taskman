package com.webank.taskman.controller.x100;

import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.base.JsonResponse;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.service.AttachFileService;
import com.webank.taskman.service.RequestInfoService;
import com.webank.taskman.service.RequestTemplateService;
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
import javax.servlet.http.HttpServletResponse;
import java.util.List;
import java.util.Map;

import static com.webank.taskman.base.JsonResponse.okayWithData;


@Api(tags = {"1、 CoreResource inteface API"})
@RestController
@RequestMapping("/v1/core-resources")
public class CoreResourceController {


    private static final Logger log = LoggerFactory.getLogger(CoreResourceController.class);

    @Autowired
    CoreServiceStub coreServiceStub;

    @ApiOperationSupport(order = 1)
    @GetMapping("/roles")
    @ApiOperation(value = "auth-role-all", notes = "")
    public JsonResponse<List<RolesDataResponse>> authRoleAll()
    {
        return okayWithData(coreServiceStub.authRoleAll());
    }

    @ApiOperationSupport(order = 2)
    @GetMapping("/users/current-user/roles")
    @ApiOperation(value = "auth-role-current-user", notes = "")
    public JsonResponse<List<RolesDataResponse>> authRoleCurrentUser()
    {
        String currentUserName = AuthenticationContextHolder.getCurrentUsername();
        return okayWithData(coreServiceStub.authRoleCurrentUser(currentUserName));
    }

    @ApiOperationSupport(order = 3)
    @GetMapping("/platform/process-definition-keys")
    @ApiOperation(value = "platform-process-all", notes = "")
    public JsonResponse<List<WorkflowDefInfoDto>> platformProcessAll()
    {
        return okayWithData(coreServiceStub.platformProcessAll());
    }

    @ApiOperationSupport(order = 4)
    @GetMapping("/platform/process-definitions-nodes/{proc-def-id}")
    @ApiOperation(value = "platform-process-nodes", notes = "")
    public JsonResponse<List<WorkflowNodeDefInfoDto>> platformProcessNodes(@PathVariable("proc-def-id") String procDefId)
    {
        return okayWithData(coreServiceStub.platformProcessNodes(procDefId));
    }


    @ApiOperationSupport(order = 5)
    @GetMapping(value = {"/platform/models","/platform/models/{package-name}"})
    @ApiOperation(value = "platform-process-models", notes = "")
    public JsonResponse<List<PluginPackageDataModelDto>> platformProcessModels(@PathVariable(value = "package-name",required = false) String packageName)
    {
        return okayWithData(coreServiceStub.platformProcessModels(packageName));
    }

    /**/
    @ApiOperationSupport(order = 6)
    @GetMapping("/platform/{proc-def-id}/root-entity")
    @ApiOperation(value = "platform-process-root-entity", notes = "")
    public JsonResponse<List<Map<String,Object>>> platformProcessRootEntity(@PathVariable("proc-def-id") String procDefId)
    {
        return okayWithData(coreServiceStub.platformProcessRootEntity(procDefId));
    }

    @ApiOperationSupport(order = 7)
    @GetMapping("/platform/models/package/{package-name}/entity/{entity-name}/attributes")
    @ApiOperation(value = "platform-process-entity-attributes", notes = "")
    public JsonResponse platformProcessEntityAttributes(
            @PathVariable("package-name") String packageName,@PathVariable("entity-name")String entity)
    {
        return okayWithData(coreServiceStub.platformProcessEntityAttributes(packageName,entity));
    }

    @ApiOperationSupport(order = 8)
    @GetMapping("/platform/packages/{package-name}/entities/{entity-name}/retrieve")
    @ApiOperation(value = "platform-process-entity-retrieve", notes = "")
    public JsonResponse<List<Object>> platformProcessEntityRetrieve(@PathVariable("package-name") String packageName,
        @PathVariable("entity-name")String entity,
        @ApiParam(value = "filters",required = false,type = "query") @RequestParam(required = false)String filters)
    {
        return okayWithData(coreServiceStub.platformProcessEntityRetrieve(packageName,entity,filters));
    }

    @Autowired
    RequestTemplateService requestTemplateService;

    @ApiOperationSupport(order = 10)
    @GetMapping("/platform/process/definitions/{proc-def-id}/preview/entities/{entity-data-id}")
    @ApiOperation(value = "platform-process-data-preview", notes = "")
    public JsonResponse<ProcessDataPreviewDto> platformProcessDataPreview( @PathVariable("proc-def-id") String procDefId,
                                        @PathVariable("entity-data-id")String entityDataId)
    {
        ProcessDataPreviewDto processDataPreviewDto = coreServiceStub.platformProcessDataPreview(procDefId,entityDataId);
        log.info("platform-process-data-preview is result:{}",processDataPreviewDto);
        return okayWithData(processDataPreviewDto);
    }


    @ApiOperationSupport(order = 11)
    @GetMapping("/platform/process/tasknodes/session/{process-session-id}/tasknode-bindings")
    @ApiOperation(value = "platform-process-tasknode-bindings", notes = "")
    public JsonResponse<List<TaskNodeDefObjectBindInfoDto>> platformProcessTasknodeBindings(@PathVariable(name = "process-session-id") String processSessionId)
    {
        return okayWithData(coreServiceStub.platformProcessTasknodeBindings(processSessionId));

    }

    @Autowired
    RequestInfoService requestInfoService;


    @ApiOperationSupport(order = 12)
    @GetMapping("/platform/crate/{proc-def-id}/{entity-data-id}")
    @ApiOperation(value = "platform-process-create-examples", notes = "just is examples")
    public JsonResponse<DynamicWorkflowInstInfoDto> platformProcessCreateExamples(
            @PathVariable("proc-def-id") String procDefId,@PathVariable("entity-data-id")String entityDataId)
    {
        return okayWithData(requestInfoService.createDynamicWorkflowInstCreationInfoDto(procDefId,entityDataId));
    }
    @ApiOperationSupport(order = 12)
    @PostMapping("/platform/crate")
    @ApiOperation(value = "platform-process-create", notes = "")
    public JsonResponse<DynamicWorkflowInstInfoDto> platformProcessCreate(
            @RequestBody DynamicWorkflowInstCreationInfoDto creationInfoDto)
    {
        return okayWithData(coreServiceStub.createNewWorkflowInstance(creationInfoDto));
    }

    @Autowired
    AttachFileService attachFileService;

    @PostMapping("/attach-file")
    @ApiOperation(value = "S3-upload-file", notes = "")
    public JsonResponse S3UploadFile(@RequestParam(value = "file") MultipartFile attachFile)throws Exception
    {
        String attachFileId  = attachFileService.uploadServiceRequestAttachFile(attachFile);

        return okayWithData(attachFileId);
    }


    @GetMapping("/{attach-id}/attach-file")
    @ApiOperation(value = "S3-download-file", notes = "")
    public void S3DownloadFile(@PathVariable(value = "attach-id") String serviceRequestId,HttpServletResponse response) throws Exception
    {
        if (serviceRequestId == null || serviceRequestId.isEmpty())
            throw new Exception("Invalid service-request-id: " + serviceRequestId);
        try {
            ServletOutputStream out = response.getOutputStream();
            DownloadAttachFileResponse attachFileInfo = attachFileService.downloadServiceRequestAttachFile(serviceRequestId);

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
