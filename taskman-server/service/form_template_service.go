package service

import (
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/exterror"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
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
			s.formItemTemplateDao.Update(session, inputItem)
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
			err = s.formDao.DeleteByFormItemTemplate(session, formItemTemplate.FormTemplate)
			if err != nil {
				return
			}
			err = s.formItemTemplateDao.Delete(session, formItemTemplate.Id)
		}
	}
	return
}

func (s FormTemplateService) GetRequestFormTemplate(id string) (result models.FormTemplateDto, err error) {
	var requestTemplate models.RequestTemplateTable
	var formTemplate *models.FormTemplateTable
	result = models.FormTemplateDto{Items: []*models.FormItemTemplateTable{}}
	requestTemplate, err = GetRequestTemplateService().GetSimpleRequestTemplate(id)
	if err != nil {
		return
	}
	formTemplate, err = s.formTemplateDao.Get(requestTemplate.GetFormTemplate())
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
	result.Items, err = s.formItemTemplateDao.QueryByFormTemplate(requestTemplate.GetFormTemplate())
	return
}

func (s FormTemplateService) CreateRequestFormTemplate(formTemplateDto models.FormTemplateDto, requestTemplateId string) (err error) {
	var requestTemplate *models.RequestTemplateTable
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return err
	}
	// 请求模板的处理不是当前用户,不允许操作
	if requestTemplate.Handler != formTemplateDto.UpdatedBy {
		return exterror.New().DataPermissionDeny
	}
	err = transactionWithoutForeignCheck(func(session *xorm.Session) error {
		// 添加表单模板
		formTemplateDto.Id, err = s.AddFormTemplate(session, formTemplateDto)
		// 更新模板
		err = GetRequestTemplateService().UpdateFormTemplate(session, requestTemplateId, formTemplateDto.Id, formTemplateDto.Description, formTemplateDto.ExpireDay)
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
		return err
	}
	// 请求模板的处理不是当前用户,不允许操作
	if requestTemplate.Handler != formTemplateDto.UpdatedBy {
		return exterror.New().DataPermissionDeny
	}
	formTemplate, err = s.formTemplateDao.Get(formTemplateDto.Id)
	if err != nil {
		return err
	}
	if formTemplate == nil {
		return exterror.Catch(exterror.New().RequestParamValidateError, fmt.Errorf("param id is invalid"))
	}
	// 前端传递表单模板更新时间必须和数据库一致才能更新
	if formTemplate.UpdatedTime != formTemplateDto.UpdatedTime {
		return exterror.New().DealWithAtTheSameTimeError
	}
	err = transaction(func(session *xorm.Session) error {
		// 更新表单项模板
		err = s.UpdateFormTemplate(session, formTemplateDto)
		// 更新模板
		err = GetRequestTemplateService().UpdateFormTemplate(session, requestTemplateId, formTemplateDto.Id, formTemplateDto.Description, formTemplateDto.ExpireDay)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

func (s FormTemplateService) DeleteRequestFormTemplate(id string) error {
	_, err := dao.X.Exec("update form_template set del_flag=1 where id=?", id)
	return err
}
