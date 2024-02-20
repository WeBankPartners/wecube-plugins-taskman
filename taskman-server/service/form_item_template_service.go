package service

import (
	"fmt"
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
	// 1. 查询表单组是否存在，不存在则新增
	formItemTemplateList, err = s.formItemTemplateDao.QueryByFormTemplateAndItemGroupName(param.FormTemplateId, param.ItemGroupName)
	if err != nil {
		return
	}
	// 新增数据
	if len(formItemTemplateList) == 0 {

		return
	}
	// 更新数据
	err = transaction(func(session *xorm.Session) error {
		for _, formItemTemplate := range formItemTemplateList {
			fmt.Printf("%v", formItemTemplate)
		}
		return nil
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

func (s FormItemTemplateService) buildFormTemplateGroupSortMap(itemGroupNameSort []string) map[string]int {
	hashMap := make(map[string]int)
	for i, groupName := range itemGroupNameSort {
		hashMap[groupName] = i
	}
	return hashMap
}
