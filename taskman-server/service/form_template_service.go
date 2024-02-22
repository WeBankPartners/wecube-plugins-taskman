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
	formTemplateDao     dao.FormTemplateDao
	formItemTemplateDao dao.FormItemTemplateDao
	formDao             dao.FormDao
}

func (s FormTemplateService) AddFormTemplate(session *xorm.Session, formTemplateDto models.FormTemplateDto) (newId string, err error) {
	newId = guid.CreateGuid()
	itemIds := guid.CreateGuidList(len(formTemplateDto.Items))
	formTemplateDto.NowTime = time.Now().Format(models.DateTimeFormat)
	formTemplateDto.Id = newId
	// 添加模板
	_, err = s.formTemplateDao.Add(session, models.CovertFormTemplateDto2Model(formTemplateDto))
	if err != nil {
		return
	}
	// 添加模板项
	for i, item := range formTemplateDto.Items {
		item.Id = itemIds[i]
		item.FormTemplate = newId
		_, err = s.formItemTemplateDao.Add(session, item)
		if err != nil {
			return
		}
	}
	return
}

func (s FormTemplateService) UpdateFormTemplate(session *xorm.Session, formTemplateDto models.FormTemplateDto) (err error) {
	var formItemTemplateList []*models.FormItemTemplateTable
	formTemplateDto.NowTime = time.Now().Format(models.DateTimeFormat)
	newItemGuidList := guid.CreateGuidList(len(formTemplateDto.Items))
	formTemplateTable := &models.FormTemplateTable{
		Id:          formTemplateDto.Id,
		Name:        formTemplateDto.Name,
		Description: formTemplateDto.Description,
		UpdatedBy:   formTemplateDto.UpdatedBy,
		UpdatedTime: formTemplateDto.UpdatedTime,
	}
	// 更新表单模板
	err = s.formTemplateDao.Update(session, formTemplateTable)
	if err != nil {
		return
	}
	formItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplate(formTemplateDto.Id)
	if err != nil {
		return
	}
	// 新增or更新表单项模板
	for i, inputItem := range formTemplateDto.Items {
		if inputItem.Id == "" {
			inputItem.Id = newItemGuidList[i]
			_, err = s.formItemTemplateDao.Add(session, inputItem)
			if err != nil {
				return
			}
		} else {
			err = s.formItemTemplateDao.Update(session, inputItem)
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
			err = s.formItemTemplateDao.DeleteByIdOrCopyId(session, formItemTemplate.Id)
		}
	}
	return
}

func (s FormTemplateService) GetRequestFormTemplate(requestTemplateId string) (result *models.FormTemplateDto, err error) {
	var requestTemplate *models.RequestTemplateTable
	var formTemplate *models.FormTemplateTable
	result = &models.FormTemplateDto{Items: []*models.FormItemTemplateTable{}}
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return
	}
	if requestTemplate == nil {
		err = fmt.Errorf("requestTemplate not exist")
		return
	}
	formTemplate, err = s.formTemplateDao.Get(requestTemplate.FormTemplate)
	if err != nil {
		return
	}
	if formTemplate == nil {
		return
	}
	result.ExpireDay = requestTemplate.ExpireDay
	result.Id = formTemplate.Id
	result.Name = formTemplate.Name
	result.Description = formTemplate.Description
	result.UpdatedTime = formTemplate.UpdatedTime
	result.UpdatedBy = formTemplate.UpdatedBy
	result.Items, err = s.formItemTemplateDao.QueryByFormTemplate(requestTemplate.FormTemplate)
	return
}

func (s FormTemplateService) GetDataFormTemplate(requestTemplateId string) (result *models.DataFormTemplateDto, err error) {
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
	// 新增数据表单
	if requestTemplate.DataFormTemplate == "" {
		err = s.CreateDataFormTemplate(models.DataFormTemplateDto{}, requestTemplateId)
		if err != nil {
			return
		}
		requestTemplate, _ = GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	}
	result = &models.DataFormTemplateDto{DataFormTemplateId: requestTemplate.DataFormTemplate, Groups: make([]*models.FormTemplateGroupDto, 0), AssociationWorkflow: associationWorkflow}
	result.Groups, err = s.getFormTemplateGroups(requestTemplate.DataFormTemplate)
	return
}

