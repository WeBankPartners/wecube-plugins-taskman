package service

import (
	"encoding/json"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/rpc"
	"math/rand"
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
	// 保存任务审批不为空,解析任务审批
	if request.TaskApprovalCache != "" {
		json.Unmarshal([]byte(request.TaskApprovalCache), &approvalList)
		if len(approvalList) > 0 {
			for _, approval := range approvalList {
				if approval.Id == taskTemplate.Id && len(approval.HandleTemplates) > 0 {
					for _, handleTemplate := range approval.HandleTemplates {
						// 组内系统分配,随机给一个
						if handleTemplate.HandlerType == string(models.TaskHandleTemplateHandlerTypeSystem) {
							if handleTemplate.Role != "" {
								userList, err := rpc.QueryRolesUsers(handleTemplate.Role, userToken, language)
								if err != nil {
									log.Logger.Error("rpcQueryRolesUsers fail", log.Error(err))
								}
								if len(userList) > 0 {
									rand.Seed(time.Now().UnixNano())
									handleTemplate.Handler = userList[rand.Intn(len(userList))].UserName
								}
							}
						}
						action = &dao.ExecAction{Sql: "insert into task_handle (id,task_handle_template,task,role,handler,handler_type,created_time,updated_time) values(?,?,?,?,?,?,?,?)"}
						action.Param = []interface{}{guid.CreateGuid(), taskHandleTemplate.Id, taskId, request.Role, handleTemplate.Handler, taskHandleTemplate.HandlerType, now, now}
					}
				}
			}
		}
	}
	return
}
