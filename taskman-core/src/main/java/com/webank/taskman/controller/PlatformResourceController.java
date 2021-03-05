package com.webank.taskman.controller;

import static com.webank.taskman.base.JsonResponse.okayWithData;

import javax.servlet.ServletOutputStream;
import javax.servlet.http.HttpServletResponse;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;

import com.webank.taskman.base.JsonResponse;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.service.AttachFileService;
import com.webank.taskman.support.core.PlatformCoreServiceRestClient;
import com.webank.taskman.support.core.dto.DataModelEntityDto;
import com.webank.taskman.support.core.dto.DynamicWorkflowInstCreationInfoDto;
import com.webank.taskman.support.core.dto.ProcessDataPreviewDto;
import com.webank.taskman.support.s3.dto.DownloadAttachFileResponse;

@RestController
@RequestMapping("/v1/core-resources")
public class PlatformResourceController {

    private static final Logger log = LoggerFactory.getLogger(PlatformResourceController.class);

    @Autowired
    private PlatformCoreServiceRestClient platformCoreServiceRestClient;

    @Autowired
    private AttachFileService attachFileService;

    /**
     * 
     * @return
     */
    @GetMapping("/roles")
    public JsonResponse getAllPlatformAuthRoles() {
        return okayWithData(platformCoreServiceRestClient.getAllPlatformAuthRoles());
    }

    /**
     * 
     * @return
     */
    @GetMapping("/users/current-user/roles")
    public JsonResponse getAllAuthRolesOfCurrentUser() {
        return okayWithData(platformCoreServiceRestClient.getAllAuthRolesOfCurrentUser());
    }

    /**
     * 
     * @return
     */
    @GetMapping("/platform/process-definitions")
    public JsonResponse getAllLatestPlatformProcesses() {
        return okayWithData(platformCoreServiceRestClient.getAllLatestPlatformProcesses());
    }

    /**
     * 
     * @param procDefId
     * @return
     */
    @GetMapping("/platform/process-definitions/{proc-def-id}/nodes")
    public JsonResponse getPlatformProcessDefinitionNodes(@PathVariable("proc-def-id") String procDefId) {
        return okayWithData(platformCoreServiceRestClient.getPlatformProcessDefinitionNodes(procDefId));
    }

    /**
     * 
     * @return
     */
    @GetMapping("/platform/models")
    public JsonResponse getAllPlatformProcessModels() {
        return okayWithData(platformCoreServiceRestClient.getAllPlatformProcessModels(null));
    }

    /**
     * 
     * @param packageName
     * @return
     */
    @GetMapping("/platform/models/{package-name}")
    public JsonResponse getAllPlatformProcessModelsOfPackage(@PathVariable("package-name") String packageName) {
        return okayWithData(platformCoreServiceRestClient.getAllPlatformProcessModels(packageName));
    }

    /**
     * 
     * @param procDefId
     * @return
     */
    @GetMapping("/platform/process/{proc-def-id}/root-entities")
    public JsonResponse getPlatformProcessRootEntities(@PathVariable("proc-def-id") String procDefId) {
        return okayWithData(platformCoreServiceRestClient.getPlatformProcessRootEntities(procDefId));
    }

    /**
     * 
     * @param packageName
     * @param entityName
     * @return
     */
    @GetMapping("/platform/models/package/{package-name}/entity/{entity-name}")
    public JsonResponse platformProcessEntityInfo(@PathVariable(value = "package-name") String packageName,
            @PathVariable(value = "entity-name") String entityName) {
        DataModelEntityDto result = platformCoreServiceRestClient.getEntityByPackageNameAndName(packageName,
                entityName);
        return okayWithData(result);
    }

    /**
     * 
     * @param packageName
     * @param entity
     * @return
     */
    @GetMapping("/platform/models/package/{package-name}/entity/{entity-name}/attributes")
    public JsonResponse platformProcessEntityAttributes(@PathVariable("package-name") String packageName,
            @PathVariable("entity-name") String entity) {
        return okayWithData(platformCoreServiceRestClient.platformProcessEntityAttributes(packageName, entity));
    }

    /**
     * 
     * @param packageName
     * @param entity
     * @param filters
     * @return
     */
    @GetMapping("/platform/packages/{package-name}/entities/{entity-name}/retrieve")
    public JsonResponse platformProcessEntityRetrieve(@PathVariable("package-name") String packageName,
            @PathVariable("entity-name") String entity, @RequestParam(required = false) String filters) {
        return okayWithData(platformCoreServiceRestClient.platformProcessEntityRetrieve(packageName, entity, filters));
    }

    /**
     * 
     * @param procDefId
     * @param entityDataId
     * @return
     */
    @GetMapping("/platform/process/definitions/{proc-def-id}/preview/entities/{entity-data-id}")
    public JsonResponse platformProcessDataPreview(@PathVariable("proc-def-id") String procDefId,
            @PathVariable("entity-data-id") String entityDataId) {
        ProcessDataPreviewDto processDataPreviewDto = platformCoreServiceRestClient
                .platformProcessDataPreview(procDefId, entityDataId);
        log.info("platform-process-data-preview is result:{}", processDataPreviewDto);
        return okayWithData(processDataPreviewDto);
    }

    /**
     * 
     * @param processSessionId
     * @return
     */
    @GetMapping("/platform/process/tasknodes/session/{process-session-id}/tasknode-bindings")
    public JsonResponse platformProcessTasknodeBindings(
            @PathVariable(name = "process-session-id") String processSessionId) {
        return okayWithData(platformCoreServiceRestClient.platformProcessTasknodeBindings(processSessionId));

    }

    /**
     * 
     * @param creationInfoDto
     * @return
     */
    @PostMapping("/platform/create")
    public JsonResponse platformProcessCreate(@RequestBody DynamicWorkflowInstCreationInfoDto creationInfoDto) {

        return okayWithData(platformCoreServiceRestClient.createNewWorkflowInstance(creationInfoDto));
    }

    /**
     * 
     * @param attachFile
     * @return
     */
    @PostMapping("/attach-file")
    public JsonResponse uploadS3File(@RequestParam(value = "file") MultipartFile attachFile){
        String attachFileId = attachFileService.uploadServiceRequestAttachFile(attachFile);

        return okayWithData(attachFileId);
    }

    /**
     * 
     * @param serviceRequestId
     * @param response
     */
    @GetMapping("/{attach-id}/attach-file")
    public void downloadS3File(@PathVariable(value = "attach-id") String serviceRequestId,
            HttpServletResponse response) {
        if (serviceRequestId == null || serviceRequestId.isEmpty()) {
            throw new TaskmanRuntimeException("Invalid service-request-id: " + serviceRequestId);
        }
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
