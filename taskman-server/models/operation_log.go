package models

import (
	"github.com/WeBankPartners/go-common-lib/guid"
	"time"
)

type OperationLogTable struct {
	Id                  string `json:"id" xorm:"id"`
	Request             string `json:"request" xorm:"request"`
	RequestName         string `json:"requestName" xorm:"request_name"`
	RequestTemplate     string `json:"requestTemplate" xorm:"request_template"`
	RequestTemplateName string `json:"requestTemplateName" xorm:"request_template_name"`
	Task                string `json:"task" xorm:"task"`
	TaskName            string `json:"taskName" xorm:"task_name"`
	Operation           string `json:"operation" xorm:"operation"`
	Operator            string `json:"operator" xorm:"operator"`
	Content             string `json:"content" xorm:"content"`
	Uri                 string `json:"uri" xorm:"uri"`
	OpTime              string `json:"opTime" xorm:"op_time"`
}

func (OperationLogTable) TableName() string {
	return "operation_log"
}

func CreateRecordRequestTemplateLog(requestTemplateId, requestTemplateName, operator, operation, uri, content string) *OperationLogTable {
	return &OperationLogTable{
		Id:                  guid.CreateGuid(),
		RequestTemplate:     requestTemplateId,
		RequestTemplateName: requestTemplateName,
		Operation:           operation,
		Operator:            operator,
		Content:             content,
		Uri:                 uri,
		OpTime:              time.Now().Format(DateTimeFormat),
	}
}
