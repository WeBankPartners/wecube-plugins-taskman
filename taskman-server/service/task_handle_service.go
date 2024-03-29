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
							actions = append(actions, &dao.ExecAction{Sql: "insert into task_handle (id,task,role,handler,created_time,updated_time) values(?,?,?,?,?,?)",
								Param: []interface{}{guid.CreateGuid(), taskId, request.Role, result[0], now, now}})
							NotifyTaskAssignMail(request.Name, taskTemplate.Name, calcExpireTime(now, taskTemplate.ExpireDay), result[0], userToken, language)
						} else {
							// 没有找到角色管理员,用本组兜底
							actions = append(actions, &dao.ExecAction{Sql: "insert into task_handle (id,task,role,created_time,updated_time) values(?,?,?,?,?,?)",
								Param: []interface{}{guid.CreateGuid(), taskId, request.Role, now, now}})
						}
						continue
					}
					for _, handleTemplate := range taskTemplateDto.HandleTemplates {
						// 组内系统分配,随机给一个
						if handleTemplate.HandlerType == string(models.TaskHandleTemplateHandlerTypeSystem) {
							if handleTemplate.Role != "" {
								var roleTableList []*models.RoleTable
								//将 roleName =>roleId
								dao.X.SQL("select * from `role` where id =?", handleTemplate.Role).Find(&roleTableList)
								if len(roleTableList) > 0 {
									userList, err := rpc.QueryRolesUsers(roleTableList[0].CoreId, userToken, language)
									if err != nil {
										log.Logger.Error("rpcQueryRolesUsers fail", log.Error(err))
									}
									if len(userList) > 0 {
										rand.Seed(time.Now().UnixNano())
										handleTemplate.Handler = userList[rand.Intn(len(userList))].UserName
									}
									NotifyTaskAssignMail(request.Name, taskTemplate.Name, calcExpireTime(now, taskTemplate.ExpireDay), handleTemplate.Handler, userToken, language)
								}
							}
						} else if handleTemplate.HandlerType == string(models.TaskHandleTemplateHandlerTypeClaim) {
							//  组内认领,给角色发送邮件
							NotifyTaskRoleMail(request.Name, taskTemplate.Name, calcExpireTime(now, taskTemplate.ExpireDay), handleTemplate.Role, userToken, language)
						} else {
							NotifyTaskAssignMail(request.Name, taskTemplate.Name, calcExpireTime(now, taskTemplate.ExpireDay), handleTemplate.Handler, userToken, language)
						}
						actions = append(actions, &dao.ExecAction{Sql: "insert into task_handle (id,task_handle_template,task,role,handler,handler_type,created_time,updated_time) values(?,?,?,?,?,?,?,?)",
							Param: []interface{}{guid.CreateGuid(), handleTemplate.Id, taskId, handleTemplate.Role, handleTemplate.Handler, handleTemplate.HandlerType, now, now}})
					}
				}
			}
		}
	}
	return
}

func (s *TaskHandleService) GetRequestCheckTaskHandle(taskId string) (taskHandle *models.TaskHandleTable, err error) {
	var taskHandleList []*models.TaskHandleTable
	err = dao.X.SQL("select * from task_handle where task = ? and latest_flag = 1", taskId).Find(&taskHandleList)
	if err != nil {
		return
	}
	if len(taskHandleList) > 0 {
		taskHandle = taskHandleList[0]
	}
	return
}

func (s *TaskHandleService) GetTaskHandleListByTaskId(taskId string) (taskHandleList []*models.TaskHandleTable, err error) {
	err = dao.X.SQL("select * from task_handle where task = ? and latest_flag = 1", taskId).Find(&taskHandleList)
	return
}

func (s *TaskHandleService) Get(id string) (taskHandle *models.TaskHandleTable, err error) {
	var taskHandleList []*models.TaskHandleTable
	err = dao.X.SQL("select * from task_handle where id = ? and latest_flag = 1", id).Find(&taskHandleList)
	if err != nil {
		return
	}
	if len(taskHandleList) > 0 {
		taskHandle = taskHandleList[0]
	}
	return
}

func (s *TaskHandleService) GetLatestRequestCheckTaskHandleByRequestId(requestId string) (taskHandle *models.TaskHandleTable, err error) {
	var taskList []*models.TaskTable
	var taskHandleList []*models.TaskHandleTable
	if requestId == "" {
		return
	}
	// 可能会有多次定版,取最新一次
	err = dao.X.SQL("select * from task   where request = ? and type = ? order by sort desc", requestId, models.TaskTypeCheck).Find(&taskList)
	if err != nil {
		return
	}
	if len(taskList) > 0 {
		dao.X.SQL("select * from task_handle   where task = ? and latest_flag = 1", taskList[0].Id).Find(&taskHandleList)
	}
	if len(taskHandleList) > 0 {
		taskHandle = taskHandleList[0]
		return
	}
	return
}

// CalcTaskResult 计算任务处理结果,当前处理节点选择 完成才调用 CalcTaskResult
func (s *TaskHandleService) CalcTaskResult(taskId, curTaskHandleId string) string {
	var result = string(models.TaskHandleResultTypeComplete)
	var taskHandleList []*models.TaskHandleTable
	var taskTemplateList []*models.TaskTemplateTable
	var handleMode string
	if taskId == "" || curTaskHandleId == "" {
		return ""
	}
	// 协同情况特殊处理
	dao.X.SQL("select * from task_template where id in (select task_template from task where id = ?)", taskId).Find(&taskTemplateList)
	if len(taskTemplateList) > 0 {
		handleMode = taskTemplateList[0].HandleMode
	}
	if handleMode == string(models.TaskTemplateHandleModeAny) {
		return result
	}

	dao.X.SQL("select * from task_handle  where task = ? and latest_flag = 1", taskId).Find(&taskHandleList)
	if len(taskHandleList) > 0 {
		for _, taskHandle := range taskHandleList {
			if taskHandle.Id == curTaskHandleId {
				continue
			}
			if taskHandle.HandleStatus != string(models.TaskHandleResultTypeComplete) {
				result = string(models.TaskHandleResultTypeUncompleted)
				return result
			}
		}
	}
	return result
}
