package service

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"strings"
)

// NotifyTaskExpireMail 定版/确认/任务/审批分配给“我”,快超时了或者超时发送邮件
func NotifyTaskExpireMail(task *models.TaskTable, expireObj *models.ExpireObj, userToken, language string) (err error) {
	var taskHandleList []*models.TaskHandleTable
	var request models.RequestTable
	var requestName string
	var mailList []string
	var roleTableList []*models.RoleTable
	var userInfo *models.UserDto
	if !checkMailEnable() {
		return
	}
	if err = dao.X.SQL("select * from task_handle where task = ?", task.Id).Find(&taskHandleList); err != nil {
		log.Logger.Error("NotifyTaskExpireMail query db err", log.Error(err))
		return
	}
	if len(taskHandleList) == 0 {
		return
	}
	if task.Request != "" {
		if request, err = GetSimpleRequest(task.Request); err != nil {
			return
		}
		requestName = request.Name
	}
	task.Name = getInternationalizationTaskName(task.Name, language)
	for _, taskHandle := range taskHandleList {
		if taskHandle.Handler != "" {
			if userInfo, err = GetRoleService().GetUserInfo(taskHandle.Handler, userToken, language); err != nil {
				return err
			}
			if userInfo == nil || strings.TrimSpace(userInfo.Email) == "" {
				log.Logger.Warn("NotifyTaskExpireMail,taskName receiver email is empty", log.String("requestName", requestName), log.String("taskName", task.Name), log.String("receiver", taskHandle.Handler))
				return
			}
			mailList = append(mailList, userInfo.Email)
		} else if taskHandle.Role != "" {
			roleTableList = append(roleTableList, &models.RoleTable{Id: taskHandle.Role})
		}
	}
	if len(roleTableList) > 0 {
		subMailList := GetRoleService().GetRoleMail(roleTableList, userToken, language)
		if len(subMailList) > 0 {
			mailList = append(mailList, subMailList...)
		}
	}

	var mailSubject = "[wecube] [Task transfer reminder] 【任务超时提醒】"
	var mailContent string
	if (expireObj.Percent >= 75) && (expireObj.DoingNotifyCount == 0) {
		expireObj.DoingNotifyCount++
		mailContent = fmt.Sprintf("The task assigned to you [Request: %s Task: %s] is about to expire and is valid until %s. Please click the link to process it", requestName, task.Name, task.ExpireTime)
		mailContent = mailContent + fmt.Sprintf("\n\n\n分配给您的任务[请求:%s-任务:%s]快过期了,有效期到%s,请点击链接处理", requestName, task.Name, task.ExpireTime)
	} else {
		expireObj.TimeoutNotifyCount++
		mailContent = fmt.Sprintf("The task assigned to you [Request: %s Task: %s] has expired. Please click the link to process it", requestName, task.Name)
		mailContent = mailContent + fmt.Sprintf("\n\n\n分配给您的任务[请求:%s-任务:%s]已过期,请点击链接尽快处理", requestName, task.Name)
	}
	mailContent = mailContent + fmt.Sprintf("\n%s/#/taskman/workbench", models.Config.WebUrl)
	log.Logger.Debug("NotifyTaskExpireMail", log.String("mailSubject", mailSubject), log.String("mailContent", mailContent), log.String("mailList", strings.Join(mailList, ",")))
	if err = models.MailSender.Send(mailSubject, mailContent, mailList); err != nil {
		log.Logger.Error("send notify email fail", log.Error(err))
		return
	}
	return nil
}

