package service

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"strings"
)

func NotifyTaskMail(taskId, userToken, language, mailSubject, mailContent string) error {
	if !models.MailEnable {
		return nil
	}
	taskObj, _ := getSimpleTask(taskId)
	log.Logger.Info("Start notify task mail", log.String("taskId", taskId))
	var roleTable []*models.RoleTable
	dao.X.SQL("select id,email from `role` where id in (select `role` from task_template_role where role_type='USE' and task_template in (select task_template from task where id=?))", taskId).Find(&roleTable)
	reportRoleString := taskObj.ReportRole
	reportRoleString = strings.ReplaceAll(reportRoleString, "[", "")
	reportRoleString = strings.ReplaceAll(reportRoleString, "]", "")
	for _, v := range strings.Split(reportRoleString, ",") {
		if v != "" {
			roleTable = append(roleTable, &models.RoleTable{Id: v})
		}
	}
	if len(roleTable) == 0 {
		return fmt.Errorf("can not find handle role with task:%s ", taskId)
	}
	mailList := GetRoleService().GetRoleMail(roleTable, userToken, language)
	if len(mailList) == 0 {
		log.Logger.Warn("Notify task mail break,email is empty", log.String("role", roleTable[0].Id))
		return fmt.Errorf("handle role email is empty ")
	}
	var taskTable []*models.TaskTable
	dao.X.SQL("select t1.id,t1.name,t1.description,t2.name as request,t1.node_name,t1.emergency,t1.reporter,t1.created_time from task t1 left join request t2 on t1.request=t2.id where t1.id=?", taskId).Find(&taskTable)
	if len(taskTable) == 0 {
		return fmt.Errorf("can not find task with id:%s ", taskId)
	}
	var subject, content string
	subject = fmt.Sprintf("Taskman task [%s] %s[%s]", models.PriorityLevelMap[taskTable[0].Emergency], taskTable[0].Name, taskTable[0].Request)
	content = fmt.Sprintf("Taskman task \nID:%s \nPriority:%s \nName:%s \nRequest:%s \nDescription:%s \nReporter:%s \nCreateTime:%s \n", taskTable[0].Id, models.PriorityLevelMap[taskTable[0].Emergency], taskTable[0].Name, taskTable[0].Request, taskTable[0].Description, taskTable[0].Reporter, taskTable[0].CreatedTime)
	if mailSubject != "" {
		subject = mailSubject
	}
	if mailContent != "" {
		content = mailContent
	}
	err := models.MailSender.Send(subject, content, mailList)
	if err != nil {
		return fmt.Errorf("send notify email fail:%s ", err.Error())
	}
	return nil
}

// NotifyTaskAssignMail 定版/确认/任务/审批分配给“我”
func NotifyTaskAssignMail(requestName, taskName, expireDate, receiver, userToken, language string) (err error) {
	var subject, content string
	var userInfo *models.SimpleLocalUserDto
	if userInfo, err = GetRoleService().GetUserInfo(receiver, userToken, language); err != nil {
		return err
	}
	if userInfo == nil || strings.TrimSpace(userInfo.EmailAddr) == "" {
		log.Logger.Warn("NotifyTaskAssignMail,taskName receiver email is empty", log.String("requestName", requestName), log.String("taskName", taskName), log.String("receiver", receiver))
		return
	}
	switch taskName {
	case RequestPending:
		taskName = "请求定版"
	}
	subject = "【任务被转单提醒】"
	content = fmt.Sprintf("您有一条待处理任务[请求:%s-任务:%s],有效期截止到%s,请尽快处理(若本人无法处理,组员可以将任务转单处理),点击查看详情", requestName, taskName, expireDate)
	err = models.MailSender.Send(subject, content, []string{userInfo.EmailAddr})
	if err != nil {
		return
	}
	return
}

// NotifyTaskRoleAdministratorMail 定版/确认/任务/审批分配给一个组,处理人为空,我是管理员
func NotifyTaskRoleAdministratorMail(requestName, taskName, expireDate, role, userToken, language string) (err error) {
	var subject, content string
	var result []string
	var userInfo *models.SimpleLocalUserDto
	var displayNameMap map[string]string
	if result, err = GetRoleService().GetRoleAdministrators(role, userToken, language); err != nil {
		return err
	}
	if len(result) == 0 {
		log.Logger.Warn("NotifyTaskAssignMail,taskName role administrator is empty", log.String("requestName", requestName), log.String("taskName", taskName), log.String("role", role))
		return
	}
	if userInfo, err = GetRoleService().GetUserInfo(result[0], userToken, language); err != nil {
		return err
	}
	if userInfo == nil || strings.TrimSpace(userInfo.EmailAddr) == "" {
		log.Logger.Warn("NotifyTaskRoleAdministratorMail,taskName receiver email is empty", log.String("requestName", requestName), log.String("taskName", taskName), log.String("receiver", result[0]))
		return
	}
	if displayNameMap, err = GetRoleService().GetRoleDisplayName(); err != nil {
		return
	}
	switch taskName {
	case RequestPending:
		taskName = "请求定版"
	}
	subject = "【新增任务提醒】"
	content = fmt.Sprintf("角色%s有一条待处理任务[请求:%s-任务:%s],有效期截止到%s,请尽快处理,点击查看详情", displayNameMap[role], requestName, taskName, expireDate)
	err = models.MailSender.Send(subject, content, []string{userInfo.EmailAddr})
	if err != nil {
		return
	}
	return
}
