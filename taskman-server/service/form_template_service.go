package service

import (
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/exterror"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/rpc"
	"sort"
	"strings"
	"time"
	"xorm.io/xorm"
)

type FormTemplateService struct {
	formTemplateDao     *dao.FormTemplateDao
	formItemTemplateDao *dao.FormItemTemplateDao
	formDao             *dao.FormDao
}

func (s *FormTemplateService) AddRequestFormTemplate(session *xorm.Session, formTemplateDto models.FormTemplateDto) (newId string, err error) {
	newId = guid.CreateGuid()
	itemIds := guid.CreateGuidList(len(formTemplateDto.Items))
	formTemplateDto.NowTime = time.Now().Format(models.DateTimeFormat)
	formTemplateDto.Id = newId
	// 新建 formTemplate
	if len(formTemplateDto.Items) > 0 && formTemplateDto.Items[0].FormTemplate == "" {
		_, err = session.Exec("insert into form_template(id,request_template,request_form_type,created_time) values(?,?,?,?)",
			newId, formTemplateDto.RequestTemplate, formTemplateDto.RequestFormType, time.Now().Format(models.DateTimeFormat))
		if err != nil {
			return
		}
	}
	// 添加模板项
	for i, item := range formTemplateDto.Items {
		item.Id = itemIds[i]
		item.FormTemplate = newId
		_, err = s.formItemTemplateDao.Add(session, models.ConvertFormItemTemplateDto2Model(item))
		if err != nil {
			return
		}
	}
	return
}

func (s *FormTemplateService) UpdateFormTemplateByDto(session *xorm.Session, formTemplateDto models.FormTemplateDto) (err error) {
	var formItemTemplateList []*models.FormItemTemplateDto
	formTemplateDto.NowTime = time.Now().Format(models.DateTimeFormat)
	newItemGuidList := guid.CreateGuidList(len(formTemplateDto.Items))
	formItemTemplateList, err = s.formItemTemplateDao.QueryDtoByFormTemplate(formTemplateDto.Id)
	if err != nil {
		return
	}
	// 新增or更新表单项模板
	for i, inputItem := range formTemplateDto.Items {
		if inputItem.Id == "" {
			inputItem.Id = newItemGuidList[i]
			if inputItem.FormTemplate == "" {
				inputItem.FormTemplate = formTemplateDto.Id
			}
			_, err = s.formItemTemplateDao.Add(session, models.ConvertFormItemTemplateDto2Model(inputItem))
			if err != nil {
				return
			}
		} else {
			if inputItem.FormTemplate == "" {
				inputItem.FormTemplate = formTemplateDto.Id
			}
			err = s.formItemTemplateDao.Update(session, models.ConvertFormItemTemplateDto2Model(inputItem))
			if err != nil {
				return
			}
		}
	}
	// 删除表单项&表单项模板
	for _, formItemTemplate := range formItemTemplateList {
		existFlag := false
		for _, inputItem := range formTemplateDto.Items {
			if formItemTemplate.Id == inputItem.Id {
				existFlag = true
				break
			}
		}
		if !existFlag {
			err = s.formItemTemplateDao.Delete(session, formItemTemplate.Id)
			if err != nil {
				return
			}
			// 审批表单、任务表单中的表单项都是copy数据表单而来,此处通过copy_id记录关系,当删除数据表单的表单项内容时候,对应任务表单、审批表单的表单项都需要删除
			err = s.formItemTemplateDao.DeleteByIdOrRefId(session, formItemTemplate.Id)
		}
	}
	return
}

func (s *FormTemplateService) GetRequestFormTemplate(requestTemplateId string) (result *models.FormTemplateDto, err error) {
	var requestTemplate *models.RequestTemplateTable
	var formTemplate *models.FormTemplateTable
	result = &models.FormTemplateDto{Items: []*models.FormItemTemplateDto{}}
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return
	}
	if requestTemplate == nil {
		err = fmt.Errorf("requestTemplate not exist")
		return
	}
	formTemplate, err = s.formTemplateDao.QueryRequestFormByRequestTemplateIdAndType(requestTemplateId, string(models.RequestFormTypeMessage))
	if err != nil {
		return
	}
	if formTemplate == nil {
		return
	}
	result.ExpireDay = requestTemplate.ExpireDay
	result.Id = formTemplate.Id
	result.Items, err = s.formItemTemplateDao.QueryDtoByFormTemplate(formTemplate.Id)
	return
}

