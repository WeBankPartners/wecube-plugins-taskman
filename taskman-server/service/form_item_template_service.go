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
	formItemTemplateDao *dao.FormItemTemplateDao
	formTemplateDao     *dao.FormTemplateDao
}

func (s *FormItemTemplateService) UpdateFormTemplateItemGroupConfig(param models.FormTemplateGroupConfigureDto, userToken string) (err error) {
	var formTemplate *models.FormTemplateTable
	var insertItems, updateItems, deleteItems []*models.FormItemTemplateTable
	var formItemTemplateList, refFormItemTemplateList []*models.FormItemTemplateTable
	var newItemGroupId string
	var systemItemExist, customItemExist bool
	var existMap = make(map[string]bool)
	var sort int
	if len(param.SystemItems) == 0 {
		param.SystemItems = []*models.ProcEntityAttributeObj{}
	}
	if len(param.CustomItems) == 0 {
		param.CustomItems = []*models.FormItemTemplateDto{}
	}
	// 1. 查询表单模板是否存在，不存在则新增
	formTemplate, err = s.formTemplateDao.Get(param.FormTemplateId)
	if err != nil {
		return
	}
	// 新增数据
	if formTemplate == nil {
		newItemGroupId = guid.CreateGuid()
		if len(param.SystemItems) > 0 {
			for _, systemItem := range param.SystemItems {
				refAttributes, tmpErr := getCMDBCiAttrDefs(systemItem.EntityName, userToken)
				if tmpErr != nil {
					err = fmt.Errorf("query remote entity:%s attr fail:%s ", systemItem.EntityName, tmpErr.Error())
					return
				}
				insertItems = append(insertItems, models.ConvertProcEntityAttributeObj2FormItemTemplate(param, systemItem, newItemGroupId, refAttributes))
				sort = systemItem.OrderNo + 2
			}
		}
		if len(param.CustomItems) > 0 {
			for _, customItem := range param.CustomItems {
				customItem.RefId = customItem.Id
				customItem.Id = guid.CreateGuid()
				customItem.FormTemplate = newItemGroupId
				customItem.ItemGroup = param.ItemGroup
				customItem.ItemGroupName = param.ItemGroupName
				customItem.ElementType = string(models.FormItemElementTypeCalculate)
				customItem.Sort = sort
				insertItems = append(insertItems, models.ConvertFormItemTemplateDto2Model(customItem))
				sort++
			}
		}
	} else {
		param.TaskTemplateId = formTemplate.TaskTemplate
		// 直接更新表单组
		formTemplate.ItemGroupName = param.ItemGroupName
		formTemplate.ItemGroup = param.ItemGroup
		formTemplate.ItemGroupRule = param.ItemGroupRule
		formTemplate.ItemGroupType = param.ItemGroupType

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
				refAttributes, tmpErr := getCMDBCiAttrDefs(systemItem.EntityName, userToken)
				if tmpErr != nil {
					err = fmt.Errorf("query remote entity:%s attr fail:%s ", systemItem.EntityName, tmpErr.Error())
					return
				}
				insertItems = append(insertItems, models.ConvertProcEntityAttributeObj2FormItemTemplate(param, systemItem, param.FormTemplateId, refAttributes))
			}
		}
		// 如果有 taskTemplateId, customItems展示的都是数据表单的 customItems数据,需要做一层转换
		if param.TaskTemplateId != "" {
			for _, customItem := range param.CustomItems {
				exist := false
				for _, formItem := range formItemTemplateList {
					if customItem.Id == formItem.RefId {
						customItem = models.ConvertFormItemTemplateModel2Dto(formItem, *formTemplate)
						exist = true
					}
				}
				if !exist {
					// 将 refId 指向 Id
					customItem.RefId = customItem.Id
				}
			}
		}
		for _, customItem := range param.CustomItems {
			customItemExist = false
			if customItem.Id != "" {
				refFormItemTemplateList = []*models.FormItemTemplateTable{}
				existMap[customItem.Id] = true
				for _, formItemTemplate := range formItemTemplateList {
					if customItem.Id == formItemTemplate.Id {
						customItemExist = true
						if customItem.FormTemplate == "" {
							customItem.FormTemplate = param.FormTemplateId
						}
						// 查询 refId指向当前ID的数据
						refFormItemTemplateList, _ = s.formItemTemplateDao.QueryByRefId(customItem.Id)
						if len(refFormItemTemplateList) > 0 {
							for _, refFormItemTemplate := range refFormItemTemplateList {
								// 主要更新 routineExpression 值
								refFormItemTemplate.RoutineExpression = customItem.RoutineExpression
								updateItems = append(updateItems, refFormItemTemplate)
							}
						}
						updateItems = append(updateItems, models.ConvertFormItemTemplateDto2Model(customItem))
						break
					}
				}
			}
			if !customItemExist {
				customItem.RefId = customItem.Id
				customItem.Id = guid.CreateGuid()
				customItem.FormTemplate = param.FormTemplateId
				customItem.ItemGroup = formTemplate.ItemGroup
				customItem.ItemGroupName = formTemplate.ItemGroupName
				insertItems = append(insertItems, models.ConvertFormItemTemplateDto2Model(customItem))
			}
		}
		for _, formItemTemplate := range formItemTemplateList {
			// 过滤掉自定义类型
			if formItemTemplate.AttrDefId == "" && formItemTemplate.ElementType != "calculate" {
				continue
			}
			// 这块只判断,编排类型和 自定义计算类型
			if existMap[formItemTemplate.Id] == false && existMap[formItemTemplate.AttrDefId] == false {
				deleteItems = append(deleteItems, formItemTemplate)
			}
		}
	}
	err = transaction(func(session *xorm.Session) error {
		if newItemGroupId != "" {
			if param.TaskTemplateId != "" {
				_, err = session.Exec("insert into form_template(id,request_template,task_template,item_group,item_group_name,item_group_type,item_group_rule,item_group_sort,"+
					"created_time,request_form_type) values(?,?,?,?,?,?,?,?,?,?)", newItemGroupId, param.RequestTemplateId, param.TaskTemplateId, param.ItemGroup, param.ItemGroupName, param.ItemGroupType,
					param.ItemGroupRule, s.CalcItemGroupSort(param.RequestTemplateId, param.TaskTemplateId), time.Now().Format(models.DateTimeFormat), models.RequestFormTypeData)
			} else {
				_, err = session.Exec("insert into form_template(id,request_template,item_group,item_group_name,item_group_type,item_group_rule,item_group_sort,"+
					"created_time,request_form_type) values(?,?,?,?,?,?,?,?,?)", newItemGroupId, param.RequestTemplateId, param.ItemGroup, param.ItemGroupName, param.ItemGroupType,
					param.ItemGroupRule, s.CalcItemGroupSort(param.RequestTemplateId, param.TaskTemplateId), time.Now().Format(models.DateTimeFormat), models.RequestFormTypeData)
			}
			if err != nil {
				return err
			}
		}
		if formTemplate != nil {
			err = s.formTemplateDao.Update(session, formTemplate)
			if err != nil {
				return err
			}
		}
		if len(insertItems) > 0 {
			newGuidList := guid.CreateGuidList(len(insertItems))
			for i, item := range insertItems {
				item.Id = "item_tpl_" + newGuidList[i]
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
				err = s.formItemTemplateDao.DeleteByIdOrRefId(session, item.Id)
				if err != nil {
					return err
				}
			}
		}
		return err
	})
	return
}

