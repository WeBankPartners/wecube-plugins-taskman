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
	taskName = getInternationalizationTaskName(taskName, language)
	subject = "[wecube] [Task transfer reminder] +【任务被转单提醒】"
	content = fmt.Sprintf("您有一条待处理任务[请求:%s-任务:%s],有效期截止到%s,请尽快处理(若本人无法处理,组员可以将任务转单处理),点击查看详情", requestName, taskName, expireDate)
	content = content + fmt.Sprintf("\n\n\nYou have a pending task [Request: %s Task: %s], which is valid until %s. Please process it as soon as possible (if you are unable to process it, team members can transfer the task to another order for processing). Click to view details", requestName, taskName, expireDate)
	err = models.MailSender.Send(subject, content, []string{userInfo.EmailAddr})
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
	taskName = getInternationalizationTaskName(taskName, language)
	subject = "[wecube] [New Task Reminder] +【新增任务提醒】"
	content = fmt.Sprintf("角色%s有一条待处理任务[请求:%s-任务:%s],有效期截止到%s,请尽快处理,点击查看详情", displayNameMap[role], requestName, taskName, expireDate)
	content = content + fmt.Sprintf("\n\n\nRole %s has a pending task [Request: %s Task: %s], which is valid until %s. Please process it as soon as possible. Click to view details", displayNameMap[role], requestName, taskName, expireDate)
	err = models.MailSender.Send(subject, content, []string{userInfo.EmailAddr})
	return
}

// NotifyTaskHandlerUpdateMail 定版/确认/任务/审批分配给“我”,但是被人点“转给我”抢单了
func NotifyTaskHandlerUpdateMail(requestName, taskName, originHandler, userToken, language string) (err error) {
	var subject, content string
	var userInfo *models.SimpleLocalUserDto
	if userInfo, err = GetRoleService().GetUserInfo(originHandler, userToken, language); err != nil {
		return err
	}
	if userInfo == nil || strings.TrimSpace(userInfo.EmailAddr) == "" {
		log.Logger.Warn("NotifyTaskHandlerUpdateMail,taskName receiver email is empty", log.String("requestName", requestName), log.String("taskName", taskName), log.String("receiver", originHandler))
		return
	}
	subject = "[wecube] [Task transfer reminder] +【任务被转单提醒】"
	content = fmt.Sprintf("分配给您的任务[请求:%s-任务:%s]已被转单给%s,点击链接查看详情", requestName, taskName, originHandler)
	content = content + fmt.Sprintf("\n\n\nThe task assigned to you [Request: %s Task: %s] has been transferred to %s. Click the link to view details", requestName, taskName, originHandler)
	err = models.MailSender.Send(subject, content, []string{userInfo.EmailAddr})
	return
}

// NotifyRequestCompleteMail 我提交的请求处理完成了
func NotifyRequestCompleteMail(requestName, creator, userToken, language string) (err error) {
	var subject, content string
	var userInfo *models.SimpleLocalUserDto
	if userInfo, err = GetRoleService().GetUserInfo(creator, userToken, language); err != nil {
		return err
	}
	if userInfo == nil || strings.TrimSpace(userInfo.EmailAddr) == "" {
		log.Logger.Warn("NotifyRequestCompleteMail,requestName creator email is empty", log.String("requestName", requestName), log.String("creator", creator))
		return
	}
	subject = "[wecube] [Request completion reminder] +【请求完成提醒】"
	content = fmt.Sprintf("您发起的[请求:%s]已处理完成,点击链接查看详情", requestName)
	content = content + fmt.Sprintf("\n\n\nThe [request: %s] you initiated has been processed. Click on the link to view details", requestName)
	err = models.MailSender.Send(subject, content, []string{userInfo.EmailAddr})
	return
}

