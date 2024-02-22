package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type FormItemTemplateDao struct {
	DB *xorm.Engine
}

func (d FormItemTemplateDao) Add(session *xorm.Session, formItemTemplate *models.FormItemTemplateTable) (affected int64, err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	affected, err = session.Insert(formItemTemplate)
	// 打印日志
	logExecuteSql(session, "FormItemTemplateDao", "Add", formItemTemplate, affected, err)
	return
}

func (d FormItemTemplateDao) Update(session *xorm.Session, formItemTemplate *models.FormItemTemplateTable) (err error) {
	var affected int64
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if formItemTemplate.Id == "" {
		return
	}
	affected, err = session.ID(formItemTemplate.Id).Update(formItemTemplate)
	// 打印日志
	logExecuteSql(session, "FormItemTemplateDao", "Update", formItemTemplate, affected, err)
	return
}

func (d FormItemTemplateDao) Get(formItemTemplateId string) (*models.FormItemTemplateTable, error) {
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

func (d FormItemTemplateDao) QueryByFormTemplate(formTemplate string) (formItemTemplateDtoList []*models.FormItemTemplateDto, err error) {
	var formItemTemplateList []*models.FormItemTemplateTable
	var formItemTemplateGroup models.FormItemTemplateGroupTable
	formItemTemplateDtoList = []*models.FormItemTemplateDto{}
	if formTemplate == "" {
		return
	}
	err = d.DB.Where("form_template = ?", formTemplate).Find(&formItemTemplateList)
	if err != nil {
		return
	}
	if len(formItemTemplateList) > 0 && formItemTemplateList[0].ItemGroupId != "" {
		_, err = d.DB.ID(formItemTemplateList[0].ItemGroupId).Get(&formItemTemplateGroup)
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

func (d FormItemTemplateDao) QueryByFormTemplateAndItemGroupName(formTemplate, itemGroupName string) (formItemTemplate []*models.FormItemTemplateTable, err error) {
	formItemTemplate = []*models.FormItemTemplateTable{}
	err = d.DB.Where("form_template = ? and item_group_name = ?", formTemplate, itemGroupName).Find(&formItemTemplate)
	return
}

func (d FormItemTemplateDao) QueryDtoByFormTemplateAndItemGroupName(formTemplate, itemGroupName string) (formItemTemplateDtoList []*models.FormItemTemplateDto, err error) {
	var formItemTemplateList []*models.FormItemTemplateTable
	var formItemTemplateGroup models.FormItemTemplateGroupTable
	formItemTemplateDtoList = []*models.FormItemTemplateDto{}
	if formTemplate == "" {
		return
	}
	err = d.DB.Where("form_template = ? and item_group_name = ?", formTemplate, itemGroupName).Find(&formItemTemplateList)
	if len(formItemTemplateList) > 0 && formItemTemplateList[0].ItemGroupId != "" {
		_, err = d.DB.ID(formItemTemplateList[0].ItemGroupId).Get(&formItemTemplateGroup)
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

func (d FormItemTemplateDao) QueryByFormTemplateAndItemGroupId(formTemplate, itemGroupId string) (formItemTemplate []*models.FormItemTemplateTable, err error) {
	formItemTemplate = []*models.FormItemTemplateTable{}
	err = d.DB.Where("form_template = ? and item_group_id = ?", formTemplate, itemGroupId).Find(&formItemTemplate)
	return
}

func (d FormItemTemplateDao) QueryDtoByFormTemplateAndItemGroupId(formTemplate, itemGroupId string) (formItemTemplateDtoList []*models.FormItemTemplateDto, err error) {
	var formItemTemplateList []*models.FormItemTemplateTable
	var formItemTemplateGroup models.FormItemTemplateGroupTable
	formItemTemplateDtoList = []*models.FormItemTemplateDto{}
	if formTemplate == "" {
		return
	}
	err = d.DB.Where("form_template = ? and item_group_id = ?", formTemplate, itemGroupId).Find(&formItemTemplateList)
	if len(formItemTemplateList) > 0 && formItemTemplateList[0].ItemGroupId != "" {
		_, err = d.DB.ID(formItemTemplateList[0].ItemGroupId).Get(&formItemTemplateGroup)
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

func (d FormItemTemplateDao) QueryByFormTemplateAndItemGroup(formTemplate, itemGroup string) (formItemTemplate []*models.FormItemTemplateTable, err error) {
	formItemTemplate = []*models.FormItemTemplateTable{}
	err = d.DB.Where("form_template = ? and item_group = ?", formTemplate, itemGroup).Find(&formItemTemplate)
	return
}

func (d FormItemTemplateDao) DeleteByIdOrCopyId(session *xorm.Session, id string) (err error) {
	var affected int64
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if id == "" {
		return
	}
	affected, err = session.Where("id = ? or copy_id = ?", id, id).Delete(&models.FormItemTemplateTable{})
	// 打印日志
	logExecuteSql(session, "FormItemTemplateDao", "DeleteByIdOrCopyId", id, affected, err)
	return
}

func (d FormItemTemplateDao) Delete(session *xorm.Session, id string) (err error) {
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
