package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.domain.AttachFile;
import com.webank.taskman.support.s3.dto.DownloadAttachFileResponse;
import org.springframework.web.multipart.MultipartFile;


public interface AttachFileService extends IService<AttachFile> {


    DownloadAttachFileResponse downloadServiceRequestAttachFile(String serviceRequestId) throws TaskmanRuntimeException;

    String uploadServiceRequestAttachFile(MultipartFile attachFile) throws TaskmanRuntimeException;
}
