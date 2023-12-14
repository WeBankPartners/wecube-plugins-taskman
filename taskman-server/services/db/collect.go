package db

import (
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"strings"
	"time"
)

func AddTemplateCollect(param *models.CollectTemplateTable) error {
	param.Id = guid.CreateGuid()
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := x.Exec("insert into collect_template(id,request_template,type,user,role,created_time) value (?,?,?,?,?,?)",
		param.Id, param.RequestTemplate, param.Type, param.User, param.Role, nowTime)
	if err != nil {
		err = fmt.Errorf("Insert database error:%s ", err.Error())
	}
	return err
}

func DeleteTemplateCollect(templateId, user string) error {
	_, err := x.Exec("delete from collect_template where request_template = ? and user = ?", templateId, user)
	if err != nil {
		err = fmt.Errorf("Delete database error:%s ", err.Error())
	}
	return err
}

// QueryTemplateCollect 查询模板收藏
func QueryTemplateCollect(param *models.QueryCollectTemplateParam, user, userToken string) (pageInfo models.PageInfo, rowData []*models.CollectDataObj, err error) {
	var result models.ProcNodeObjList
	var collectTemplateList []*models.CollectTemplateTable
	var disableTemplateVersionMap = getAllDisableTemplateVersionMap()
	var roleTemplateMap = make(map[string]string)
	var templateUserRoleMap = make(map[string]bool)
	var resultList []string
	var templateType int
	if param.Action == 1 {
		// 发布
		templateType = 1
	} else if param.Action == 2 {
		// 请求
		templateType = 0
	}
	// 查询该用户收藏的所有模板id
	err = x.SQL("select * from collect_template where user = ? and type = ?", user, templateType).Find(&collectTemplateList)
	if err != nil {
		return
	}
	// 遍历模板id,查询 当前最新的发布模板id
	for _, collectTemplate := range collectTemplateList {
		var tempList []string
		x.SQL("select id from request_template where status='confirm' and parent_id = ? order by created_time desc limit 0,1", collectTemplate.RequestTemplate).Find(&tempList)
		resultList = append(resultList, tempList...)
		roleTemplateMap[collectTemplate.RequestTemplate] = collectTemplate.Role
	}
	if len(resultList) == 0 {
		return
	}
	sql := fmt.Sprintf("select * from (select rt.id,rt.parent_id,rt.name,rtg.id as template_group_id,rtg.name as template_group,rtg.manage_role as template_group_role,rt.operator_obj_type,"+
		"rt.proc_def_name,rt.handler as owner,rt.tags,rt.created_time,rt.updated_time from request_template rt "+
		"join request_template_group rtg on rt.group= rtg.id where rt.id in ("+getSQL(resultList)+")) t %s", transCollectConditionToSQL(param))
	// 排序处理
	if param.Sorting != nil {
		hashMap, _ := getJsonToXormMap(models.CollectDataObj{})
		if len(hashMap) > 0 {
			if param.Sorting.Asc {
				sql += fmt.Sprintf(" ORDER BY %s ASC", hashMap[param.Sorting.Field])
			} else {
				sql += fmt.Sprintf(" ORDER BY %s DESC", hashMap[param.Sorting.Field])
			}
		}
	}
	// 分页处理
	pageInfo.PageSize = param.PageSize
	pageInfo.StartIndex = param.StartIndex
	pageInfo.TotalRows = queryCount(sql)
	err = x.SQL(sql+" limit ?,?", param.StartIndex, param.PageSize).Find(&rowData)
	if err != nil {
		return
	}
	if len(rowData) > 0 {
		for _, collectObj := range rowData {
			templateUserRoleMap = make(map[string]bool, 0)
			template, err := GetSimpleRequestTemplate(collectObj.Id)
			if err != nil {
				continue
			}
			if template.Status != "confirm" {
				collectObj.Version = "beta"
			} else {
				collectObj.Version = template.Version
			}
			requestTemplateRoleList, _ := getRequestTemplateRole(collectObj.Id)
			if len(requestTemplateRoleList) == 0 {
				continue
			}
			for _, requestTemplateRole := range requestTemplateRoleList {
				if requestTemplateRole.RoleType == "MGMT" {
					collectObj.ManageRole = requestTemplateRole.Role
				} else if requestTemplateRole.RoleType == "USE" {
					templateUserRoleMap[requestTemplateRole.Role] = true
				}
			}
			collectObj.UseRole = roleTemplateMap[collectObj.ParentId]
			collectObj.Status = 1
			// 判断 收藏模板是否被禁用. 禁用版本 大于等于当前模板版本表示禁用
			if disableTemplateVersionMap[template.ParentId] != "" && strings.Compare(disableTemplateVersionMap[template.ParentId], template.Version) >= 0 {
				collectObj.Status = 2
			} else if !templateUserRoleMap[collectObj.UseRole] {
				// 模板使用权限变更,导致收藏模板时候角色,没权限新建请求
				collectObj.Status = 3
			}
			result, err = GetProcessNodesByProc(models.RequestTemplateTable{Id: collectObj.Id}, userToken, "template")
			if err != nil {
				continue
			}
			for _, item := range result {
				collectObj.WorkNode = append(collectObj.WorkNode, item.NodeName)
			}
		}
	}
	return
}

