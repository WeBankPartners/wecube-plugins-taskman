package service

import (
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

const defaultCustomFormItemName = "taskman-custom-form"

type FormItemTemplateService struct {
	formItemTemplateDao      dao.FormItemTemplateDao
	formItemTemplateGroupDao dao.FormItemTemplateGroupDao
}

func (s FormItemTemplateService) UpdateFormTemplateItemGroupConfig(param models.FormTemplateGroupConfigureDto) (err error) {
	var formItemTemplateList []*models.FormItemTemplateTable
	var insertItems, updateItems, deleteItems []*models.FormItemTemplateTable
	var updateItemGroup *models.FormItemTemplateGroupTable
	var newItemGroupId string
	var systemItemExist, customItemExist bool
	var existMap = make(map[string]bool)
	if len(param.SystemItems) == 0 {
		param.SystemItems = []*models.ProcEntityAttributeObj{}
	}
	if len(param.CustomItems) == 0 {
		param.CustomItems = []*models.FormItemTemplateDto{}
	}
	// 1. 查询表单组是否存在，不存在则新增
	formItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplateAndItemGroupId(param.FormTemplateId, param.ItemGroupId)
	if err != nil {
		return
	}
	// 新增数据
	if len(formItemTemplateList) == 0 {
		newItemGroupId = guid.CreateGuid()
		if len(param.SystemItems) > 0 {
			for _, systemItem := range param.SystemItems {
				insertItems = append(insertItems, models.ConvertProcEntityAttributeObj2FormItemTemplate(param, systemItem, newItemGroupId))
			}
		}
		if len(param.CustomItems) > 0 {
			for _, customItem := range param.CustomItems {
				customItem.FormTemplate = param.FormTemplateId
				customItem.ElementType = string(models.FormItemElementTypeCalculate)
				insertItems = append(insertItems, models.ConvertFormItemTemplateDto2Model(customItem))
			}
		}
	} else {
		updateItemGroup, err = s.formItemTemplateGroupDao.Get(param.ItemGroupId)
		if err != nil {
			return
		}
		if updateItemGroup != nil {
			updateItemGroup.ItemGroupName = param.ItemGroupName
			updateItemGroup.ItemGroup = param.ItemGroup
			updateItemGroup.ItemGroupRule = param.ItemGroupRule
			updateItemGroup.ItemGroupType = param.ItemGroupType
		}
		for _, systemItem := range param.SystemItems {
			systemItemExist = false
			existMap[systemItem.Id] = true
			for _, formItemTemplate := range formItemTemplateList {
				if systemItem.Id == formItemTemplate.AttrDefId {
					systemItemExist = true
					break
				}
			}
			if !systemItemExist {
				insertItems = append(insertItems, models.ConvertProcEntityAttributeObj2FormItemTemplate(param, systemItem, param.ItemGroupId))
			}
		}
		for _, customItem := range param.CustomItems {
			customItemExist = false
			existMap[customItem.Id] = true
			for _, formItemTemplate := range formItemTemplateList {
				if customItem.Id == formItemTemplate.Id {
					customItemExist = true
					updateItems = append(updateItems, models.ConvertFormItemTemplateDto2Model(customItem))
					break
				}
			}
			if !customItemExist {
				customItem.FormTemplate = param.FormTemplateId
				insertItems = append(insertItems, models.ConvertFormItemTemplateDto2Model(customItem))
			}
		}
		for _, formItemTemplate := range formItemTemplateList {
			if existMap[formItemTemplate.Id] == false && existMap[formItemTemplate.AttrDefId] == false {
				deleteItems = append(deleteItems, formItemTemplate)
			}
		}
		if param.ItemGroupType == string(models.FormItemGroupTypeCustom) {

		}
	}
	err = transaction(func(session *xorm.Session) error {
		if newItemGroupId != "" {
			_, err = s.formItemTemplateGroupDao.Add(session, &models.FormItemTemplateGroupTable{
				Id:            newItemGroupId,
				ItemGroup:     param.ItemGroup,
				ItemGroupType: param.ItemGroupType,
				ItemGroupName: param.ItemGroupName,
				ItemGroupSort: param.ItemGroupSort,
				ItemGroupRule: param.ItemGroupRule,
				FormTemplate:  param.FormTemplateId,
			})
		}
		if updateItemGroup != nil {
			err = s.formItemTemplateGroupDao.Update(session, updateItemGroup)
		}
		if len(insertItems) > 0 {
			for _, item := range insertItems {
				_, err = s.formItemTemplateDao.Add(session, item)
				if err != nil {
					return err
				}
			}
		}
		if len(updateItems) > 0 {
			for _, item := range updateItems {
				err = s.formItemTemplateDao.Update(session, item)
				if err != nil {
					return err
				}
			}
		}
		if len(deleteItems) > 0 {
			for _, item := range deleteItems {
				err = s.formItemTemplateDao.Delete(session, item.Id)
				if err != nil {
					return err
				}
			}
		}
		return err
	})
	return
}

func (s FormItemTemplateService) DeleteFormTemplateItemGroup(formTemplateId, itemGroupId string) (err error) {
	var formItemTemplateList []*models.FormItemTemplateTable
	formItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplateAndItemGroupId(formTemplateId, itemGroupId)
	if err != nil {
		return err
	}
	if len(formItemTemplateList) > 0 {
		err = transaction(func(session *xorm.Session) error {
			for _, formItemTemplate := range formItemTemplateList {
				err = s.formItemTemplateDao.DeleteByIdOrCopyId(session, formItemTemplate.Id)
				if err != nil {
					return err
				}
				err = s.formItemTemplateGroupDao.DeleteByIdOrCopyId(session, itemGroupId)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	return
}

func (s FormItemTemplateService) UpdateFormTemplateItemGroup(param models.FormTemplateGroupCustomDataDto) (err error) {
	if len(param.Items) > 0 {
		err = transaction(func(session *xorm.Session) error {
			for _, item := range param.Items {
				// Id 为空新增
				if item.Id == "" {
					item.Id = guid.CreateGuid()
					item.FormTemplate = param.FormTemplateId
					s.formItemTemplateDao.Add(session, models.ConvertFormItemTemplateDto2Model(item))
				} else {
					err = s.formItemTemplateDao.Update(session, models.ConvertFormItemTemplateDto2Model(item))
				}
				if err != nil {
					return err
				}
			}
			return err
		})
	}
	return
}

func (s FormItemTemplateService) SortFormTemplateItemGroup(param models.FormTemplateGroupSortDto) (err error) {
	var formItemTemplateList []*models.FormItemTemplateDto
	var formItemTemplateGroupSortMap = make(map[string]int)
	formItemTemplateGroupSortMap = s.buildFormTemplateGroupSortMap(param.ItemGroupIdSort)
	formItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplate(param.FormTemplateId)
	if err != nil {
		return
	}
	if len(formItemTemplateList) > 0 {
		err = transaction(func(session *xorm.Session) error {
			for _, formItemTemplate := range formItemTemplateList {
				formItemTemplate.Sort = formItemTemplateGroupSortMap[formItemTemplate.ItemGroupName]
				err = s.formItemTemplateDao.Update(session, models.ConvertFormItemTemplateDto2Model(formItemTemplate))
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	return
}

func (s FormItemTemplateService) CopyDataFormTemplateItemGroup(formTemplateId, itemGroupId string) (err error) {
	var formItemTemplateList []*models.FormItemTemplateTable
	// 1. 查询表单组是否存在，不存在则新增
	formItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplateAndItemGroupId(formTemplateId, itemGroupId)
	if err != nil {
		return
	}
	// 新增数据
	if len(formItemTemplateList) > 0 {
		err = transaction(func(session *xorm.Session) error {
			for _, formItemTemplate := range formItemTemplateList {
				formItemTemplate.CopyId = formItemTemplate.Id
				formItemTemplate.Id = guid.CreateGuid()
				_, err = s.formItemTemplateDao.Add(session, formItemTemplate)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	return
}

func (s FormItemTemplateService) buildFormTemplateGroupSortMap(itemGroupNameSort []string) map[string]int {
	hashMap := make(map[string]int)
	for i, groupName := range itemGroupNameSort {
		hashMap[groupName] = i
	}
	return hashMap
}
