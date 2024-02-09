package service

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"xorm.io/xorm"
)

var (
	// db操作
	db *xorm.Engine
	// 文件上传 service
	attachFileService AttachFileService
	// 模板收藏 service
	collectTemplateService CollectTemplateService
	// 表单 service
	formService FormService
	// 表单模板 service
	formTemplateService FormTemplateService
	// 日志记录 service
	operationLogService OperationLogService
	// 编排 service
	procDefService ProcDefService
	// CMDB service
	refSelectService RefSelectService
	// 请求 service
	requestService RequestService
	// 请求模板 service
	requestTemplateService RequestTemplateService
	// 请求模板组 service
	requestTemplateGroupService RequestTemplateGroupService
	// 任务 service
	taskService TaskService
	// 任务模板 service
	taskTemplateService TaskTemplateService
)

func New() (err error) {
	var engine *xorm.Engine
	// 初始化Db
	engine, err = dao.InitDatabase()
	if err != nil {
		return
	}
	// 初始化Dao
	attachFileDao := dao.AttachFileDao{DB: engine}
	collectTemplateDao := dao.CollectTemplateDao{DB: engine}
	formDao := dao.FormDao{DB: engine}
	formItemDao := dao.FormItemDao{DB: engine}
	formItemTemplateDao := dao.FormItemTemplateDao{DB: engine}
	operationLogDao := dao.OperationLogDao{DB: engine}
	formTemplateDao := dao.FormTemplateDao{DB: engine}
	requestDao := dao.RequestDao{DB: engine}
	requestTemplateDao := dao.RequestTemplateDao{DB: engine}
	requestTemplateGroupDao := dao.RequestTemplateGroupDao{DB: engine}
	requestTemplateRoleDao := dao.RequestTemplateRoleDao{DB: engine}
	taskDao := dao.TaskDao{DB: engine}
	taskTemplateDao := dao.TaskTemplateDao{DB: engine}
	taskTemplateRoleDao := dao.TaskTemplateRoleDao{DB: engine}
	// 初始化Service
	collectTemplateService = CollectTemplateService{collectTemplateDao: collectTemplateDao}
	requestService = RequestService{requestDao: requestDao}
	attachFileService = AttachFileService{attachFileDao: attachFileDao}
	operationLogService = OperationLogService{operationLogDao: operationLogDao}
	procDefService = ProcDefService{}
	refSelectService = RefSelectService{}
	requestService = RequestService{requestDao: requestDao}
	taskService = TaskService{taskDao: taskDao}
	taskTemplateService = TaskTemplateService{taskTemplateDao: taskTemplateDao, taskTemplateRoleDao: taskTemplateRoleDao}
	formService = FormService{formDao: formDao, formItemDao: formItemDao}
	requestTemplateService = RequestTemplateService{requestTemplateDao: requestTemplateDao, operationLogDao: operationLogDao, requestTemplateRoleDao: requestTemplateRoleDao}
	requestTemplateGroupService = RequestTemplateGroupService{requestTemplateGroupDao: requestTemplateGroupDao}
	formTemplateService = FormTemplateService{formTemplateDao: formTemplateDao, formItemTemplateDao: formItemTemplateDao, formDao: formDao}
	db = engine
	return
}

// GetRequestService 获取请求service
func GetRequestService() RequestService {
	return requestService
}

// GetRequestTemplateService 获取请求模板 service
func GetRequestTemplateService() RequestTemplateService {
	return requestTemplateService
}

// GetOperationLogService 获取记录日志 service
func GetOperationLogService() OperationLogService {
	return operationLogService
}

// GetProcDefService 获取编排 service
func GetProcDefService() ProcDefService {
	return procDefService
}

// GetAttachFileService 文件上传 service
func GetAttachFileService() AttachFileService {
	return attachFileService
}

// GetCollectTemplateService 模板收藏 service
func GetCollectTemplateService() CollectTemplateService {
	return collectTemplateService
}

// GetFormService 表单 service
func GetFormService() FormService {
	return formService
}

// GetFormTemplateService 表单模板 service
func GetFormTemplateService() FormTemplateService {
	return formTemplateService
}

// GetRefSelectService service
func GetRefSelectService() RefSelectService {
	return refSelectService
}

// GetTaskService 任务 service
func GetTaskService() TaskService {
	return taskService
}

// GetTaskTemplateService 获取任务模板 service
func GetTaskTemplateService() TaskTemplateService {
	return taskTemplateService
}

// GetRequestTemplateGroupService 获取请求模板组 service
func GetRequestTemplateGroupService() RequestTemplateGroupService {
	return requestTemplateGroupService
}

// transaction 事务处理
func transaction(f func(session *xorm.Session) error) (err error) {
	session := db.NewSession()
	defer session.Close()
	session.Begin()
	err = f(session)
	if err != nil {
		session.Rollback()
	}
	err = session.Commit()
	if err != nil {
		return
	}
	return
}

// transactionWithoutForeignCheck 事务处理
func transactionWithoutForeignCheck(f func(session *xorm.Session) error) (err error) {
	session := db.NewSession()
	defer session.Close()
	session.Begin()
	session.Exec("SET FOREIGN_KEY_CHECKS=0")
	err = f(session)
	if err != nil {
		session.Rollback()
	}
	err = session.Commit()
	if err != nil {
		return
	}
	session.Exec("SET FOREIGN_KEY_CHECKS=1")
	return
}
