package service

import (
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"strings"
	"time"
	"xorm.io/xorm"
)

type FormTemplateLibraryService struct {
	formTemplateLibraryDao     *dao.FormTemplateLibraryDao
	formItemTemplateLibraryDao *dao.FormItemTemplateLibraryDao
}

func (s *FormTemplateLibraryService) CheckNameExist(name string) (exist bool, err error) {
	var list []*models.FormTemplateLibraryTable
	list, err = s.formTemplateLibraryDao.QueryByName(name)
	exist = len(list) > 0
	return
}

func (s *FormTemplateLibraryService) Get(id string) (formTemplateLibraryTable *models.FormTemplateLibraryTable, err error) {
	formTemplateLibraryTable, err = s.formTemplateLibraryDao.Get(id)
	return
}

func (s *FormTemplateLibraryService) QueryAllFormType() (formTypes []string, err error) {
	var formTypeMap = make(map[string]bool)
	var list []*models.FormTemplateLibraryTable
	formTypes = []string{}
	if list, err = s.formTemplateLibraryDao.QueryAll(); err != nil {
		return
	}
	if len(list) > 0 {
		for _, formTemplateLibrary := range list {
			formTypeMap[formTemplateLibrary.Name] = true
		}
		for key, _ := range formTypeMap {
			formTypes = append(formTypes, key)
		}
	}
	return
}

func (s *FormTemplateLibraryService) AddFormTemplateLibrary(param models.FormTemplateLibraryParam, user string) (err error) {
	now := time.Now().Format(models.DateTimeFormat)
	formTemplateLibraryId := guid.CreateGuid()
	err = transaction(func(session *xorm.Session) error {
		formTemplateLibrary := &models.FormTemplateLibraryTable{
			Id:          formTemplateLibraryId,
			Name:        param.Name,
			FormType:    param.FormType,
			CreatedTime: now,
			UpdatedTime: now,
			CreatedBy:   user,
			UpdatedBy:   user,
		}
		// 添加表单组件库
		if _, err = s.formTemplateLibraryDao.Add(session, formTemplateLibrary); err != nil {
			return err
		}
		// 添加表单项组件库
		if len(param.Items) > 0 {
			for _, item := range param.Items {
				item.FormTemplateLibrary = formTemplateLibraryId
				if _, err = s.formItemTemplateLibraryDao.Add(session, item); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return
}

// DeleteFormTemplateLibrary  删除表单组件库,逻辑删除. 删除表单项组件库&表单项引用的外键需要删除
func (s *FormTemplateLibraryService) DeleteFormTemplateLibrary(id string) (err error) {
	var formItemTemplateLibraryList []*models.FormItemTemplateLibraryTable
	if formItemTemplateLibraryList, err = s.formItemTemplateLibraryDao.QueryByFormTemplateLibrary(id); err != nil {
		return
	}
	err = transaction(func(session *xorm.Session) error {
		if err = s.formTemplateLibraryDao.Disable(session, id); err != nil {
			return err
		}
		if len(formItemTemplateLibraryList) > 0 {
			for _, formItemTemplateLibrary := range formItemTemplateLibraryList {
				// 清空表单项引用外键
				if _, err = session.Exec("update form_item_template set form_item_library = null where form_item_library =? ", formItemTemplateLibrary.Id); err != nil {
					return err
				}
				// 删除表单项组件库
				if err = s.formItemTemplateLibraryDao.Delete(session, formItemTemplateLibrary.Id); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return
}

func (s *FormTemplateLibraryService) QueryFormTemplateLibrary(param models.QueryFormTemplateLibraryParam) (pageInfo models.PageInfo, list []*models.FormTemplateLibraryDto, err error) {
	var formTemplateLibraryList []*models.FormTemplateLibraryTable
	var formItemTemplateLibraryList []*models.FormItemTemplateLibraryTable
	var items []string
	list = []*models.FormTemplateLibraryDto{}
	if pageInfo, formTemplateLibraryList, err = s.formTemplateLibraryDao.QueryListByCondition(param); err != nil {
		return
	}
	if len(formTemplateLibraryList) > 0 {
		for _, formTemplateLibrary := range formTemplateLibraryList {
			items = []string{}
			formItemTemplateLibraryList = []*models.FormItemTemplateLibraryTable{}
			if formItemTemplateLibraryList, err = s.formItemTemplateLibraryDao.QueryByFormTemplateLibrary(formTemplateLibrary.Id); err != nil {
				return
			}
			if len(formItemTemplateLibraryList) > 0 {
				for _, item := range formItemTemplateLibraryList {
					items = append(items, item.Name)
				}
			}
			list = append(list, &models.FormTemplateLibraryDto{
				Id:          formTemplateLibrary.Id,
				Name:        formTemplateLibrary.Name,
				FormType:    formTemplateLibrary.FormType,
				CreatedTime: formTemplateLibrary.CreatedTime,
				CreatedBy:   formTemplateLibrary.CreatedBy,
				FormItems:   strings.Join(items, ","),
			})
		}
	}
	return
}