func (s *FormTemplateService) QueryRequestFormByRequestTemplateIdAndType(requestTemplateId, requestFormType string) (formTemplate *models.FormTemplateTable, err error) {
	formTemplate, err = s.formTemplateDao.QueryRequestFormByRequestTemplateIdAndType(requestTemplateId, requestFormType)
	return
}

func (s *FormTemplateService) GetDataFormTemplate(requestTemplateId string) (result *models.DataFormTemplateDto, err error) {
	var requestTemplate *models.RequestTemplateTable
	var associationWorkflow bool
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return
	}
	if requestTemplate == nil {
		err = fmt.Errorf("requestTemplate not exist")
		return
	}
	// 关联编排
	if strings.TrimSpace(requestTemplate.ProcDefId) != "" {
		associationWorkflow = true
	}
	result = &models.DataFormTemplateDto{Groups: make([]*models.FormTemplateGroupDto, 0), AssociationWorkflow: associationWorkflow}
	result.Groups, err = s.getFormTemplateGroups(requestTemplateId, "", string(models.RequestFormTypeData))
	return
}

func (s *FormTemplateService) GetFormTemplate(requestTemplateId, taskTemplateId string) (result *models.SimpleFormTemplateDto, err error) {
	result = &models.SimpleFormTemplateDto{TaskTemplateId: taskTemplateId, Groups: make([]*models.FormTemplateGroupDto, 0)}
	result.Groups, err = s.getFormTemplateGroups(requestTemplateId, taskTemplateId, "")
	return
}

func (s *FormTemplateService) getFormTemplateGroups(requestTemplateId, taskTemplateId, requestFormType string) (groups []*models.FormTemplateGroupDto, err error) {
	var formItemTemplateGroupList []*models.FormTemplateTable
	var formItemTemplateList []*models.FormItemTemplateTable
	groups = []*models.FormTemplateGroupDto{}
	formItemTemplateGroupList, err = s.formTemplateDao.QueryListByRequestTemplateAndTaskTemplate(requestTemplateId, taskTemplateId, requestFormType)
	if err != nil {
		return
	}
	for _, group := range formItemTemplateGroupList {
		formItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplate(group.Id)
		groups = append(groups, &models.FormTemplateGroupDto{
			ItemGroupId:   group.Id,
			ItemGroup:     group.ItemGroup,
			ItemGroupType: group.ItemGroupType,
			ItemGroupName: group.ItemGroupName,
			ItemGroupSort: group.ItemGroupSort,
			Items:         models.ConvertFormItemTemplateModelList2Dto(formItemTemplateList, group),
		})
	}
	// 设置排序,保证前端展示数据顺序一致
	for _, FormTemplateGroupDto := range groups {
		if len(FormTemplateGroupDto.Items) > 0 {
			sort.Sort(models.FormItemTemplateDtoSort(FormTemplateGroupDto.Items))
		}
	}
	sort.Sort(models.FormTemplateGroupDtoSort(groups))
	return
}

func (s *FormTemplateService) GetDataFormTemplateItemGroups(requestTemplateId string) (result []*models.FormTemplateTable, err error) {
	var requestTemplate *models.RequestTemplateTable
	result = []*models.FormTemplateTable{}
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return
	}
	if requestTemplate == nil {
		err = fmt.Errorf("requestTemplate not exist")
		return
	}
	result, err = s.formTemplateDao.QueryListByRequestTemplateAndTaskTemplate(requestTemplateId, "", string(models.RequestFormTypeData))
	// 排序
	sort.Sort(models.FormTemplateTableSort(result))
	return
}

