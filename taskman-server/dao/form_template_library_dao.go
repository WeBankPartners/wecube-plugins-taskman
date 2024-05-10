package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"time"
	"xorm.io/xorm"
)

type FormTemplateLibraryDao struct {
	DB *xorm.Engine
}

func (d *FormTemplateLibraryDao) Add(session *xorm.Session, formTemplateLibrary *models.FormTemplateLibraryTable) (affected int64, err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	affected, err = session.Insert(formTemplateLibrary)
	// 打印日志
	logExecuteSql(session, "FormTemplateLibraryDao", "Add", formTemplateLibrary, affected, err)
	return
}

func (d *FormTemplateLibraryDao) Delete(session *xorm.Session, id string) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	_, err = session.ID(id).Delete(&models.FormTemplateLibraryTable{})
	return
}

func (d *FormTemplateLibraryDao) Get(formTemplateLibraryId string) (*models.FormTemplateLibraryTable, error) {
	var formTemplateLibrary *models.FormTemplateLibraryTable
	var found bool
	var err error
	formTemplateLibrary = &models.FormTemplateLibraryTable{}
	found, err = d.DB.ID(formTemplateLibraryId).Get(formTemplateLibrary)
	if err != nil {
		return nil, err
	}
	if found {
		return formTemplateLibrary, nil
	}
	return nil, nil
}

func (d *FormTemplateLibraryDao) Disable(session *xorm.Session, id string) (err error) {
	var affected int64
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if id == "" {
		return
	}
	affected, err = session.ID(id).Update(&models.FormTemplateLibraryTable{DelFlag: 1, UpdatedTime: time.Now().Format(models.DateTimeFormat)})
	// 打印日志
	logExecuteSql(session, "FormTemplateLibraryDao", "Disable", id, affected, err)
	if err != nil {
		return
	}
	return
}
func (d *FormTemplateLibraryDao) QueryAll() (list []*models.FormTemplateLibraryTable, err error) {
	err = d.DB.Where("del_flag = 0").Find(&list)
	return
}

func (d *FormTemplateLibraryDao) QueryByName(name string) (list []*models.FormTemplateLibraryTable, err error) {
	err = d.DB.Where("name=? and del_flag = 0", name).Find(&list)
	return
}

func (d *FormTemplateLibraryDao) QueryListByCondition(condition models.QueryFormTemplateLibraryParam) (pageInfo models.PageInfo, list []*models.FormTemplateLibraryTable, err error) {
	var params []interface{}
	sql := "select * from form_template_library where del_flag = 0"
	if condition.FormType != "" {
		sql = sql + " and form_type = ?"
		params = append(params, condition.FormType)
	}
	if condition.CreatedBy != "" {
		sql = sql + " and created_by = ?"
		params = append(params, condition.CreatedBy)
	}
	if condition.Name != "" {
		sql = sql + " and name like '%" + condition.Name + "%'"
	}
	pageInfo.StartIndex = condition.StartIndex
	pageInfo.PageSize = condition.PageSize
	pageInfo.TotalRows = QueryCount(sql, params...)
	sql = sql + " limit ?,? "
	params = append(params, condition.StartIndex, condition.PageSize)
	err = d.DB.SQL(sql, params...).Find(&list)
	return
}
