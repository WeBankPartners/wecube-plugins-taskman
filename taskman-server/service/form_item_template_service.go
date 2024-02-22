package service

import (
	"fmt"

	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

const defaultCustomFormItemName = "taskman-custom-form"

type FormItemTemplateService struct {
	formItemTemplateDao dao.FormItemTemplateDao
}

func (s FormItemTemplateService) UpdateFormTemplateItemGroupConfig(param models.FormTemplateGroupConfigureDto) (err error) {
	var formItemTemplateList []*models.FormItemTemplateTable
	var insertItems, updateItems, deleteItems []*models.FormItemTemplateTable
	var systemItemExist, customItemExist bool
	var existMap = make(map[string]bool)
	if len(param.SystemItems) == 0 {
		param.SystemItems = []*models.ProcEntityAttributeObj{}
	}
	if len(param.CustomItems) == 0 {
		param.CustomItems = []*models.FormItemTemplateTable{}
	}
	// 自定义类型 单独处理
	if param.ItemGroupType == string(models.FormItemGroupTypeCustom) {
		return s.UpdateFormTemplateCustomItemGroupConfig(param)
	}
	// 1. 查询表单组是否存在，不存在则新增
	formItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplateAndItemGroupName(param.FormTemplateId, param.ItemGroupName)
	if err != nil {
		return
	}
	// 新增数据
	if len(formItemTemplateList) == 0 {
		if len(param.SystemItems) > 0 {
			for _, systemItem := range param.SystemItems {
				insertItems = append(insertItems, models.ConvertProcEntityAttributeObj2FormItemTemplate(param, systemItem))
			}
		}
		if len(param.CustomItems) > 0 {
			for _, customItem := range param.CustomItems {
				customItem.FormTemplate = param.FormTemplateId
				customItem.ElementType = string(models.FormItemElementTypeCalculate)
				insertItems = append(insertItems, customItem)
			}
		}
	} else {
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
				insertItems = append(insertItems, models.ConvertProcEntityAttributeObj2FormItemTemplate(param, systemItem))
			}
		}
		for _, customItem := range param.CustomItems {
			customItemExist = false
			existMap[customItem.Id] = true
			for _, formItemTemplate := range formItemTemplateList {
				if customItem.Id == formItemTemplate.Id {
					customItemExist = true
					updateItems = append(updateItems, customItem)
					break
				}
			}
			if !customItemExist {
				customItem.FormTemplate = param.FormTemplateId
				insertItems = append(insertItems, customItem)
			}
		}
		for _, formItemTemplate := range formItemTemplateList {
			if existMap[formItemTemplate.Id] == false && existMap[formItemTemplate.AttrDefId] == false {
				deleteItems = append(deleteItems, formItemTemplate)
			}
		}
	}
	if len(insertItems) > 0 || len(updateItems) > 0 || len(deleteItems) > 0 {
		err = transaction(func(session *xorm.Session) error {
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
	}
	return
}

func (s FormItemTemplateService) UpdateFormTemplateCustomItemGroupConfig(param models.FormTemplateGroupConfigureDto) (err error) {
	var formItemTemplateList, tempFormItemTemplateList []*models.FormItemTemplateTable
	var addFormItemTemplate *models.FormItemTemplateTable
	// 1. 查询表单组是否存在，不存在则新增
	formItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplateAndItemGroup(param.FormTemplateId, param.ItemGroup)
	if err != nil {
		return
	}
	if len(formItemTemplateList) == 0 {
		formItemTemplateList = []*models.FormItemTemplateTable{}
	}
	formItemTemplateMap := models.ConvertFormItemTemplateList2Map(formItemTemplateList)
	tempFormItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplateAndItemGroupName(param.FormTemplateId, param.ItemGroupName)
	if err != nil {
		return
	}

	if len(tempFormItemTemplateList) > 0 {
		if len(tempFormItemTemplateList) != len(formItemTemplateList) {
			err = fmt.Errorf("itemGroupName:%s has exisit", param.ItemGroupName)
			return
		}
		for _, tempFormItemTemplate := range tempFormItemTemplateList {
			if formItemTemplateMap[tempFormItemTemplate.Id] == nil {
				err = fmt.Errorf("itemGroupName:%s has exisit", param.ItemGroupName)
				return
			}
		}
	}
	if len(formItemTemplateList) > 0 {
		// 更新模板
		for _, formItemTemplate := range formItemTemplateList {
			formItemTemplate.ItemGroupName = param.ItemGroupName
			formItemTemplate.ItemGroupRule = param.ItemGroupRule
			formItemTemplate.ItemGroupType = param.ItemGroupType
			err = s.formItemTemplateDao.Update(nil, formItemTemplate)
			if err != nil {
				return err
			}
		}
	} else {
		// 新增自定义组,同时给一个初始化假数据(后续更新模板组时候删除掉)
		addFormItemTemplate = &models.FormItemTemplateTable{Id: guid.CreateGuid(), Name: defaultCustomFormItemName, ItemGroup: guid.CreateGuid(), ItemGroupName: param.ItemGroupName,
			ItemGroupType: param.ItemGroupType, ItemGroupRule: param.ItemGroupRule, FormTemplate: param.FormTemplateId}
		_, err = s.formItemTemplateDao.Add(nil, addFormItemTemplate)
	}
	return
}

func (s FormItemTemplateService) DeleteFormTemplateItemGroup(formTemplateId, itemGroupName string) (err error) {
	var formItemTemplateList []*models.FormItemTemplateTable
	formItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplateAndItemGroupName(formTemplateId, itemGroupName)
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
					s.formItemTemplateDao.Add(session, item)
				} else {
					err = s.formItemTemplateDao.Update(session, item)
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
	var formItemTemplateList []*models.FormItemTemplateTable
	var formItemTemplateGroupSortMap = make(map[string]int)
	formItemTemplateGroupSortMap = s.buildFormTemplateGroupSortMap(param.ItemGroupNameSort)
	formItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplate(param.FormTemplateId)
	if err != nil {
		return
	}
	if len(formItemTemplateList) > 0 {
		err = transaction(func(session *xorm.Session) error {
			for _, formItemTemplate := range formItemTemplateList {
				formItemTemplate.Sort = formItemTemplateGroupSortMap[formItemTemplate.ItemGroupName]
				err = s.formItemTemplateDao.Update(session, formItemTemplate)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	return
}

func (s FormItemTemplateService) CopyDataFormTemplateItemGroup(formTemplateId, itemGroupName string) (err error) {
	var formItemTemplateList []*models.FormItemTemplateTable
	// 1. 查询表单组是否存在，不存在则新增
	formItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplateAndItemGroupName(formTemplateId, itemGroupName)
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
