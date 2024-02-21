package service

import (
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

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
				insertItems = append(insertItems, customItem)
			}
		}
		for _, formItemTemplate := range formItemTemplateList {
			if existMap[formItemTemplate.Id] == false && existMap[formItemTemplate.AttrDefId] == false {
				deleteItems = append(deleteItems, formItemTemplate)
			}
		}
	}
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
	return
}

func (s FormItemTemplateService) buildFormTemplateGroupSortMap(itemGroupNameSort []string) map[string]int {
	hashMap := make(map[string]int)
	for i, groupName := range itemGroupNameSort {
		hashMap[groupName] = i
	}
	return hashMap
}