func (s FormTemplateService) GetFormTemplate(formTemplateId string) (result *models.SimpleFormTemplateDto, err error) {
	result = &models.SimpleFormTemplateDto{FormTemplateId: formTemplateId, Groups: make([]*models.FormTemplateGroupDto, 0)}
	result.Groups, err = s.getFormTemplateGroups(formTemplateId)
	return
}

func (s FormTemplateService) getFormTemplateGroups(formTemplateId string) (groups []*models.FormTemplateGroupDto, err error) {
	var formItemTemplateList []*models.FormItemTemplateTable
	var itemGroupMap = make(map[string][]*models.FormItemTemplateTable)
	var itemGroupType, itemGroupName string
	var itemGroupSort int
	groups = []*models.FormTemplateGroupDto{}
	formItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplate(formTemplateId)
	if err != nil {
		return
	}
	if len(formItemTemplateList) == 0 {
		return
	}
	for _, formItemTemplate := range formItemTemplateList {
		if _, ok := itemGroupMap[formItemTemplate.ItemGroup]; !ok {
			itemGroupMap[formItemTemplate.ItemGroup] = make([]*models.FormItemTemplateTable, 0)
		}
	}
	for itemGroup, formItemTemplateArr := range itemGroupMap {
		for _, formItemTemplate := range formItemTemplateList {
			if itemGroup == formItemTemplate.ItemGroup {
				if formItemTemplate.ItemGroupType == string(models.FormItemGroupTypeCustom) && formItemTemplate.Name == defaultCustomFormItemName {
					continue
				}
				formItemTemplateArr = append(formItemTemplateArr, formItemTemplate)
			}
		}
		if len(formItemTemplateArr) > 0 {
			itemGroupType = formItemTemplateArr[0].ItemGroupType
			itemGroupName = formItemTemplateArr[0].ItemGroupName
			itemGroupSort = formItemTemplateArr[0].ItemGroupSort
		}
		groups = append(groups, &models.FormTemplateGroupDto{
			ItemGroup:     itemGroup,
			ItemGroupType: itemGroupType,
			ItemGroupName: itemGroupName,
			ItemGroupSort: itemGroupSort,
			Items:         formItemTemplateArr,
		})
	}
	// 设置排序,保证前端展示数据顺序一致
	for _, FormTemplateGroupDto := range groups {
		if len(FormTemplateGroupDto.Items) > 0 {
			sort.Sort(models.FormItemTemplateTableSort(FormTemplateGroupDto.Items))
		}
	}
	sort.Sort(models.FormTemplateGroupDtoSort(groups))
	return
}

func (s FormTemplateService) GetDataFormTemplateItemGroups(requestTemplateId string) (entityList []string, err error) {
	var itemGroupNameMap = make(map[string]bool)
	var requestTemplate *models.RequestTemplateTable
	var formItemTemplateList []*models.FormItemTemplateTable
	entityList = []string{}
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return
	}
	if requestTemplate == nil {
		err = fmt.Errorf("requestTemplate not exist")
		return
	}
	formItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplate(requestTemplate.DataFormTemplate)
	if err != nil {
		return
	}
	for _, formItemTemplate := range formItemTemplateList {
		itemGroupNameMap[formItemTemplate.ItemGroupName] = true
	}
	for groupName, _ := range itemGroupNameMap {
		entityList = append(entityList, groupName)
	}
	// 排序
	sort.Strings(entityList)
	return
}

