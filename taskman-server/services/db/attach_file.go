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
	requestTemplateId := getAttachFileRequestTemplate(requestId, taskId)
	if requestTemplateId != "" {
		requestTemplateId = requestTemplateId + "/"
	}
	if requestId != "" {
		uploadParam.ObjectName = fmt.Sprintf("%s%s_%s", requestTemplateId, fileName, requestId)
	} else {
		uploadParam.ObjectName = fmt.Sprintf("%s%s_%s", requestTemplateId, fileName, taskId)
	}
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

func RemoveAttachFile(fileId string) (fileObj models.AttachFileTable, err error) {
	fileObj, err = getAttachFileInfo(fileId)
	if err != nil {
		return
	}
	ms, minioErr := getMinioServerObj()
	if minioErr != nil {
		err = minioErr
		return
	}
	removeParam := file_server.MinioParam{Ctx: context.Background(), Bucket: models.Config.AttachFile.Bucket}
	removeParam.ObjectName = fileObj.S3KeyName
	err = ms.Remove(removeParam)
	if err != nil {
		return
	}
	_, err = x.Exec("delete from attach_file where id=?", fileId)
	return
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

func CheckAttachFilePermission(fileId, operator, operation string, roles []string) error {
	fileObj, err := getAttachFileInfo(fileId)
	if err != nil {
		return err
	}
	var legalRoles []string
	if fileObj.Request != "" {
		var requestTemplateRoles []*models.RequestTemplateRoleTable
		if operation == "download" {
			x.SQL("select distinct t1.`role` from (select `role` from task_template_role where task_template in (select id from task_template where request_template in (select request_template from request where id=?)) union select `role` from request_template_role where request_template in (select request_template from request where id=?)) t1", fileObj.Request, fileObj.Request).Find(&requestTemplateRoles)
			for _, v := range requestTemplateRoles {
				legalRoles = append(legalRoles, v.Role)
			}
		} else {
			var requestTable []*models.RequestTable
			x.SQL("select id,created_by from request where id=?", fileObj.Request).Find(&requestTable)
			if len(requestTable) > 0 {
				if requestTable[0].CreatedBy == operator {
					return nil
				}
			}
			x.SQL("select * from request_template_role where request_template in (select request_template from request where id=?)", fileObj.Request).Find(&requestTemplateRoles)
			for _, v := range requestTemplateRoles {
				if v.RoleType == "USE" {
					legalRoles = append(legalRoles, v.Role)
				}
			}
		}
	}
	if fileObj.Task != "" {
		var taskTemplateRoles []*models.TaskTemplateRoleTable
		if operation == "download" {
			x.SQL("select distinct t1.`role` from (select `role` from task_template_role where task_template in (select id from task_template where request_template in (select request_template from request where id in (select request from task where id=?))) union select `role` from request_template_role where request_template in (select request_template from request where id in (select request from task where id=?))) t1", fileObj.Task, fileObj.Task).Find(&taskTemplateRoles)
			for _, v := range taskTemplateRoles {
				legalRoles = append(legalRoles, v.Role)
			}
		} else {
			x.SQL("select * from task_template_role where task_template in (select task_template from task where id=?)", fileObj.Task).Find(&taskTemplateRoles)
			for _, v := range taskTemplateRoles {
				legalRoles = append(legalRoles, v.Role)
			}
		}
	}
	legalFlag := false
	for _, v := range legalRoles {
		for _, vv := range roles {
			if v == vv {
				legalFlag = true
				break
			}
		}
		if legalFlag {
			break
		}
	}
	if !legalFlag {
		return fmt.Errorf("Permission illegal ")
	}
	return nil
}

func getAttachFileRequestTemplate(requestId, taskId string) string {
	var requestTable []*models.RequestTable
	if requestId != "" {
		x.SQL("select request_template from request where id=?", requestId).Find(&requestTable)
	} else {
		x.SQL("select request_template from request where id in (select request from task where id=?)", taskId).Find(&requestTable)
	}
	if len(requestTable) > 0 {
		return requestTable[0].RequestTemplate
	}
	return ""
}