// NotifyTaskBackMail 我提交的请求被定版退回/审批退回
func NotifyTaskBackMail(requestName, taskName, creator, approval, userToken, language string) (err error) {
	var subject, content string
	var userInfo *models.SimpleLocalUserDto
	if userInfo, err = GetRoleService().GetUserInfo(creator, userToken, language); err != nil {
		return err
	}
	if userInfo == nil || strings.TrimSpace(userInfo.EmailAddr) == "" {
		log.Logger.Warn("NotifyTaskBackMail,requestName creator email is empty", log.String("requestName", requestName), log.String("creator", creator))
		return
	}
	taskName = getInternationalizationTaskName(taskName, language)
	subject = "[wecube] [Request completion reminder] +【请求退回提醒】"
	content = fmt.Sprintf("您发起的[请求:%s],在%s节点被%s退回到草稿,请修改之后重新提交,点击链接查看详情", requestName, taskName, approval)
	content = content + fmt.Sprintf("The [request: %s] you initiated was returned to the draft by %s at node %s. Please make the necessary modifications and resubmit. Click the link to view details", requestName, taskName, approval)
	err = models.MailSender.Send(subject, content, []string{userInfo.EmailAddr})
	return
}

// NotifyTaskDenyMail 我提交的请求被审批拒绝
func NotifyTaskDenyMail(requestName, taskName, creator, approval, userToken, language string) (err error) {
	var subject, content string
	var userInfo *models.SimpleLocalUserDto
	if userInfo, err = GetRoleService().GetUserInfo(creator, userToken, language); err != nil {
		return err
	}
	if userInfo == nil || strings.TrimSpace(userInfo.EmailAddr) == "" {
		log.Logger.Warn("NotifyTaskDenyMail,requestName creator email is empty", log.String("requestName", requestName), log.String("creator", creator))
		return
	}
	taskName = getInternationalizationTaskName(taskName, language)
	subject = "[wecube] [Request termination reminder] +【请求终止提醒】"
	content = fmt.Sprintf("您发起的[请求:%s],在%s审批节点被%s拒绝,请求已终止,请点击链接查看详情", requestName, taskName, approval)
	content = content + fmt.Sprintf("\n\n\nThe [request: %s] you initiated was rejected by %s at the %s approval node, and the request has been terminated. Please click the link to view details", requestName, taskName, approval)
	err = models.MailSender.Send(subject, content, []string{userInfo.EmailAddr})
	return
}

// NotifyTaskWorkflowFailMail 我提交的请求在编排执行中被手动终止
func NotifyTaskWorkflowFailMail(requestName, procDefName, status, creator, userToken, language string) (err error) {
	var subject, content string
	var userInfo *models.SimpleLocalUserDto
	if userInfo, err = GetRoleService().GetUserInfo(creator, userToken, language); err != nil {
		return err
	}
	if userInfo == nil || strings.TrimSpace(userInfo.EmailAddr) == "" {
		log.Logger.Warn("NotifyTaskWorkflowTerminationMail,requestName creator email is empty", log.String("requestName", requestName), log.String("creator", creator))
		return
	}
	subject = "[wecube] [Request Termination Reminder] +【请求终止提醒】"
	if status == string(models.RequestStatusTermination) {
		content = fmt.Sprintf("因为编排%s被管理员[手动终止],导致\n您发起的[请求:%s]请求已终止,请点击链接查看详情", procDefName, requestName)
		content = content + fmt.Sprintf("\n\n\nDue to scheduling %s reaching the [Auto Exit] node, it caused Your [Request: %s] request has been terminated. Please click on the link to view details", procDefName, requestName)
	} else if status == string(models.RequestStatusFaulted) {
		content = fmt.Sprintf("因为编排%s走到[自动退出]节点,导致\n您发起的[请求:%s]请求终止,请点击链接查看详情", procDefName, requestName)
		content = content + fmt.Sprintf("\n\n\nDue to scheduling %s being manually terminated by the administrator, it resulted in The [Request: %s] request you initiated has been terminated. Please click on the link to view details", procDefName, requestName)
	}
	err = models.MailSender.Send(subject, content, []string{userInfo.EmailAddr})
	return
}

func getInternationalizationTaskName(taskName, language string) string {
	switch taskName {
	case RequestPending:
		taskName = "Pending(定版)"
	case Confirm:
		taskName = "Confirm(确认)"
	}
	return taskName
}
