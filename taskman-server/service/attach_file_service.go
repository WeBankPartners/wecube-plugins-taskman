package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/go-common-lib/file_server"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/exterror"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"time"
)

type AttachFileService struct {
	attachFileDao *dao.AttachFileDao
}

func UploadAttachFile(requestId, taskId, taskHandleId, fileName, operator string, fileContent []byte) error {
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
		uploadParam.ObjectName = fmt.Sprintf("%s%s_%s", requestTemplateId, fileName, taskHandleId)
	}
	err = ms.Upload(uploadParam)
	if err != nil {
		return err
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	var actions []*dao.ExecAction
	if requestId != "" {
		actions = append(actions, &dao.ExecAction{Sql: "insert into attach_file(id,name,s3_bucket_name,s3_key_name,request,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?,?,?)", Param: []interface{}{fileGuid, fileName, models.Config.AttachFile.Bucket, uploadParam.ObjectName, requestId, operator, nowTime, operator, nowTime}})
	} else if taskId != "" {
		actions = append(actions, &dao.ExecAction{Sql: "insert into attach_file(id,name,s3_bucket_name,s3_key_name,task,created_by,created_time,updated_by,updated_time,task_handle) value (?,?,?,?,?,?,?,?,?,?)", Param: []interface{}{fileGuid, fileName, models.Config.AttachFile.Bucket, uploadParam.ObjectName, taskId, operator, nowTime, operator, nowTime, taskHandleId}})
	}
	return dao.Transaction(actions)
}

func DownloadAttachFile(fileObj models.AttachFileTable) (fileContent []byte, fileName string, err error) {
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
	fileObj, err = GetAttachFileInfo(fileId)
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
	_, err = dao.X.Exec("delete from attach_file where id=?", fileId)
	return
}

func GetAttachFileInfo(fileId string) (fileObj models.AttachFileTable, err error) {
	var attachFileTable []*models.AttachFileTable
	err = dao.X.SQL("select * from attach_file where id=?", fileId).Find(&attachFileTable)
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
	dao.X.SQL("select id,name from attach_file where request=?", requestId).Find(&attachFileTable)
	return attachFileTable
}

func GetTaskAttachFileList(taskId string) []*models.AttachFileTable {
	attachFileTable := []*models.AttachFileTable{}
	dao.X.SQL("select id,name from attach_file where task=?", taskId).Find(&attachFileTable)
	return attachFileTable
}

func GetAttachFileListByTaskHandleId(taskHandleId string) []*models.AttachFileTable {
	var attachFileTable []*models.AttachFileTable
	dao.X.SQL("select id,name from attach_file where task_handle=?", taskHandleId).Find(&attachFileTable)
	return attachFileTable
}

func getMinioServerObj() (ms file_server.MinioServer, err error) {
	ms = file_server.MinioServer{ServerAddress: models.Config.AttachFile.MinioAddress, AccessKey: models.Config.AttachFile.MinioAccessKey, SecretKey: models.Config.AttachFile.MinioSecretKey, SSL: models.Config.AttachFile.SSL}
	err = ms.Init()
	return
}

func getAttachFileRequestTemplate(requestId, taskId string) string {
	var requestTable []*models.RequestTable
	if requestId != "" {
		dao.X.SQL("select request_template from request where id=?", requestId).Find(&requestTable)
	} else {
		dao.X.SQL("select request_template from request where id in (select request from task where id=?)", taskId).Find(&requestTable)
	}
	if len(requestTable) > 0 {
		return requestTable[0].RequestTemplate
	}
	return ""
}

func GetTaskHandleAttachFileList(taskHandleIds []string) (result []*models.AttachFileTable, err error) {
	result = make([]*models.AttachFileTable, 0)
	var attachFiles []*models.AttachFileTable
	taskHandleIdsFilterSql, taskHandleIdsFilterParams := dao.CreateListParams(taskHandleIds, "")
	err = dao.X.SQL("select * from attach_file where task_handle in ("+taskHandleIdsFilterSql+")", taskHandleIdsFilterParams...).Find(&attachFiles)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	if len(attachFiles) > 0 {
		result = attachFiles
	}
	return
}

// CheckDownloadPermission 检查下载权限
func CheckDownloadPermission(attachFile models.AttachFileTable, roles []string) (checkPermission bool, err error) {
	var task models.TaskTable
	var request models.RequestTable
	var taskTemplateDtoList []*models.TaskTemplateDto
	var roleMap = make(map[string]bool)
	var requestTemplateRoleList []*models.RequestTemplateRoleTable
	var taskHandleList []*models.TaskHandleTable
	var requestId string
	// 权限校验,下载需要文件模版使用角色和任务模版角色权限
	if attachFile.Request == "" {
		if task, err = GetSimpleTask(attachFile.Task); err != nil {
			return
		}
		requestId = task.Request
	} else {
		requestId = attachFile.Request
	}
	if requestId == "" {
		// 可能编排创建的任务,没有请求的任务
		if taskHandleList, err = GetTaskHandleService().GetTaskHandleListByTaskId(attachFile.Task); err != nil {
			return
		}
		if len(taskHandleList) > 0 {
			for _, taskHandle := range taskHandleList {
				if taskHandle.Role != "" {
					roleMap[taskHandle.Role] = true
				}
			}
		}
	} else {
		if request, err = GetSimpleRequest(requestId); err != nil {
			return
		}
		// 统计审批、任务配置的处理角色
		if request.TaskApprovalCache != "" {
			json.Unmarshal([]byte(request.TaskApprovalCache), &taskTemplateDtoList)
			if len(taskTemplateDtoList) > 0 {
				for _, dto := range taskTemplateDtoList {
					if dto != nil && len(dto.HandleTemplates) > 0 {
						for _, template := range dto.HandleTemplates {
							if template.Role != "" {
								roleMap[template.Role] = true
							}
						}
					}
				}
			}
		}
		// 统计模版使用角色
		if requestTemplateRoleList, err = GetRequestTemplateService().GetRequestTemplateRole(request.RequestTemplate); err != nil {
			return
		}
		if len(requestTemplateRoleList) > 0 {
			for _, requestTemplateRole := range requestTemplateRoleList {
				if requestTemplateRole.RoleType == string(models.RolePermissionUse) {
					roleMap[requestTemplateRole.Role] = true
				}
			}
		}
	}
	for _, role := range roles {
		if _, ok := roleMap[role]; ok {
			checkPermission = true
			break
		}
	}
	return
}

func CheckAttachFilePermission(fileId, operator string, roles []string) error {
	fileObj, err := GetAttachFileInfo(fileId)
	if err != nil {
		return err
	}
	var legalRoles []string
	if fileObj.Request != "" {
		var requestTemplateRoles []*models.RequestTemplateRoleTable
		var requestTable []*models.RequestTable
		dao.X.SQL("select id,created_by from request where id=?", fileObj.Request).Find(&requestTable)
		if len(requestTable) > 0 {
			if requestTable[0].CreatedBy == operator {
				return nil
			}
		}
		dao.X.SQL("select * from request_template_role where request_template in (select request_template from request where id=?)", fileObj.Request).Find(&requestTemplateRoles)
		for _, v := range requestTemplateRoles {
			if v.RoleType == "USE" {
				legalRoles = append(legalRoles, v.Role)
			}
		}
	}
	if fileObj.Task != "" {
		var taskHandleList []*models.TaskHandleTable
		dao.X.SQL("select * from task_handle where task=?", fileObj.Task).Find(&taskHandleList)
		for _, v := range taskHandleList {
			legalRoles = append(legalRoles, v.Role)
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
