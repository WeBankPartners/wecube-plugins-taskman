package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.AttachFile;
import com.webank.taskman.dto.DownloadAttachFileResponse;
import org.springframework.web.multipart.MultipartFile;


public interface AttachFileService extends IService<AttachFile> {


    DownloadAttachFileResponse downloadServiceRequestAttachFile(String serviceRequestId) throws Exception;

    String uploadServiceRequestAttachFile(MultipartFile attachFile) throws Exception;
}