// NotifyTaskRoleMail 给角色发邮件
func NotifyTaskRoleMail(requestName, taskName, expireDate, role, userToken, language string) (err error) {
	var mailList []string
	if !checkMailEnable() {
		return
	}
	var subject, content string
	var displayNameMap map[string]string
	taskName = getInternationalizationTaskName(taskName, language)
	if displayNameMap, err = GetRoleService().GetRoleDisplayName(userToken, language); err != nil {
		return
	}
	subject = "[wecube] [New Task Reminder]  【新增任务提醒】"
	content = fmt.Sprintf("Role %s has a pending task [Request: %s Task: %s], which is valid until %s. Please process it as soon as possible. Click to view details", displayNameMap[role], requestName, taskName, expireDate)
	content = content + fmt.Sprintf("\n\n\n角色%s有一条待处理任务[请求:%s-任务:%s],有效期截止到%s,请尽快处理,点击查看详情", displayNameMap[role], requestName, taskName, expireDate)
	content = content + fmt.Sprintf("\n%s/#/taskman/workbench", models.Config.WebUrl)
	mailList = GetRoleService().GetRoleMail([]*models.RoleTable{{Id: role}}, userToken, language)
	if len(mailList) > 0 {
		log.Logger.Debug("NotifyTaskRoleMail", log.String("mailSubject", subject), log.String("mailContent", content), log.String("mailList", mailList[0]))
		err = models.MailSender.Send(subject, content, mailList)
		if err != nil {
			log.Logger.Error("send mail err", log.Error(err))
		}
	}
	return
}

// NotifyTaskAssignMail 定版/确认/任务/审批分配给“我”
func NotifyTaskAssignMail(requestName, taskName, expireDate, receiver, userToken, language string) (err error) {
	if !checkMailEnable() {
		return
	}
	var subject, content string
	var userInfo *models.UserDto
	if userInfo, err = GetRoleService().GetUserInfo(receiver, userToken, language); err != nil {
		return err
	}
	if userInfo == nil || strings.TrimSpace(userInfo.Email) == "" {
		log.Logger.Warn("NotifyTaskAssignMail,taskName receiver email is empty", log.String("requestName", requestName), log.String("taskName", taskName), log.String("receiver", receiver))
		return
	}
	taskName = getInternationalizationTaskName(taskName, language)
	subject = "[wecube] [New Task reminder]  【新增任务提醒】"
	content = fmt.Sprintf("You have a pending task [Request: %s Task: %s], which is valid until %s. Please process it as soon as possible (if you are unable to process it, team members can transfer the task to another order for processing). Click to view details", requestName, taskName, expireDate)
	content = content + fmt.Sprintf("\n\n\n您有一条待处理任务[请求:%s-任务:%s],有效期截止到%s,请尽快处理(若本人无法处理,组员可以将任务转单处理),点击查看详情", requestName, taskName, expireDate)
	content = content + fmt.Sprintf("\n%s/#/taskman/workbench", models.Config.WebUrl)
	log.Logger.Debug("NotifyTaskAssignMail", log.String("mailSubject", subject), log.String("mailContent", content), log.String("mailList", userInfo.Email))
	err = models.MailSender.Send(subject, content, []string{userInfo.Email})
	if err != nil {
		log.Logger.Error("send mail err", log.Error(err))
	}
	return
}

// NotifyTaskHandlerUpdateMail 定版/确认/任务/审批分配给“我”,但是被人点“转给我”抢单了
func NotifyTaskHandlerUpdateMail(requestName, taskName, originHandler, targetHandler, userToken, language string) (err error) {
	if !checkMailEnable() {
		return
	}
	var subject, content string
	var userInfo *models.UserDto
	if userInfo, err = GetRoleService().GetUserInfo(originHandler, userToken, language); err != nil {
		return err
	}
	if userInfo == nil || strings.TrimSpace(userInfo.Email) == "" {
		log.Logger.Warn("NotifyTaskHandlerUpdateMail,taskName receiver email is empty", log.String("requestName", requestName), log.String("taskName", taskName), log.String("receiver", originHandler))
		return
	}
	taskName = getInternationalizationTaskName(taskName, language)
	subject = "[wecube] [Task transfer reminder]  【任务被转单提醒】"
	content = fmt.Sprintf("The task assigned to you [Request: %s Task: %s] has been transferred to %s. Click the link to view details", requestName, taskName, targetHandler)
	content = content + fmt.Sprintf("\n\n\n分配给您的任务[请求:%s-任务:%s]已被转单给%s,点击链接查看详情", requestName, taskName, targetHandler)
	content = content + fmt.Sprintf("\n%s/#/taskman/workbench", models.Config.WebUrl)
	log.Logger.Debug("NotifyTaskHandlerUpdateMail", log.String("mailSubject", subject), log.String("mailContent", content), log.String("mailList", userInfo.Email))
	err = models.MailSender.Send(subject, content, []string{userInfo.Email})
	if err != nil {
		log.Logger.Error("send mail err", log.Error(err))
	}
	return
}

