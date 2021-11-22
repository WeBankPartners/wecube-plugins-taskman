package db

import (
	"context"
	"fmt"
	"github.com/WeBankPartners/go-common-lib/file_server"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"time"
)

func UploadAttachFile(requestId, taskId, fileName, operator string, fileContent []byte) error {
	if requestId == "" && taskId == "" {
		return fmt.Errorf("requestId or taskId can not empty ")
	}
	ms, err := getMinioServerObj()
	if err != nil {
		return err
	}
	fileGuid := guid.CreateGuid()
	uploadParam := file_server.MinioParam{Ctx: context.Background(), Bucket: models.Config.AttachFile.Bucket}
	uploadParam.FileContent = fileContent
	uploadParam.ObjectName = fmt.Sprintf("%s_%s", fileName, fileGuid)
	err = ms.Upload(uploadParam)
	if err != nil {
		return err
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	var actions []*execAction
	if requestId != "" {
		actions = append(actions, &execAction{Sql: "insert into attach_file(id,name,s3_bucket_name,s3_key_name,request,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?,?,?)", Param: []interface{}{fileGuid, fileName, models.Config.AttachFile.Bucket, uploadParam.ObjectName, requestId, operator, nowTime, operator, nowTime}})
	} else if taskId != "" {
		actions = append(actions, &execAction{Sql: "insert into attach_file(id,name,s3_bucket_name,s3_key_name,task,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?,?,?)", Param: []interface{}{fileGuid, fileName, models.Config.AttachFile.Bucket, uploadParam.ObjectName, taskId, operator, nowTime, operator, nowTime}})
	}
	return transaction(actions)
}

func DownloadAttachFile(fileId string) (fileContent []byte, fileName string, err error) {
	fileObj, tmpErr := getAttachFileInfo(fileId)
	if tmpErr != nil {
		err = tmpErr
		return
	}
	fileName = fileObj.Name
	ms, minioErr := getMinioServerObj()
	if minioErr != nil {
		err = minioErr
		return
	}
	downloadParam := file_server.MinioParam{Ctx: context.Background(), Bucket: models.Config.AttachFile.Bucket}
	downloadParam.ObjectName = fileObj.S3KeyName
	fileContent, err = ms.Download(downloadParam)
	return
}

func RemoveAttachFile(fileId string) error {
	fileObj, err := getAttachFileInfo(fileId)
	if err != nil {
		return err
	}
	ms, minioErr := getMinioServerObj()
	if minioErr != nil {
		return minioErr
	}
	removeParam := file_server.MinioParam{Ctx: context.Background(), Bucket: models.Config.AttachFile.Bucket}
	removeParam.ObjectName = fileObj.S3KeyName
	err = ms.Remove(removeParam)
	if err != nil {
		return err
	}
	_, err = x.Exec("delete from attach_file where id=?", fileId)
	return err
}

func getAttachFileInfo(fileId string) (fileObj models.AttachFileTable, err error) {
	var attachFileTable []*models.AttachFileTable
	err = x.SQL("select * from attach_file where id=?", fileId).Find(&attachFileTable)
	if err != nil {
		return fileObj, err
	}
	if len(attachFileTable) == 0 {
		return fileObj, fmt.Errorf("Can not find attach file with id:%s ", fileId)
	}
	fileObj = *attachFileTable[0]
	return fileObj, nil
}

func GetRequestAttachFileList(requestId string) []*models.AttachFileTable {
	attachFileTable := []*models.AttachFileTable{}
	x.SQL("select id,name from attach_file where request=?", requestId).Find(&attachFileTable)
	return attachFileTable
}

func GetTaskAttachFileList(taskId string) []*models.AttachFileTable {
	attachFileTable := []*models.AttachFileTable{}
	x.SQL("select id,name from attach_file where task=?", taskId).Find(&attachFileTable)
	return attachFileTable
}

func getMinioServerObj() (ms file_server.MinioServer, err error) {
	ms = file_server.MinioServer{ServerAddress: models.Config.AttachFile.MinioAddress, AccessKey: models.Config.AttachFile.MinioAccessKey, SecretKey: models.Config.AttachFile.MinioSecretKey, SSL: models.Config.AttachFile.SSL}
	err = ms.Init()
	return
}
