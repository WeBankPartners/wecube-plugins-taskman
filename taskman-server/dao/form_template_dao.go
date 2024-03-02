package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type FormTemplateDao struct {
	DB *xorm.Engine
}

func (d *FormTemplateDao) Add(session *xorm.Session, formTemplate *models.FormTemplateTable) (affected int64, err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	affected, err = session.Insert(formTemplate)
	// 打印日志
	logExecuteSql(session, "FormTemplateDao", "Add", formTemplate, affected, err)
	return
}

func (d *FormTemplateDao) Update(session *xorm.Session, formTemplate *models.FormTemplateTable) (err error) {
	var affected int64
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if formTemplate == nil || formTemplate.Id == "" {
		return
	}
	affected, err = session.ID(formTemplate.Id).Update(formTemplate)
	// 打印日志
	logExecuteSql(session, "FormTemplateDao", "Update", formTemplate, affected, err)
	if err != nil {
		return
	}
	return
}

func (d *FormTemplateDao) Get(formTemplateId string) (*models.FormTemplateTable, error) {
	var formTemplate *models.FormTemplateTable
	var found bool
	var err error
	formTemplate = &models.FormTemplateTable{}
	found, err = d.DB.ID(formTemplateId).Get(formTemplate)
	if err != nil {
		return nil, err
	}
	if found {
		return formTemplate, nil
	}
	return nil, nil
}

func (d *FormTemplateDao) Delete(session *xorm.Session, id string) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	_, err = d.DB.ID(id).Delete(&models.FormTemplateTable{})
	return
}

func (d *FormTemplateDao) QueryListByRequestTemplateAndTaskTemplate(requestTemplateId, taskTemplateId, requestFormType string) (list []*models.FormTemplateTable, err error) {
	list = []*models.FormTemplateTable{}
	// taskTemplateId 为空,则只查询请求表单,根据 requestFormType区分类型
	if taskTemplateId == "" {
		err = d.DB.Where("request_template=? and request_form_type = ? and del_flag = 0", requestTemplateId, requestFormType).Find(&list)
	} else {
		err = d.DB.Where("request_template=?  and task_template=? and del_flag = 0", requestTemplateId, taskTemplateId).Find(&list)
	}
	return
}

// QueryRequestFormByRequestTemplateIdAndType 查询请求表单
func (d *FormTemplateDao) QueryRequestFormByRequestTemplateIdAndType(requestTemplateId, requestFormType string) (result *models.FormTemplateTable, err error) {
	var list []*models.FormTemplateTable
	err = d.DB.Where("request_template=? and request_form_type = ? and del_flag = 0", requestTemplateId, requestFormType).Find(&list)
	if len(list) > 0 {
		result = list[0]
	}
	return
}

func (d *FormTemplateDao) QueryListByRequestTemplateAndItemGroupType(requestTemplateId, itemGroupType string) (list []*models.FormTemplateTable, err error) {
	list = []*models.FormTemplateTable{}
	err = d.DB.Where("request_template=? and  item_group_type = ? and del_flag = 0", requestTemplateId, itemGroupType).Find(&list)
	return
}

func (d *FormTemplateDao) QueryListByIdOrRefId(id string) (list []*models.FormTemplateTable, err error) {
	list = []*models.FormTemplateTable{}
	err = d.DB.Where("id=? or ref_id=?", id, id).Find(&list)
	return
}
