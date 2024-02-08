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
	// 表单模板 service
	formTemplateService FormTemplateService
	// CMDB service
	refSelectService RefSelectService
	// 请求 service
	requestService RequestService
	// 请求模板 service
	requestTemplateService RequestTemplateService
	// 任务 service
	taskService TaskService
	// 记录 service
	operationLogService OperationLogService
	// 编排 service
	procDefService ProcDefService
)

func New() (err error) {
	var engine *xorm.Engine
	// 初始化Db
	engine, err = dao.InitDatabase()
	if err != nil {
		return
	}
	// 初始化Dao
	requestDao := dao.RequestDao{DB: engine}
	attachFileDao := dao.AttachFileDao{DB: engine}
	requestTemplateDao := dao.RequestTemplateDao{DB: engine}
	requestTemplateRoleDao := dao.RequestTemplateRoleDao{DB: engine}
	operationLogDao := dao.OperationLogDao{DB: engine}
	// 初始化Service
	requestService = RequestService{requestDao: requestDao}
	attachFileService = AttachFileService{attachFileDao: attachFileDao}
	requestTemplateService = RequestTemplateService{requestTemplateDao: requestTemplateDao, requestTemplateRoleDao: requestTemplateRoleDao, operationLogDao: operationLogDao}
	operationLogService = OperationLogService{operationLogDao: operationLogDao}
	procDefService = ProcDefService{}
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

// transaction 事务处理
func transaction(f func(session *xorm.Session) error) (err error) {
	session := db.NewSession()
	defer session.Close()
	session.Begin()
	err = f(session)
	if err != nil {
		session.Rollback()
	}
	session.Commit()
	return
}