func (s FormTemplateService) CreateRequestFormTemplate(formTemplateDto models.FormTemplateDto, requestTemplateId string) (err error) {
	var requestTemplate *models.RequestTemplateTable
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return err
	}
	if requestTemplate == nil {
		return exterror.Catch(exterror.New().RequestParamValidateError, fmt.Errorf("param id is invalid"))
	}
	// 请求模板的处理不是当前用户,不允许操作
	if requestTemplate.Handler != formTemplateDto.UpdatedBy {
		return exterror.New().DataPermissionDeny
	}
	err = transactionWithoutForeignCheck(func(session *xorm.Session) error {
		// 添加表单模板
		formTemplateDto.Id, err = s.AddFormTemplate(session, formTemplateDto)
		if err != nil {
			return err
		}
		// 更新模板
		err = GetRequestTemplateService().UpdateRequestTemplateBase(session, requestTemplateId, formTemplateDto.Id, formTemplateDto.Description, formTemplateDto.UpdatedBy, formTemplateDto.ExpireDay)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

func (s FormTemplateService) UpdateRequestFormTemplate(formTemplateDto models.FormTemplateDto, requestTemplateId string) (err error) {
	// 需要对当前用户进行校验&操作时间进行校验
	var requestTemplate *models.RequestTemplateTable
	var formTemplate *models.FormTemplateTable
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return
	}
	if requestTemplate == nil {
		return exterror.Catch(exterror.New().RequestParamValidateError, fmt.Errorf("param id is invalid"))
	}
	// 请求模板的处理不是当前用户,不允许操作
	if requestTemplate.Handler != formTemplateDto.UpdatedBy {
		return exterror.New().DataPermissionDeny
	}
	formTemplate, err = s.formTemplateDao.Get(formTemplateDto.Id)
	if err != nil {
		return
	}
	if formTemplate == nil {
		return exterror.Catch(exterror.New().RequestParamValidateError, fmt.Errorf("param form_template_id is invalid"))
	}
	// 前端传递表单模板更新时间必须和数据库一致才能更新
	if formTemplate.UpdatedTime != formTemplateDto.UpdatedTime {
		return exterror.New().DealWithAtTheSameTimeError
	}
	err = transactionWithoutForeignCheck(func(session *xorm.Session) error {
		// 更新表单项模板
		err = s.UpdateFormTemplate(session, formTemplateDto)
		if err != nil {
			return err
		}
		// 更新模板
		err = GetRequestTemplateService().UpdateRequestTemplateBase(session, requestTemplateId, formTemplateDto.Id, formTemplateDto.Description, formTemplate.UpdatedBy, formTemplateDto.ExpireDay)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// CreateDataFormTemplate 创建数据表单
func (s FormTemplateService) CreateDataFormTemplate(formTemplateDto models.DataFormTemplateDto, requestTemplateId string) (err error) {
	var requestTemplate *models.RequestTemplateTable
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return err
	}
	if requestTemplate == nil {
		return exterror.Catch(exterror.New().RequestParamValidateError, fmt.Errorf("param id is invalid"))
	}
	err = transactionWithoutForeignCheck(func(session *xorm.Session) error {
		// 添加表单模板
		formTemplateDto.DataFormTemplateId, err = s.AddFormTemplate(session, models.ConvertDataFormTemplate2FormTemplateDto(formTemplateDto))
		if err != nil {
			return err
		}
		// 更新模板
		err = GetRequestTemplateService().UpdateRequestTemplateDataForm(session, requestTemplateId, formTemplateDto.DataFormTemplateId, formTemplateDto.UpdatedBy)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// GetFormConfig 获取配置表单,数据基于数据表单数据
func (s FormTemplateService) GetFormConfig(requestTemplateId, formTemplateId, itemGroupName, userToken, language string) (configureDto *models.FormTemplateGroupConfigureDto, err error) {
	var requestTemplate *models.RequestTemplateTable
	var dataFormConfigureDto *models.FormTemplateGroupConfigureDto
	var formItemTemplateList []*models.FormItemTemplateTable
	var existAttrMap = make(map[string]bool)
	var existCustomItemsMap = make(map[string]string)
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return
	}
	if requestTemplate == nil {
		err = exterror.Catch(exterror.New().RequestParamValidateError, fmt.Errorf("param requestTemplateId is invalid"))
		return
	}
	if requestTemplate.DataFormTemplate == "" {
		err = fmt.Errorf("requestTemplate:%s DataFormTemplate is empty", requestTemplate.Id)
		return
	}
	configureDto = &models.FormTemplateGroupConfigureDto{FormTemplateId: formTemplateId, SystemItems: []*models.ProcEntityAttributeObj{}, CustomItems: []*models.FormItemTemplateTable{}}
	// 1.先查询用户配置数据
	formItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplateAndItemGroupName(formTemplateId, itemGroupName)
	if err != nil {
		return
	}
	for _, formItemTemplate := range formItemTemplateList {
		if formItemTemplate.AttrDefId != "" {
			existAttrMap[formItemTemplate.AttrDefId] = true
		} else {
			configureDto.CustomItems = append(configureDto.CustomItems, formItemTemplate)
			existCustomItemsMap[formItemTemplate.CopyId] = formItemTemplate.Id
		}
	}
	// 2. 查询数据表单
	dataFormConfigureDto, err = s.GetDataFormConfig(requestTemplate.DataFormTemplate, itemGroupName, userToken, language)
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
		for _, customItem := range dataFormConfigureDto.CustomItems {
			if existCustomItemsMap[customItem.Id] != "" {
				customItem.Id = existCustomItemsMap[customItem.Id]
			} else {
				customItem.Id = ""
			}
			configureDto.CustomItems = append(configureDto.CustomItems, customItem)
		}
	}
	return
}

// GetDataFormConfig 获取数据表单配置
func (s FormTemplateService) GetDataFormConfig(formTemplateId, itemGroupName, userToken, language string) (configureDto *models.FormTemplateGroupConfigureDto, err error) {
	var formItemTemplate []*models.FormItemTemplateTable
	var entitiesList []*models.ExpressionEntities
	var entity *models.ExpressionEntities
	var existAttrMap = make(map[string]bool)
	configureDto = &models.FormTemplateGroupConfigureDto{FormTemplateId: formTemplateId, SystemItems: []*models.ProcEntityAttributeObj{}, CustomItems: []*models.FormItemTemplateTable{}}
	// 1.先查询用户配置数据
	formItemTemplate, err = s.formItemTemplateDao.QueryByFormTemplateAndItemGroupName(formTemplateId, itemGroupName)
	if err != nil {
		return
	}
	if len(formItemTemplate) > 0 {
		configureDto.ItemGroup = formItemTemplate[0].ItemGroup
		configureDto.ItemGroupName = formItemTemplate[0].ItemGroupName
		configureDto.ItemGroupType = formItemTemplate[0].ItemGroupType
		configureDto.ItemGroupRule = formItemTemplate[0].ItemGroupRule
		for _, formItem := range formItemTemplate {
			if formItem.ItemGroupType == string(models.FormItemGroupTypeCustom) {
				configureDto.CustomItems = append(configureDto.CustomItems, formItem)
			} else {
				existAttrMap[formItem.AttrDefId] = true
			}
		}
	}
	// 2.查询entity 属性集合
	entitiesList, err = rpc.QueryEntityAttributes(models.QueryExpressionDataParam{DataModelExpression: itemGroupName}, userToken, language)
	if err != nil {
		return
	}
	if len(entitiesList) > 0 && len(entitiesList[0].Attributes) > 0 {
		entity = entitiesList[0]
		if configureDto.ItemGroup == "" {
			configureDto.ItemGroup = itemGroupName
			configureDto.ItemGroupName = itemGroupName
			configureDto.ItemGroupType = string(models.FormItemGroupTypeOptional)
		}
		for _, attribute := range entitiesList[0].Attributes {
			attribute.Id = fmt.Sprintf("%s:%s:%s", entitiesList[0].PackageName, entitiesList[0].EntityName, attribute.Name)
			attribute.EntityName = entity.EntityName
			attribute.EntityPackage = entity.PackageName
			if existAttrMap[attribute.Id] {
				attribute.Active = true
			}
			configureDto.SystemItems = append(configureDto.SystemItems, attribute)
		}
	}
	return
}