// NotifyRequestCompleteMail 我提交的请求处理完成了
func NotifyRequestCompleteMail(requestName, creator, userToken, language string) (err error) {
	if !checkMailEnable() {
		return
	}
	var subject, content string
	var userInfo *models.UserDto
	if userInfo, err = GetRoleService().GetUserInfo(creator, userToken, language); err != nil {
		return err
	}
	if userInfo == nil || strings.TrimSpace(userInfo.Email) == "" {
		log.Logger.Warn("NotifyRequestCompleteMail,requestName creator email is empty", log.String("requestName", requestName), log.String("creator", creator))
		return
	}
	subject = "[wecube] [Request completion reminder]  【请求完成提醒】"
	content = fmt.Sprintf("The [request: %s] you initiated has been processed. Click on the link to view details", requestName)
	content = content + fmt.Sprintf("\n\n\n您发起的[请求:%s]已处理完成,点击链接查看详情", requestName)
	content = content + fmt.Sprintf("\n%s/#/taskman/workbench", models.Config.WebUrl)
	log.Logger.Debug("NotifyRequestCompleteMail", log.String("mailSubject", subject), log.String("mailContent", content), log.String("mailList", userInfo.Email))

	err = models.MailSender.Send(subject, content, []string{userInfo.Email})
	if err != nil {
		log.Logger.Error("send mail err", log.Error(err))
	}
	return
}

// NotifyTaskBackMail 我提交的请求被定版退回/审批退回
func NotifyTaskBackMail(requestName, taskName, creator, approval, userToken, language string) (err error) {
	if !checkMailEnable() {
		return
	}
	var subject, content string
	var userInfo *models.UserDto
	if userInfo, err = GetRoleService().GetUserInfo(creator, userToken, language); err != nil {
		return err
	}
	if userInfo == nil || strings.TrimSpace(userInfo.Email) == "" {
		log.Logger.Warn("NotifyTaskBackMail,requestName creator email is empty", log.String("requestName", requestName), log.String("creator", creator))
		return
	}
	taskName = getInternationalizationTaskName(taskName, language)
	subject = "[wecube] [Request completion reminder]  【请求退回提醒】"
	content = fmt.Sprintf("The [request: %s] you initiated was returned to the draft by %s at node %s. Please make the necessary modifications and resubmit. Click the link to view details", requestName, taskName, approval)
	content = content + fmt.Sprintf("\n\n\n您发起的[请求:%s],在%s节点被%s退回到草稿,请修改之后重新提交,点击链接查看详情", requestName, taskName, approval)
	content = content + fmt.Sprintf("\n%s/#/taskman/workbench", models.Config.WebUrl)
	err = models.MailSender.Send(subject, content, []string{userInfo.Email})
	if err != nil {
		log.Logger.Error("send mail err", log.Error(err))
	}
	return
}

