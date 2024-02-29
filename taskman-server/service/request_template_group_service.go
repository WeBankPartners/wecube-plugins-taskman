package service

import (
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"time"
)

type RequestTemplateGroupService struct {
	requestTemplateGroupDao *dao.RequestTemplateGroupDao
}

func (s *RequestTemplateGroupService) QueryRequestTemplateGroup(param *models.QueryRequestParam, userRoles []string) (pageInfo models.PageInfo, rowData []*models.RequestTemplateGroupTable, err error) {
	rowData = []*models.RequestTemplateGroupTable{}
	filterSql, queryColumn, queryParam := dao.TransFiltersToSQL(param, &models.TransFiltersParam{IsStruct: true, StructObj: models.RequestTemplateGroupTable{}, PrimaryKey: "id"})
	userRoleFilterSql, userRoleFilterParams := dao.CreateListParams(userRoles, "")
	baseSql := fmt.Sprintf("SELECT %s FROM request_template_group WHERE manage_role in ("+userRoleFilterSql+") and del_flag=0 %s ", queryColumn, filterSql)
	queryParam = append(userRoleFilterParams, queryParam...)
	if param.Paging {
		pageInfo.StartIndex = param.Pageable.StartIndex
		pageInfo.PageSize = param.Pageable.PageSize
		pageInfo.TotalRows = dao.QueryCount(baseSql, queryParam...)
		pageSql, pageParam := dao.TransPageInfoToSQL(*param.Pageable)
		baseSql += pageSql
		queryParam = append(queryParam, pageParam...)
	}
	err = dao.X.SQL(baseSql, queryParam...).Find(&rowData)
	if len(rowData) > 0 {
		roleMap, _ := getRoleMap()
		for _, row := range rowData {
			if row.ManageRole != "" {
				row.ManageRoleObj = models.RoleTable{Id: row.ManageRole, DisplayName: roleMap[row.ManageRole].DisplayName}
			}
		}
	}
	return
}

func (s *RequestTemplateGroupService) CreateRequestTemplateGroup(param *models.RequestTemplateGroupTable) error {
	param.Id = guid.CreateGuid()
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := dao.X.Exec("insert into request_template_group(id,name,description,manage_role,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?,?)",
		param.Id, param.Name, param.Description, param.ManageRole, param.CreatedBy, nowTime, param.CreatedBy, nowTime)
	if err != nil {
		err = fmt.Errorf("Insert database error:%s ", err.Error())
	}
	return err
}

func (s *RequestTemplateGroupService) UpdateRequestTemplateGroup(param *models.RequestTemplateGroupTable) error {
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := dao.X.Exec("update request_template_group set name=?,description =?,manage_role=?,updated_by=?,updated_time=? where id=?",
		param.Name, param.Description, param.ManageRole, param.UpdatedBy, nowTime, param.Id)
	if err != nil {
		err = fmt.Errorf("Update database error:%s ", err.Error())
	}
	return err
}

func (s *RequestTemplateGroupService) DeleteRequestTemplateGroup(id string) error {
	_, err := dao.X.Exec("update request_template_group set del_flag=1 where id=?", id)
	if err != nil {
		err = fmt.Errorf("Delete database error:%s ", err.Error())
	}
	return err
}

func (s *RequestTemplateGroupService) CheckRequestTemplateGroupRoles(id string, roles []string) error {
	rolesFilterSql, rolesFilterParam := dao.CreateListParams(roles, "")
	var requestTemplateGroupRows []*models.RequestTemplateGroupTable
	rolesFilterParam = append([]interface{}{id}, rolesFilterParam...)
	err := dao.X.SQL("select id from request_template_group where id=? and manage_role in ("+rolesFilterSql+")", rolesFilterParam...).Find(&requestTemplateGroupRows)
	if err != nil {
		return fmt.Errorf("Try to query database data fail,%s ", err.Error())
	}
	if len(requestTemplateGroupRows) == 0 {
		return fmt.Errorf(models.RowDataPermissionErr)
	}
	return nil
}
