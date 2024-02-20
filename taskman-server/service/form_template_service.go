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
	var formItemTemplateList []*models.FormItemTemplateTable
	var itemGroupMap = make(map[string][]*models.FormItemTemplateTable)
	var itemGroupType, itemGroupName string
	var itemGroupSort int
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
	result = &models.DataFormTemplateDto{DataFormTemplateId: requestTemplate.DataFormTemplate, Groups: make([]*models.DataFormTemplateGroupDto, 0), AssociationWorkflow: associationWorkflow}
	formItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplate(requestTemplate.DataFormTemplate)
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
				formItemTemplateArr = append(formItemTemplateArr, formItemTemplate)
			}
		}
		if len(formItemTemplateArr) > 0 {
			itemGroupType = formItemTemplateArr[0].ItemGroupType
			itemGroupName = formItemTemplateArr[0].ItemGroupName
			itemGroupSort = formItemTemplateArr[0].ItemGroupSort
		}
		result.Groups = append(result.Groups, &models.DataFormTemplateGroupDto{
			ItemGroup:     itemGroup,
			ItemGroupType: itemGroupType,
			ItemGroupName: itemGroupName,
			ItemGroupSort: itemGroupSort,
			Items:         formItemTemplateArr,
		})
	}
	// 设置排序,保证前端展示数据顺序一致
	for _, DataFormTemplateGroupDto := range result.Groups {
		if len(DataFormTemplateGroupDto.Items) > 0 {
			sort.Sort(models.FormItemTemplateTableSort(DataFormTemplateGroupDto.Items))
		}
	}
	sort.Sort(models.DataFormTemplateGroupDtoSort(result.Groups))
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
	// 请求模板的处理不是当前用户,不允许操作
	if requestTemplate.Handler != formTemplateDto.UpdatedBy {
		return exterror.New().DataPermissionDeny
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

// UpdateDataFormTemplate 更新数据表单
func (s FormTemplateService) UpdateDataFormTemplate(dataFormTemplateDto models.DataFormTemplateDto, requestTemplateId string) (err error) {
	var requestTemplate *models.RequestTemplateTable
	var formTemplateDto = models.FormTemplateDto{Items: []*models.FormItemTemplateTable{}}
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return err
	}
	if requestTemplate == nil {
		return exterror.Catch(exterror.New().RequestParamValidateError, fmt.Errorf("param id is invalid"))
	}
	// 请求模板的处理不是当前用户,不允许操作
	if requestTemplate.Handler != dataFormTemplateDto.UpdatedBy {
		return exterror.New().DataPermissionDeny
	}
	// 将数据表单各分组中的表单项,整合一起更新表单模板
	formTemplateDto.Id = dataFormTemplateDto.DataFormTemplateId
	if len(dataFormTemplateDto.Groups) > 0 {
		for _, groupDto := range dataFormTemplateDto.Groups {
			if groupDto != nil && len(groupDto.Items) > 0 {
				formTemplateDto.Items = append(formTemplateDto.Items, groupDto.Items...)
			}
		}
	}
	err = transactionWithoutForeignCheck(func(session *xorm.Session) error {
		// 更新表单项模板
		err = s.UpdateFormTemplate(session, formTemplateDto)
		if err != nil {
			return err
		}
		err = GetRequestTemplateService().UpdateRequestTemplateUpdatedBy(session, requestTemplateId, formTemplateDto.UpdatedBy)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// GetConfigureForm 获取配置表单
func (s FormTemplateService) GetConfigureForm(formTemplateId, itemGroupName, userToken, language string) (configureDto *models.FormTemplateGroupConfigureDto, err error) {
	var formItemTemplate []*models.FormItemTemplateTable
	var entitiesList []*models.ExpressionEntities
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
				existAttrMap[formItem.Name] = true
			}
		}
	}
	// 2.查询entity 属性集合
	entitiesList, err = rpc.QueryEntityAttributes(models.QueryExpressionDataParam{DataModelExpression: itemGroupName}, userToken, language)
	if err != nil {
		return
	}
	if len(entitiesList) > 0 && len(entitiesList[0].Attributes) > 0 {
		if configureDto.ItemGroup == "" {
			configureDto.ItemGroup = itemGroupName
			configureDto.ItemGroupName = itemGroupName
			configureDto.ItemGroupType = string(models.FormItemGroupTypeOptional)
		}
		for _, attribute := range entitiesList[0].Attributes {
			if existAttrMap[attribute.Name] {
				attribute.Active = true
			}
			configureDto.SystemItems = append(configureDto.SystemItems, attribute)
		}
	}
	return
}
