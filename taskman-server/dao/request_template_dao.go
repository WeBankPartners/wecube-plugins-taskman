package dao

import (
	"database/sql"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type RequestTemplateDao struct {
	DB *xorm.Engine
}

// Add 添加模板
func (d RequestTemplateDao) Add(session *xorm.Session, requestTemplate *models.RequestTemplateTable) (affected int64, err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	affected, err = session.Insert(requestTemplate)
	// 打印日志
	logExecuteSql(session, "RequestTemplateDao", "Add", requestTemplate, affected, err)
	return
}

// AddBasicInfo 添加模板基础信息(此处用SQL形式添加,由于RequestTemplateTable中包含外键字段,外键form_template传递"",新增数据会报错)
func (d RequestTemplateDao) AddBasicInfo(session *xorm.Session, template *models.RequestTemplateTable) (affected int64, err error) {
	var result sql.Result
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	result, err = session.Exec("insert into request_template(id,`group`,name,description,tags,package_name,entity_name,proc_def_key,proc_def_id,"+
		"proc_def_name,expire_day,handler,created_by,created_time,updated_by,updated_time,type,operator_obj_type,parent_id,approve_by,pending_switch,"+
		"pending_role,pending_handler,confirm_switch,confirm_expire_day,rollback_desc) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", template.Id,
		template.Group, template.Name, template.Description, template.Tags, template.PackageName, template.EntityName, template.ProcDefKey, template.ProcDefId,
		template.ProcDefName, template.ExpireDay, template.Handler, template.CreatedBy, template.CreatedTime, template.UpdatedBy, template.UpdatedTime, template.OperatorObjType,
		template.Id, template.ApproveBy, template.PendingSwitch, template.PendingRole, template.PendingHandler, template.ConfirmSwitch, template.ConfirmExpireDay, template.RollbackDesc)
	if err != nil {
		return
	}
	affected, err = result.RowsAffected()
	return
}

func (d RequestTemplateDao) Update(session *xorm.Session, requestTemplate *models.RequestTemplateTable) (err error) {
	var affected int64
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if requestTemplate == nil || requestTemplate.Id == "" {
		return
	}
	// 由于RequestTemplateTable里面包含version字段,此处需要去掉xorm自带版本校验
	session.NoVersionCheck()
	affected, err = session.ID(requestTemplate.Id).Update(requestTemplate)
	// 打印日志
	logExecuteSql(session, "RequestTemplateDao", "Update", requestTemplate, affected, err)
	if err != nil {
		return
	}
	return
}

func (d RequestTemplateDao) Get(requestTemplateId string) (*models.RequestTemplateTable, error) {
	var requestTemplate *models.RequestTemplateTable
	var found bool
	var err error
	requestTemplate = &models.RequestTemplateTable{}
	found, err = d.DB.ID(requestTemplateId).Get(requestTemplate)
	if err != nil {
		return nil, err
	}
	if found {
		return requestTemplate, nil
	}
	return nil, nil
}