func getAllDisableTemplateVersionMap() map[string]string {
	var hashMap = make(map[string]string)
	var list []*models.RequestTemplateTable
	x.SQL("select * from request_template where status = 'disable'").Find(&list)
	if len(list) > 0 {
		for _, requestTemplate := range list {
			hashMap[requestTemplate.ParentId] = requestTemplate.Version
		}
	}
	return hashMap
}

func QueryAllTemplateCollect(user string) (collectMap map[string]bool, err error) {
	collectMap = make(map[string]bool)
	var idList []string
	err = x.SQL("select request_template from collect_template where user = ?", user).Find(&idList)
	if err != nil {
		return
	}
	for _, id := range idList {
		collectMap[id] = true
	}
	return
}

func GetCollectFilterItem(param *models.FilterRequestParam, user string) (data *models.CollectFilterItem, err error) {
	data = &models.CollectFilterItem{}
	var pairList []models.KeyValuePair
	var rowsData []*models.CollectDataObj
	var templateGroupMap = make(map[string]string)
	var operatorObjTypeMap = make(map[string]bool)
	var procDefNameMap = make(map[string]bool)
	var ownerMap = make(map[string]bool)
	var tagMap = make(map[string]bool)
	var manageRoleMap = make(map[string]bool)
	var useRoleMap = make(map[string]bool)
	var sql = "select rt.id,rt.name,rtg.id as template_group_id,rtg.name  as template_group ,rt.operator_obj_type,rt.proc_def_name,rtg.manage_role,rt.handler as owner,rt.tags,rt.created_time from request_template rt " +
		"join request_template_group rtg on rt.group= rtg.id where rt.id in (select request_template from collect_template where user = ?) and rt.created_time > ?"
	err = x.SQL(sql, user, param.StartTime).Find(&rowsData)
	if err != nil {
		return
	}
	if len(rowsData) > 0 {
		for _, row := range rowsData {
			templateGroupMap[row.TemplateGroup] = row.TemplateGroupId
			operatorObjTypeMap[row.OperatorObjType] = true
			procDefNameMap[row.ProcDefName] = true
			ownerMap[row.Owner] = true
			tagMap[row.Tags] = true
			manageRoleMap[row.ManageRole] = true
			var roleList []string
			err = x.SQL("select role from request_template_role where role_type='USE' and request_template= ?", row.Id).Find(&roleList)
			if err != nil || len(roleList) == 0 {
				continue
			}
			for _, role := range roleList {
				useRoleMap[role] = true
			}
		}
	}
	for key, value := range templateGroupMap {
		pairList = append(pairList, models.KeyValuePair{TemplateId: value, TemplateName: key})
	}
	data.TemplateGroupList = pairList
	data.OperatorObjTypeList = convertMap2Array(operatorObjTypeMap)
	data.ProcDefNameList = convertMap2Array(procDefNameMap)
	data.OwnerList = convertMap2Array(ownerMap)
	data.TagList = convertMap2Array(tagMap)
	data.ManageRoleList = convertMap2Array(manageRoleMap)
	data.UseRoleList = convertMap2Array(useRoleMap)
	return
}

// CheckUserCollectExist 检查用户模板id是否收藏
func CheckUserCollectExist(templateId, user string) bool {
	var idList []string
	err := x.SQL("select id from collect_template where request_template=? and user=?", templateId, user).Find(&idList)
	if err != nil {
		return true
	}
	if len(idList) > 0 {
		return true
	}
	return false
}
