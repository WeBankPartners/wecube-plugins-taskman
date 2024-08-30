package service

import (
	"context"
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/exterror"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/gin-gonic/gin"
	"sort"
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
			formTypeMap[formTemplateLibrary.FormType] = true
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
		}
		// 添加表单组件库
		if _, err = s.formTemplateLibraryDao.Add(session, formTemplateLibrary); err != nil {
			return err
		}
		// 添加表单项组件库
		if len(param.Items) > 0 {
			ids := guid.CreateGuidList(len(param.Items))
			for i, item := range param.Items {
				item.Id = ids[i]
				item.FormTemplateLibrary = formTemplateLibraryId
				if _, err = s.formItemTemplateLibraryDao.Add(session, models.ConvertFormItemTemplateLibraryDto2Model(item)); err != nil {
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
			sort.Sort(models.FormItemTemplateLibraryTableSort(formItemTemplateLibraryList))
			if len(formItemTemplateLibraryList) > 0 {
				for _, item := range formItemTemplateLibraryList {
					items = append(items, item.Title)
				}
			}
			list = append(list, &models.FormTemplateLibraryDto{
				Id:          formTemplateLibrary.Id,
				Name:        formTemplateLibrary.Name,
				FormType:    formTemplateLibrary.FormType,
				CreatedTime: formTemplateLibrary.CreatedTime,
				CreatedBy:   formTemplateLibrary.CreatedBy,
				FormItems:   strings.Join(items, "、"),
				Items:       models.ConvertFormItemTemplateLibraryModel2Dto(formItemTemplateLibraryList),
			})
		}
	}
	return
}

func ExportFormTemplateLibrary(ctx context.Context) (result []*models.FormTemplateLibraryTableData, err error) {
	var formTemplateLibraryData []*models.FormTemplateLibraryTableData
	baseSql := "SELECT * FROM form_template_library WHERE del_flag=0"
	err = dao.X.SQL(baseSql).Find(&formTemplateLibraryData)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}

	if len(formTemplateLibraryData) == 0 {
		err = fmt.Errorf("query form template library empty")
		return
	}

	result = formTemplateLibraryData

	// query formItemTemplateLibrary
	var formTemplateLibraryIds []string
	for _, formTempalteLib := range formTemplateLibraryData {
		formTemplateLibraryIds = append(formTemplateLibraryIds, formTempalteLib.Id)
	}

	var formItemTemplateLibraryData []*models.FormItemTemplateLibraryTableData
	filterSql, filterParam := dao.CreateListParams(formTemplateLibraryIds, "")
	baseSql = fmt.Sprintf("SELECT * FROM form_item_template_library WHERE form_template_library IN (%s)", filterSql)
	err = dao.X.SQL(baseSql, filterParam...).Find(&formItemTemplateLibraryData)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}

	formTemplateLibraryIdMapItemInfo := make(map[string][]*models.FormItemTemplateLibraryTableData)
	for _, itemInfo := range formItemTemplateLibraryData {
		formTemplateLibraryIdMapItemInfo[itemInfo.FormTemplateLibrary] = append(
			formTemplateLibraryIdMapItemInfo[itemInfo.FormTemplateLibrary], itemInfo)
	}

	for _, formTemplateLib := range formTemplateLibraryData {
		formTemplateLib.Items = []*models.FormItemTemplateLibraryTableData{}
		if _, isExisted := formTemplateLibraryIdMapItemInfo[formTemplateLib.Id]; isExisted {
			formTemplateLib.Items = formTemplateLibraryIdMapItemInfo[formTemplateLib.Id]
		}
	}
	return
}

func ImportFormTemplateLibrary(c *gin.Context, formTemplateLibraryData []*models.FormTemplateLibraryTableData) (err error) {
	var actions []*dao.ExecAction
	now := time.Now().Format(models.DateTimeFormat)
	user := middleware.GetRequestUser(c)

	if len(formTemplateLibraryData) == 0 {
		return
	}

	tableNameFormTemplateLibrary := "form_template_library"
	tableNameFormItemTemplateLibrary := "form_item_template_library"

	for _, formTemplateLibInfo := range formTemplateLibraryData {
		queryFormTemplateLibData := &models.FormTemplateLibraryTableData{}
		var exists bool
		querySql := fmt.Sprintf("SELECT * FROM form_template_library WHERE id = ?")
		queryParam := []interface{}{formTemplateLibInfo.Id}
		exists, err = dao.X.SQL(querySql, queryParam...).Get(queryFormTemplateLibData)
		if err != nil {
			err = exterror.Catch(exterror.New().DatabaseQueryError, err)
			return
		}

		if !exists {
			formTemplateLibData := &models.FormTemplateLibraryTableData{
				Id:          formTemplateLibInfo.Id,
				Name:        formTemplateLibInfo.Name,
				FormType:    formTemplateLibInfo.FormType,
				CreatedTime: now,
				UpdatedTime: now,
				CreatedBy:   user,
			}
			action, tmpErr := dao.GetInsertTableExecAction(tableNameFormTemplateLibrary, *formTemplateLibData, nil)
			if tmpErr != nil {
				err = fmt.Errorf("get insert sql for formTemplateLibData.Id: %s failed: %s", formTemplateLibData.Id, tmpErr.Error())
				return
			}
			actions = append(actions, action)
		} else {
			// update
			updateColumnStr := "`name`=?,`form_type`=?,`updated_time`=?,`created_by`=?,`del_flag`=?,`custom_flag`=?"
			action := &dao.ExecAction{
				Sql: dao.CombineDBSql("UPDATE ", tableNameFormTemplateLibrary, " SET ", updateColumnStr, " WHERE id=?"),
				Param: []interface{}{formTemplateLibInfo.Name, formTemplateLibInfo.FormType, now, user,
					formTemplateLibInfo.DelFlag, formTemplateLibInfo.CustomFlag, formTemplateLibInfo.Id},
			}
			actions = append(actions, action)
		}

		formTemplateLibId := formTemplateLibInfo.Id
		// update formItemTemplateLibrary
		// firstly delete original formItemTemplateLibrary and then create new formItemTemplateLibrary
		action := &dao.ExecAction{
			Sql:   dao.CombineDBSql("DELETE FROM ", tableNameFormItemTemplateLibrary, " WHERE form_template_library=?"),
			Param: []interface{}{formTemplateLibId},
		}
		actions = append(actions, action)

		var formItemTemplateLibList []*models.FormItemTemplateLibraryTableData
		formItemTemplateLibList = formTemplateLibInfo.Items
		for i := range formItemTemplateLibList {
			formItemTemplateLibList[i].FormTemplateLibrary = formTemplateLibId
			action, tmpErr := dao.GetInsertTableExecAction(tableNameFormItemTemplateLibrary, *formItemTemplateLibList[i], nil)
			if tmpErr != nil {
				err = fmt.Errorf("get insert sql of formItemTemplateLibrary for formTemplateLibData.Id: %s failed: %s", formTemplateLibInfo.Id, tmpErr.Error())
				return
			}
			actions = append(actions, action)
		}
	}

	err = dao.Transaction(actions)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseExecuteError, err)
		return
	}
	return
}
