package service

import (
	"encoding/json"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"time"
)

type TaskHandleService struct {
}

// CreateTaskHandleByTemplate  根据模板创建任务处理
func (s *TaskHandleService) CreateTaskHandleByTemplate(taskId, userToken, language string, request *models.RequestTable, taskTemplate *models.TaskTemplateTable, taskHandleTemplate *models.TaskHandleTemplateTable) (actions []*dao.ExecAction) {
	var approvalList []*models.TaskTemplateDto
	now := time.Now().Format(models.DateTimeFormat)
	actions = []*dao.ExecAction{}
	var action *dao.ExecAction
	// 角色管理员
	if taskTemplate.HandleMode == string(models.TaskTemplateHandleModeAdmin) {
		result, _ := GetRoleService().GetRoleAdministrators(request.Role, userToken, language)
		if len(result) > 0 && result[0] != "" {
			action = &dao.ExecAction{Sql: "insert into task_handle (id,task_handle_template,task,role,handler,handler_type,created_time,updated_time) values(?,?,?,?,?,?,?,?)"}
			action.Param = []interface{}{guid.CreateGuid(), taskHandleTemplate.Id, taskId, request.Role, result[0], taskHandleTemplate.HandlerType, now, now}
		}
		return
	}
	if taskHandleTemplate.Assign == string(models.TaskHandleTemplateAssignTypeTemplate) {

	} else if taskHandleTemplate.Assign == string(models.TaskHandleTemplateAssignTypeCustom) {
		// 提交人指定,则需要读取请求表,获取提交人指定审批
		if request.TaskApprovalCache != "" {
			json.Unmarshal([]byte(request.TaskApprovalCache), &approvalList)
			if len(approvalList) > 0 {
				for _, approval := range approvalList {
					if approval.Id == taskTemplate.Id && len(approval.HandleTemplates) > 0 {
						for _, handleTemplate := range approval.HandleTemplates {
							action = &dao.ExecAction{Sql: "insert into task_handle (id,task_handle_template,task,role,handler,handler_type,created_time,updated_time) values(?,?,?,?,?,?,?,?)"}
							action.Param = []interface{}{guid.CreateGuid(), taskHandleTemplate.Id, taskId, request.Role, handleTemplate.Role, taskHandleTemplate.HandlerType, now, now}
						}
					}
				}
			}
		}
	}
	return
}
