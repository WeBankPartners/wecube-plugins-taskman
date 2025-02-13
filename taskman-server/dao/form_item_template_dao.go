package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type FormItemTemplateDao struct {
	DB *xorm.Engine
}

func (d *FormItemTemplateDao) Add(session *xorm.Session, formItemTemplate *models.FormItemTemplateTable) (affected int64, err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	affected, err = session.Insert(formItemTemplate)
	// 打印日志
	logExecuteSql(session, "FormItemTemplateDao", "Add", formItemTemplate, affected, err)
	return
}

func (d *FormItemTemplateDao) Update(session *xorm.Session, formItemTemplate *models.FormItemTemplateTable) (err error) {
	var affected int64
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if formItemTemplate.Id == "" {
		return
	}
	affected, err = session.ID(formItemTemplate.Id).AllCols().Update(formItemTemplate)
	// 打印日志
	logExecuteSql(session, "FormItemTemplateDao", "Update", formItemTemplate, affected, err)
	return
}

func (d *FormItemTemplateDao) UpdateCmdbAttribute(session *xorm.Session, formItemTemplateId, cmdbAttr string) (err error) {
	var affected int64
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if formItemTemplateId == "" {
		return
	}
	affected, err = session.ID(formItemTemplateId).Cols("cmdb_attr").Update(cmdbAttr)
	// 打印日志
	logExecuteSql(session, "FormItemTemplateDao", "UpdateCmdbAttribute", formItemTemplateId, affected, err)
	return
}

func (d *FormItemTemplateDao) UpdateByRefId(session *xorm.Session, formItemTemplate *models.FormItemTemplateTable, refId string) (err error) {
	var affected int64
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	affected, err = session.Where("ref_id = ?", refId).Update(formItemTemplate)
	// 打印日志
	logExecuteSql(session, "FormItemTemplateDao", "UpdateByRefId", formItemTemplate, affected, err)
	return
}

func (d *FormItemTemplateDao) Get(formItemTemplateId string) (*models.FormItemTemplateTable, error) {
	var formItemTemplate *models.FormItemTemplateTable
	var found bool
	var err error
	formItemTemplate = &models.FormItemTemplateTable{}
	found, err = d.DB.ID(formItemTemplateId).Get(formItemTemplate)
	if err != nil {
		return nil, err
	}
	if found {
		return formItemTemplate, nil
	}
	return nil, nil
}

func (d *FormItemTemplateDao) QueryDtoByFormTemplate(formTemplate string) (formItemTemplateDtoList []*models.FormItemTemplateDto, err error) {
	var formItemTemplateList []*models.FormItemTemplateTable
	var formItemTemplateGroup models.FormTemplateTable
	formItemTemplateDtoList = []*models.FormItemTemplateDto{}
	if formTemplate == "" {
		return
	}
	err = d.DB.Where("form_template = ?", formTemplate).OrderBy("sort").Find(&formItemTemplateList)
	if err != nil {
		return
	}
	if len(formItemTemplateList) > 0 && formItemTemplateList[0].FormTemplate != "" {
		_, err = d.DB.ID(formItemTemplateList[0].FormTemplate).Get(&formItemTemplateGroup)
		if err != nil {
			return
		}
	}
	for _, formItemTemplate := range formItemTemplateList {
		dto := models.ConvertFormItemTemplateModel2Dto(formItemTemplate, formItemTemplateGroup)
		if dto != nil {
			formItemTemplateDtoList = append(formItemTemplateDtoList, dto)
		}
	}
	return
}

func (d *FormItemTemplateDao) QueryByFormTemplate(formTemplate string) (formItemTemplate []*models.FormItemTemplateTable, err error) {
	formItemTemplate = []*models.FormItemTemplateTable{}
	err = d.DB.Where("form_template = ?", formTemplate).Find(&formItemTemplate)
	return
}

func (d *FormItemTemplateDao) DeleteByIdOrRefId(session *xorm.Session, id string) (err error) {
	var affected int64
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if id == "" {
		return
	}
	affected, err = session.Where("id = ? or ref_id = ?", id, id).Delete(&models.FormItemTemplateTable{})
	// 打印日志
	logExecuteSql(session, "FormItemTemplateDao", "DeleteByIdOrRefId", id, affected, err)
	return
}

func (d *FormItemTemplateDao) Delete(session *xorm.Session, id string) (err error) {
	var affected int64
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if id == "" {
		return
	}
	affected, err = session.ID(id).Delete(&models.FormItemTemplateTable{})
	// 打印日志
	logExecuteSql(session, "FormItemTemplateDao", "Delete", id, affected, err)
	return
}

func (d *FormItemTemplateDao) DeleteByFormTemplate(session *xorm.Session, formTemplate string) (err error) {
	var affected int64
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if formTemplate == "" {
		return
	}
	affected, err = session.Where("form_template = ?", formTemplate).Delete(&models.FormItemTemplateTable{})
	// 打印日志
	logExecuteSql(session, "FormItemTemplateDao", "DeleteByFormTemplate", formTemplate, affected, err)
	return
}

func (d *FormItemTemplateDao) QueryByRefId(refId string) (formItemTemplateList []*models.FormItemTemplateTable, err error) {
	err = d.DB.Where("ref_id = ?", refId).Find(&formItemTemplateList)
	return
}
