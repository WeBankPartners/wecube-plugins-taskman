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
func (s *TaskHandleService) CreateTaskHandleByTemplate(taskId, userToken, language string, request *models.RequestTable, taskTemplate *models.TaskTemplateTable) (actions []*dao.ExecAction) {
	var taskTemplateDtoList []*models.TaskTemplateDto
	now := time.Now().Format(models.DateTimeFormat)
	actions = []*dao.ExecAction{}

	// 保存任务审批不为空,解析任务审批
	if request.TaskApprovalCache != "" {
		json.Unmarshal([]byte(request.TaskApprovalCache), &taskTemplateDtoList)
		if len(taskTemplateDtoList) > 0 {
			for _, taskTemplateDto := range taskTemplateDtoList {
				if taskTemplateDto.Id == taskTemplate.Id && len(taskTemplateDto.HandleTemplates) > 0 {
					// 角色管理员
					if taskTemplate.HandleMode == string(models.TaskTemplateHandleModeAdmin) {
						result, _ := GetRoleService().GetRoleAdministrators(request.Role, userToken, language)
						if len(result) > 0 && result[0] != "" {
							action := &dao.ExecAction{Sql: "insert into task_handle (id,task,role,handler,created_time,updated_time) values(?,?,?,?,?,?)"}
							action.Param = []interface{}{guid.CreateGuid(), taskId, request.Role, result[0], now, now}
						}
						continue
					}
					for _, handleTemplate := range taskTemplateDto.HandleTemplates {
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
						action := &dao.ExecAction{Sql: "insert into task_handle (id,task_handle_template,task,role,handler,handler_type,created_time,updated_time) values(?,?,?,?,?,?,?,?)"}
						action.Param = []interface{}{guid.CreateGuid(), handleTemplate.Id, taskId, handleTemplate.Role, handleTemplate.Handler, handleTemplate.HandlerType, now, now}
						actions = append(actions, action)
					}
				}
			}
		}
	}
	return
}

func (s *TaskHandleService) GetRequestCheckTaskHandle(taskId string) (taskHandle *models.TaskHandleTable, err error) {
	var taskHandleList []*models.TaskHandleTable
	err = dao.X.SQL("select * from task_handle where task = ?", taskId).Find(&taskHandleList)
	if err != nil {
		return
	}
	if len(taskHandleList) > 0 {
		taskHandle = taskHandleList[0]
	}
	return
}

func (s *TaskHandleService) Get(id string) (taskHandle *models.TaskHandleTable, err error) {
	var taskHandleList []*models.TaskHandleTable
	err = dao.X.SQL("select * from task_handle where id = ?", id).Find(&taskHandleList)
	if err != nil {
		return
	}
	if len(taskHandleList) > 0 {
		taskHandle = taskHandleList[0]
	}
	return
}
