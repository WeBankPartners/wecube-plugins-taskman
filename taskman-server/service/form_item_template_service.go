package service

import (
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"time"
	"xorm.io/xorm"
)

type FormItemTemplateService struct {
	formItemTemplateDao      *dao.FormItemTemplateDao
	formItemTemplateGroupDao *dao.FormItemTemplateGroupDao
}

func (s FormItemTemplateService) UpdateFormTemplateItemGroupConfig(param models.FormTemplateGroupConfigureDto) (err error) {
	var formItemTemplateGroup *models.FormTemplateNewTable
	var insertItems, updateItems, deleteItems []*models.FormItemTemplateTable
	var formItemTemplateList []*models.FormItemTemplateTable
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
	formItemTemplateGroup, err = s.formItemTemplateGroupDao.Get(param.FormTemplateId)
	if err != nil {
		return
	}
	// 新增数据
	if formItemTemplateGroup == nil {
		newItemGroupId = guid.CreateGuid()
		if len(param.SystemItems) > 0 {
			for _, systemItem := range param.SystemItems {
				insertItems = append(insertItems, models.ConvertProcEntityAttributeObj2FormItemTemplate(param, systemItem, newItemGroupId))
			}
		}
		if len(param.CustomItems) > 0 {
			for _, customItem := range param.CustomItems {
				customItem.Id = guid.CreateGuid()
				customItem.FormTemplate = newItemGroupId
				customItem.ItemGroup = param.ItemGroup
				customItem.ItemGroupName = param.ItemGroupName
				customItem.ElementType = string(models.FormItemElementTypeCalculate)
				insertItems = append(insertItems, models.ConvertFormItemTemplateDto2Model(customItem))
			}
		}
	} else {
		// 直接更新表单组
		formItemTemplateGroup.ItemGroupName = param.ItemGroupName
		formItemTemplateGroup.ItemGroup = param.ItemGroup
		formItemTemplateGroup.ItemGroupRule = param.ItemGroupRule
		formItemTemplateGroup.ItemGroupType = param.ItemGroupType

		formItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplate(param.FormTemplateId)
		if err != nil {
			return
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
				insertItems = append(insertItems, models.ConvertProcEntityAttributeObj2FormItemTemplate(param, systemItem, param.FormTemplateId))
			}
		}
		for _, customItem := range param.CustomItems {
			customItemExist = false
			if customItem.Id != "" {
				existMap[customItem.Id] = true
				for _, formItemTemplate := range formItemTemplateList {
					if customItem.Id == formItemTemplate.Id {
						customItemExist = true
						if customItem.FormTemplate == "" {
							customItem.FormTemplate = param.FormTemplateId
						}
						updateItems = append(updateItems, models.ConvertFormItemTemplateDto2Model(customItem))
						break
					}
				}
			}
			if !customItemExist {
				customItem.Id = guid.CreateGuid()
				customItem.FormTemplate = param.FormTemplateId
				if customItem.FormTemplate == "" {
					customItem.FormTemplate = param.FormTemplateId
				}
				insertItems = append(insertItems, models.ConvertFormItemTemplateDto2Model(customItem))
			}
		}
		for _, formItemTemplate := range formItemTemplateList {
			if existMap[formItemTemplate.Id] == false && existMap[formItemTemplate.AttrDefId] == false {
				deleteItems = append(deleteItems, formItemTemplate)
			}
		}
	}
	err = transaction(func(session *xorm.Session) error {
		if newItemGroupId != "" {
			if param.TaskTemplateId != "" {
				_, err = session.Exec("insert into form_template_new(id,request_template,task_template,item_group,item_group_name,item_group_type,item_group_rule,item_group_sort,"+
					"created_time) values(?,?,?,?,?,?,?,?,?)", newItemGroupId, param.RequestTemplateId, param.TaskTemplateId, param.ItemGroup, param.ItemGroupName, param.ItemGroupType,
					param.ItemGroupRule, s.CalcItemGroupSort(param.FormTemplateId), time.Now().Format(models.DateTimeFormat))
			} else {
				_, err = session.Exec("insert into form_template_new(id,request_template,item_group,item_group_name,item_group_type,item_group_rule,item_group_sort,"+
					"created_time) values(?,?,?,?,?,?,?,?)", newItemGroupId, param.RequestTemplateId, param.ItemGroup, param.ItemGroupName, param.ItemGroupType,
					param.ItemGroupRule, s.CalcItemGroupSort(param.FormTemplateId), time.Now().Format(models.DateTimeFormat))
			}
			if err != nil {
				return err
			}
		}
		if formItemTemplateGroup != nil {
			err = s.formItemTemplateGroupDao.Update(session, formItemTemplateGroup)
			if err != nil {
				return err
			}
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

func (s FormItemTemplateService) CalcItemGroupSort(formTemplateId string) int {
	var max = 0
	list, err := s.formItemTemplateGroupDao.QueryByFormTemplateId(formTemplateId)
	if err != nil {
		log.Logger.Error("CalcItemGroupSort err", log.Error(err))
		return 0
	}
	if len(list) > 0 {
		for _, group := range list {
			if group.ItemGroupSort >= max {
				max = group.ItemGroupSort
			}
		}
	}
	return max + 1
}

func (s FormItemTemplateService) UpdateFormTemplateItemGroup(param models.FormTemplateGroupCustomDataDto) (err error) {
	var formItemTemplateList []*models.FormItemTemplateTable
	var insertItems, updateItems, deleteItems []*models.FormItemTemplateTable
	var exist bool
	formItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplate(param.FormTemplateId)
	if err != nil {
		return
	}
	if len(formItemTemplateList) == 0 {
		formItemTemplateList = []*models.FormItemTemplateTable{}
	}

	if len(param.Items) == 0 {
		param.Items = []*models.FormItemTemplateDto{}
	}
	for _, item := range param.Items {
		// Id 为空新增
		if item.Id == "" {
			item.Id = guid.CreateGuid()
			item.FormTemplate = param.FormTemplateId
			item.FormTemplate = param.FormTemplateId
			insertItems = append(insertItems, models.ConvertFormItemTemplateDto2Model(item))
		} else {
			item.FormTemplate = param.FormTemplateId
			updateItems = append(updateItems, models.ConvertFormItemTemplateDto2Model(item))
		}
	}
	for _, formItemTemplate := range formItemTemplateList {
		exist = false
		for _, item := range param.Items {
			if item.Id == formItemTemplate.Id {
				exist = true
			}
		}
		if !exist {
			deleteItems = append(deleteItems, formItemTemplate)
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

func (s FormItemTemplateService) CopyDataFormTemplateItemGroup(requestTemplate, formTemplateId, taskTemplate string) (err error) {
	var formItemTemplateList []*models.FormItemTemplateTable
	var formItemTemplateGroup *models.FormTemplateNewTable
	var newItemGroupId string
	formItemTemplateGroup, err = s.formItemTemplateGroupDao.Get(formTemplateId)
	if err != nil {
		return
	}
	if formItemTemplateGroup == nil {
		err = fmt.Errorf("item-group-id is invalid")
		return
	}
	// 1. 查询表单组是否存在
	formItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplate(formTemplateId)
	if err != nil {
		return
	}
	// 新增数据
	err = transaction(func(session *xorm.Session) error {
		newItemGroupId = guid.CreateGuid()
		_, err = s.formItemTemplateGroupDao.Add(session, &models.FormTemplateNewTable{
			Id:              newItemGroupId,
			RequestTemplate: requestTemplate,
			TaskTemplate:    taskTemplate,
			ItemGroup:       formItemTemplateGroup.ItemGroup,
			ItemGroupType:   formItemTemplateGroup.ItemGroupType,
			ItemGroupName:   formItemTemplateGroup.ItemGroupName,
			ItemGroupSort:   s.CalcItemGroupSort(formTemplateId),
			ItemGroupRule:   formItemTemplateGroup.ItemGroupRule,
			RefId:           formTemplateId,
			CreatedTime:     time.Now().Format(models.DateTimeFormat),
		})
		if err != nil {
			return err
		}
		if len(formItemTemplateList) > 0 {
			for _, formItemTemplate := range formItemTemplateList {
				formItemTemplate.RefId = formItemTemplate.Id
				formItemTemplate.Id = guid.CreateGuid()
				formItemTemplate.FormTemplate = newItemGroupId
				_, err = s.formItemTemplateDao.Add(session, formItemTemplate)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	return
}