// NotifyTaskDenyMail 我提交的请求被审批拒绝
func NotifyTaskDenyMail(requestName, taskName, creator, approval, userToken, language string) (err error) {
	if !checkMailEnable() {
		return
	}
	var subject, content string
	var userInfo *models.UserDto
	if userInfo, err = GetRoleService().GetUserInfo(creator, userToken, language); err != nil {
		return err
	}
	if userInfo == nil || strings.TrimSpace(userInfo.Email) == "" {
		log.Logger.Warn("NotifyTaskDenyMail,requestName creator email is empty", log.String("requestName", requestName), log.String("creator", creator))
		return
	}
	taskName = getInternationalizationTaskName(taskName, language)
	subject = "[wecube] [Request termination reminder]  【请求终止提醒】"
	content = fmt.Sprintf("The [request: %s] you initiated was rejected by %s at the %s approval node, and the request has been terminated. Please click the link to view details", requestName, taskName, approval)
	content = content + fmt.Sprintf("\n\n\n您发起的[请求:%s],在%s审批节点被%s拒绝,请求已终止,请点击链接查看详情", requestName, taskName, approval)
	content = content + fmt.Sprintf("\n%s/#/taskman/workbench", models.Config.WebUrl)
	log.Logger.Debug("NotifyTaskDenyMail", log.String("mailSubject", subject), log.String("mailContent", content), log.String("mailList", userInfo.Email))

	err = models.MailSender.Send(subject, content, []string{userInfo.Email})
	if err != nil {
		log.Logger.Error("send mail err", log.Error(err))
	}
	return
}

// NotifyTaskWorkflowFailMail 我提交的请求在编排执行中被手动终止
func NotifyTaskWorkflowFailMail(requestName, procDefName, status, creator, userToken, language string) (err error) {
	if !checkMailEnable() {
		return
	}
	var subject, content string
	var userInfo *models.UserDto
	if userInfo, err = GetRoleService().GetUserInfo(creator, userToken, language); err != nil {
		return err
	}
	if userInfo == nil || strings.TrimSpace(userInfo.Email) == "" {
		log.Logger.Warn("NotifyTaskWorkflowTerminationMail,requestName creator email is empty", log.String("requestName", requestName), log.String("creator", creator))
		return
	}
	subject = "[wecube] [Request Termination Reminder]  【请求终止提醒】"
	if status == string(models.RequestStatusTermination) {
		content = fmt.Sprintf("Due to scheduling %s reaching the [Auto Exit] node, it caused Your [Request: %s] request has been terminated. Please click on the link to view details", procDefName, requestName)
		content = content + fmt.Sprintf("\n\n\n因为编排%s被管理员[手动终止],导致\n您发起的[请求:%s]请求已终止,请点击链接查看详情", procDefName, requestName)
		content = content + fmt.Sprintf("\n%s/#/taskman/workbench", models.Config.WebUrl)
	} else if status == string(models.RequestStatusFaulted) {
		content = fmt.Sprintf("Due to scheduling %s being manually terminated by the administrator, it resulted in The [Request: %s] request you initiated has been terminated. Please click on the link to view details", procDefName, requestName)
		content = content + fmt.Sprintf("\n\n\n因为编排%s走到[自动退出]节点,导致\n您发起的[请求:%s]请求终止,请点击链接查看详情", procDefName, requestName)
		content = content + fmt.Sprintf("\n%s/#/taskman/workbench", models.Config.WebUrl)
	}
	log.Logger.Debug("NotifyTaskWorkflowFailMail", log.String("mailSubject", subject), log.String("mailContent", content), log.String("mailList", userInfo.Email))

	err = models.MailSender.Send(subject, content, []string{userInfo.Email})
	if err != nil {
		log.Logger.Error("send mail err", log.Error(err))
	}
	return
}

func getInternationalizationTaskName(taskName, language string) string {
	switch taskName {
	case RequestPending:
		taskName = "Check(定版)"
	case "check":
		taskName = "Check(定版)"
	case Confirm:
		taskName = "Confirm(确认)"
	}
	return taskName
}

func checkMailEnable() bool {
	return models.MailEnable
}