func (s *FormItemTemplateService) CalcItemGroupSort(requestTemplateId, taskTemplateId string) int {
	var max = 0
	// 如果任务模板ID,则只查询数据模板
	list, err := s.formTemplateDao.QueryListByRequestTemplateAndTaskTemplate(requestTemplateId, taskTemplateId, string(models.RequestFormTypeData))
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

func (s *FormItemTemplateService) UpdateFormTemplateItemGroup(param models.FormTemplateGroupCustomDataDto) (err error) {
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

func (s *FormItemTemplateService) CopyDataFormTemplateItemGroup(requestTemplateId, formTemplateId, taskTemplateId string) (err error) {
	var formItemTemplateList []*models.FormItemTemplateTable
	var formTemplate *models.FormTemplateTable
	var newItemGroupId string
	formTemplate, err = s.formTemplateDao.Get(formTemplateId)
	if err != nil {
		return
	}
	if formTemplate == nil {
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
		_, err = s.formTemplateDao.Add(session, &models.FormTemplateTable{
			Id:              newItemGroupId,
			RequestTemplate: requestTemplateId,
			TaskTemplate:    taskTemplateId,
			ItemGroup:       formTemplate.ItemGroup,
			ItemGroupType:   formTemplate.ItemGroupType,
			ItemGroupName:   formTemplate.ItemGroupName,
			ItemGroupSort:   s.CalcItemGroupSort(requestTemplateId, taskTemplateId),
			ItemGroupRule:   formTemplate.ItemGroupRule,
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

func (s *FormItemTemplateService) GetFormItemTemplate(formItemTemplateId string) (result *models.FormItemTemplateTable, err error) {
	result, err = s.formItemTemplateDao.Get(formItemTemplateId)
	return
}
