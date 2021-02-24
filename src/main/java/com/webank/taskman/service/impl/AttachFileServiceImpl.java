package com.webank.taskman.service.impl;


import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.commons.AppProperties.ServiceTaskmanProperties;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.domain.AttachFile;
import com.webank.taskman.mapper.AttachFileMapper;
import com.webank.taskman.service.AttachFileService;
import com.webank.taskman.support.s3.S3Client;
import com.webank.taskman.support.s3.dto.DownloadAttachFileResponse;
import com.webank.taskman.utils.SystemUtils;
import org.apache.commons.io.FileUtils;
import org.apache.commons.io.FilenameUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import java.io.File;


@Service
public class AttachFileServiceImpl extends ServiceImpl<AttachFileMapper, AttachFile> implements AttachFileService {


    @Autowired
    private ServiceTaskmanProperties serviceTaskmanProperties;

    public DownloadAttachFileResponse downloadServiceRequestAttachFile(String attachFileId){

        AttachFile attachFile = getById(attachFileId);
        if (null == attachFile){
            throw new TaskmanRuntimeException("3011", "This service request has no attach file");
        }
        try {
            String fileName = attachFile.getAttachFileName();
            String tempDownloadFilePath = SystemUtils.getTempFolderPath() + fileName;
            File downloadFile = new File(tempDownloadFilePath);

            new S3Client(serviceTaskmanProperties).downFile(fileName, tempDownloadFilePath);
            DownloadAttachFileResponse response = new DownloadAttachFileResponse(
                    FileUtils.readFileToByteArray(downloadFile), fileName);

            FileUtils.forceDelete(downloadFile);
            return response;
        }catch (Exception e){
            throw new TaskmanRuntimeException(e.getMessage());
        }
    }

    @Override
    public String uploadServiceRequestAttachFile(MultipartFile attachFile) {
        if (attachFile.isEmpty()) {
            throw new TaskmanRuntimeException("3008", "Empty file!");
        }
        try {
            String fileExtension = FilenameUtils.getExtension(attachFile.getOriginalFilename());
//        if (!fileExtension.equals("xlsx") && !fileExtension.equals("xls")) {
//            throw new ServiceMgmtException("3009", "Only support Excel file");
//        }

            String tmpFileName = String.valueOf(System.currentTimeMillis());
            File tempUploadFile = new File(SystemUtils.getTempFolderPath() + tmpFileName);
            attachFile.transferTo(tempUploadFile);

            String uploadFileName = FilenameUtils.getBaseName(attachFile.getOriginalFilename()) + "-" + tmpFileName + "."
                    + fileExtension;

            String s3Url = new S3Client(serviceTaskmanProperties).uploadFile(uploadFileName, tempUploadFile);
            AttachFile attachFileObject = new AttachFile(uploadFileName, s3Url,
                    serviceTaskmanProperties.getS3DefaultBucket(), uploadFileName);
            save(attachFileObject);
            FileUtils.forceDelete(tempUploadFile);
            return attachFileObject.getId();
        }catch (Exception e){
            throw new TaskmanRuntimeException(e.getMessage());
        }
    }

}