func (s *FormTemplateService) CreateRequestFormTemplate(formTemplateDto models.FormTemplateDto) (err error) {
	var requestTemplate *models.RequestTemplateTable
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(formTemplateDto.RequestTemplate)
	if err != nil {
		return err
	}
	if requestTemplate == nil {
		return exterror.Catch(exterror.New().RequestParamValidateError, fmt.Errorf("param id is invalid"))
	}
	err = transactionWithoutForeignCheck(func(session *xorm.Session) error {
		// 设置请求表单类型
		formTemplateDto.RequestFormType = string(models.RequestFormTypeMessage)
		formTemplateDto.Id, err = s.AddRequestFormTemplate(session, formTemplateDto)
		if err != nil {
			return err
		}
		// 更新模板
		err = GetRequestTemplateService().UpdateRequestTemplateUpdatedBy(session, formTemplateDto.RequestTemplate, formTemplateDto.UpdatedBy)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

func (s *FormTemplateService) UpdateRequestFormTemplate(formTemplateDto models.FormTemplateDto) (err error) {
	// 需要对当前用户进行校验&操作时间进行校验
	var requestTemplate *models.RequestTemplateTable
	var formTemplate *models.FormTemplateTable
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(formTemplateDto.RequestTemplate)
	if err != nil {
		return
	}
	if requestTemplate == nil {
		return exterror.Catch(exterror.New().RequestParamValidateError, fmt.Errorf("param id is invalid"))
	}
	formTemplate, err = s.formTemplateDao.Get(formTemplateDto.Id)
	if err != nil {
		return
	}
	if formTemplate == nil {
		return exterror.Catch(exterror.New().RequestParamValidateError, fmt.Errorf("param form_template_id is invalid"))
	}
	err = transactionWithoutForeignCheck(func(session *xorm.Session) error {
		// 更新表单项模板
		err = s.UpdateFormTemplateByDto(session, formTemplateDto)
		if err != nil {
			return err
		}
		// 更新模板
		err = GetRequestTemplateService().UpdateRequestTemplateUpdatedBy(session, formTemplateDto.RequestTemplate, formTemplateDto.UpdatedBy)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// GetFormConfig 获取配置表单,数据基于数据表单数据
func (s *FormTemplateService) GetFormConfig(requestTemplateId, taskTemplateId, formTemplateId, userToken, language string) (configureDto *models.FormTemplateGroupConfigureDto, err error) {
	var requestTemplate *models.RequestTemplateTable
	var dataFormConfigureDto *models.FormTemplateGroupConfigureDto
	var formItemTemplateList []*models.FormItemTemplateDto
	var existAttrMap = make(map[string]bool)
	var existCustomItemsMap = make(map[string]bool)
	var formItemTemplateGroup *models.FormTemplateTable
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return
	}
	if requestTemplate == nil {
		err = exterror.Catch(exterror.New().RequestParamValidateError, fmt.Errorf("param requestTemplateId is invalid"))
		return
	}
	configureDto = &models.FormTemplateGroupConfigureDto{RequestTemplateId: requestTemplateId, TaskTemplateId: taskTemplateId, FormTemplateId: formTemplateId, SystemItems: []*models.ProcEntityAttributeObj{}, CustomItems: []*models.FormItemTemplateDto{}}
	// 1.先查询用户配置数据
	if formTemplateId != "" {
		// 1.先查询用户配置数据
		formItemTemplateList, err = s.formItemTemplateDao.QueryDtoByFormTemplate(formTemplateId)
		if err != nil {
			return
		}
		// 查询表单组
		formItemTemplateGroup, err = s.formTemplateDao.Get(formTemplateId)
		if err != nil {
			return
		}
		if formItemTemplateGroup != nil {
			configureDto.ItemGroup = formItemTemplateGroup.ItemGroup
			configureDto.ItemGroupName = formItemTemplateGroup.ItemGroupName
			configureDto.ItemGroupType = formItemTemplateGroup.ItemGroupType
			configureDto.ItemGroupRule = formItemTemplateGroup.ItemGroupRule
			configureDto.ItemGroupSort = formItemTemplateGroup.ItemGroupSort
		}
		for _, formItemTemplate := range formItemTemplateList {
			if formItemTemplate.AttrDefId != "" {
				existAttrMap[formItemTemplate.AttrDefId] = true
			} else {
				existCustomItemsMap[formItemTemplate.RefId] = true
			}
		}
		// 2. 查询数据表单
		if formItemTemplateGroup != nil {
			dataFormConfigureDto, err = s.GetDataFormConfig(requestTemplateId, taskTemplateId, formItemTemplateGroup.RefId, "", "", userToken, language)
			if err != nil {
				return
			}
			// 将数据表单的选中作为 审批、任务表单的全量
			if len(dataFormConfigureDto.SystemItems) > 0 {
				for _, systemItem := range dataFormConfigureDto.SystemItems {
					if systemItem.Active {
						configureDto.SystemItems = append(configureDto.SystemItems, systemItem)
						if !existAttrMap[systemItem.Id] {
							systemItem.Active = false
						}
					}
				}
			}
			if len(dataFormConfigureDto.CustomItems) > 0 {
				configureDto.CustomItems = dataFormConfigureDto.CustomItems
				for _, customItem := range dataFormConfigureDto.CustomItems {
					customItem.Active = false
					if existCustomItemsMap[customItem.Id] {
						customItem.Active = true
					}
				}
			}
		}
	}
	return
}

// GetDataFormConfig 获取数据表单配置
func (s *FormTemplateService) GetDataFormConfig(requestTemplateId, taskTemplateId, formTemplateId, formType, entity, userToken, language string) (configureDto *models.FormTemplateGroupConfigureDto, err error) {
	var formItemTemplateList []*models.FormItemTemplateDto
	var entitiesList []*models.ExpressionEntities
	var expressEntity *models.ExpressionEntities
	var existAttrMap = make(map[string]bool)
	var formItemTemplateGroup *models.FormTemplateTable
	configureDto = &models.FormTemplateGroupConfigureDto{RequestTemplateId: requestTemplateId, TaskTemplateId: taskTemplateId, FormTemplateId: formTemplateId, SystemItems: []*models.ProcEntityAttributeObj{}, CustomItems: []*models.FormItemTemplateDto{}}
	// 1.先查询用户配置数据
	if formTemplateId != "" {
		formItemTemplateList, err = s.formItemTemplateDao.QueryDtoByFormTemplate(formTemplateId)
		if err != nil {
			return
		}
		// 查询表单组
		formItemTemplateGroup, err = s.formTemplateDao.Get(formTemplateId)
		if err != nil {
			return
		}
		if formItemTemplateGroup != nil {
			configureDto.ItemGroup = formItemTemplateGroup.ItemGroup
			configureDto.ItemGroupName = formItemTemplateGroup.ItemGroupName
			configureDto.ItemGroupType = formItemTemplateGroup.ItemGroupType
			configureDto.ItemGroupRule = formItemTemplateGroup.ItemGroupRule
			configureDto.ItemGroupSort = formItemTemplateGroup.ItemGroupSort
		}
		if len(formItemTemplateList) > 0 {
			for _, formItem := range formItemTemplateList {
				if formItem.ElementType == string(models.FormItemElementTypeCalculate) {
					formItem.Active = true
					configureDto.CustomItems = append(configureDto.CustomItems, formItem)
				} else {
					existAttrMap[formItem.AttrDefId] = true
				}
			}
			// 自定义表单组的 entity 为空
			if configureDto.ItemGroupType != string(models.FormItemGroupTypeCustom) {
				entity = configureDto.ItemGroupName
			}
		}
	}
	// 2.查询entity 属性集合
	if entity != "" {
		entitiesList, err = rpc.QueryEntityAttributes(models.QueryExpressionDataParam{DataModelExpression: entity}, userToken, language)
		if err != nil {
			return
		}
		if len(entitiesList) > 0 && len(entitiesList[0].Attributes) > 0 {
			expressEntity = entitiesList[0]
			if configureDto.ItemGroup == "" {
				configureDto.ItemGroup = entity
				configureDto.ItemGroupName = entity
				configureDto.ItemGroupType = formType
			}
			for _, attribute := range entitiesList[0].Attributes {
				attribute.Id = fmt.Sprintf("%s:%s:%s", entitiesList[0].PackageName, entitiesList[0].EntityName, attribute.Name)
				attribute.EntityName = expressEntity.EntityName
				attribute.EntityPackage = expressEntity.PackageName
				if existAttrMap[attribute.Id] {
					attribute.Active = true
				}
				configureDto.SystemItems = append(configureDto.SystemItems, attribute)
			}
		}
	}
	return
}

func (s *FormTemplateService) CleanDataForm(requestTemplateId string) (err error) {
	var list []*models.FormTemplateTable
	var formItemTemplateList []*models.FormItemTemplateTable
	list, err = s.formTemplateDao.QueryListByRequestTemplateAndTaskTemplate(requestTemplateId, "", string(models.RequestFormTypeData))
	if err != nil {
		return
	}
	if len(list) > 0 {
		for _, formTemplate := range list {
			formItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplate(formTemplate.Id)
			if err != nil {
				return err
			}
			if len(formItemTemplateList) == 0 {
				// 需要删除该表单
				err = s.formTemplateDao.Delete(nil, formTemplate.Id)
				if err != nil {
					return err
				}
			}
		}
	}
	return
}

func (s *FormTemplateService) DeleteFormTemplateItemGroup(formTemplateId string) (err error) {
	var formTemplateList []*models.FormTemplateTable
	formTemplateList, err = s.formTemplateDao.QueryListByIdOrRefId(formTemplateId)
	if err != nil {
		return err
	}
	err = transaction(func(session *xorm.Session) error {
		if len(formTemplateList) > 0 {
			for _, formTemplate := range formTemplateList {
				err = s.formItemTemplateDao.DeleteByFormTemplate(session, formTemplate.Id)
				if err != nil {
					return err
				}
				err = s.formTemplateDao.Delete(session, formTemplate.Id)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	return
}

func (s *FormTemplateService) DeleteFormTemplateItemGroupTransaction(session *xorm.Session, formTemplateId string) (err error) {
	var formTemplateList []*models.FormTemplateTable
	formTemplateList, err = s.formTemplateDao.QueryListByIdOrRefId(formTemplateId)
	if err != nil {
		return
	}
	if len(formTemplateList) > 0 {
		for _, formTemplate := range formTemplateList {
			err = s.formItemTemplateDao.DeleteByFormTemplate(session, formTemplate.Id)
			if err != nil {
				return
			}
			err = s.formTemplateDao.Delete(session, formTemplate.Id)
			if err != nil {
				return
			}
		}
	}
	return
}
func (s *FormTemplateService) SortFormTemplateItemGroup(param models.FormTemplateGroupSortDto) (err error) {
	var formTemplateList []*models.FormTemplateTable
	var formTemplateSortMap = s.buildFormTemplateGroupSortMap(param.ItemGroupIdSort)
	formTemplateList, err = s.formTemplateDao.QueryListByRequestTemplateAndTaskTemplate(param.RequestTemplateId, param.TaskTemplateId, "")
	if err != nil {
		return
	}
	if len(formTemplateList) > 0 {
		err = transaction(func(session *xorm.Session) error {
			for _, formTemplate := range formTemplateList {
				formTemplate.ItemGroupSort = formTemplateSortMap[formTemplate.Id]
				err = s.formTemplateDao.Update(session, formTemplate)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	return
}

func (s *FormTemplateService) buildFormTemplateGroupSortMap(itemGroupIdSort []string) map[string]int {
	hashMap := make(map[string]int)
	for i, groupId := range itemGroupIdSort {
		hashMap[groupId] = i
	}
	return hashMap
}

func (s *FormTemplateService) DeleteWorkflowFormTemplateGroupSql(requestTemplateId string) (actions []*dao.ExecAction, err error) {
	var formTemplateList []*models.FormTemplateTable
	actions = []*dao.ExecAction{}
	err = dao.X.SQL("select * from form_template where request_template=?  and item_group_type='workflow'", requestTemplateId).Find(&formTemplateList)
	if err != nil {
		return
	}
	if len(formTemplateList) > 0 {
		for _, formTemplate := range formTemplateList {
			// 过滤掉 任务类型
			if strings.HasPrefix(formTemplate.TaskTemplate, "im_") {
				continue
			}
			actions = append(actions, &dao.ExecAction{Sql: "delete from form_item_template where form_template=?", Param: []interface{}{formTemplate.Id}})
			actions = append(actions, &dao.ExecAction{Sql: "delete from form_template where id=?", Param: []interface{}{formTemplate.Id}})
		}
	}
	return
}
