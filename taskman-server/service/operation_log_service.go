package service

import (
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"go.uber.org/zap"
	"time"
)

type OperationLogService struct {
	operationLogDao *dao.OperationLogDao
}

func (s *OperationLogService) RecordRequestTemplateLog(requestTemplateId, requestTemplateName, operator, operation, uri, content string) {
	_, err := dao.X.Exec("insert into operation_log(id,request_template,request_template_name,operation,uri,content,operator,op_time) values (?,?,?,?,?,?,?,?)",
		guid.CreateGuid(), requestTemplateId, requestTemplateName, operation, uri, content, operator, time.Now().Format(models.DateTimeFormat))
	if err != nil {
		log.Error(nil, log.LOGGER_APP, "Record request operation log fail", zap.Error(err))
	}
}

func (s *OperationLogService) RecordRequestLog(requestId, requestName, operator, operation, uri, content string) {
	_, err := dao.X.Exec("insert into operation_log(id,request,request_name,operation,uri,content,operator,op_time) values (?,?,?,?,?,?,?,?)",
		guid.CreateGuid(), requestId, requestName, operation, uri, content, operator, time.Now().Format(models.DateTimeFormat))
	if err != nil {
		log.Error(nil, log.LOGGER_APP, "Record request operation log fail", zap.Error(err))
	}
}

func (s *OperationLogService) RecordTaskLog(taskId, taskName, operator, operation, uri, content string) {
	_, err := dao.X.Exec("insert into operation_log(id,task,task_name,operation,uri,content,operator,op_time) values (?,?,?,?,?,?,?,?)",
		guid.CreateGuid(), taskId, taskName, operation, uri, content, operator, time.Now().Format(models.DateTimeFormat))
	if err != nil {
		log.Error(nil, log.LOGGER_APP, "Record request operation log fail", zap.Error(err))
	}
}
