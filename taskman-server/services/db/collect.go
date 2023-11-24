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
	_, err := x.Exec("insert into collect_template(id,request_template,user,created_time) value (?,?,?,?)",
		param.Id, param.RequestTemplate, param.User, nowTime)
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
func QueryTemplateCollect(param *models.QueryCollectTemplateObj, user, userToken string) (pageInfo models.PageInfo, rowData []*models.CollectDataObj, err error) {
	var result models.ProcNodeObjList
	sql := "select rt.id,rt.name,rtg.name  as template_group ,rt.proc_def_name,rtg.manage_role,rt.handler as owner,rt.tags,rt.created_time from request_template rt " +
		"join request_template_group rtg on rt.group= rtg.id where rt.id in (select request_template from collect_template where user = ?) order by rt.created_time desc"
	pageInfo.PageSize = param.PageSize
	pageInfo.StartIndex = param.StartIndex
	pageInfo.TotalRows = queryCount(sql, user)
	err = x.SQL(sql+" limit ?,?", user, param.StartIndex, param.PageSize).Find(&rowData)
	if err != nil {
		return
	}
	if len(rowData) > 0 {
		for _, collectObj := range rowData {
			template, err := GetSimpleRequestTemplate(collectObj.Id)
			if err != nil {
				continue
			}
			if template.Status != "confirm" {
				collectObj.Name = fmt.Sprintf("%s(beta)", template.Name)
			} else {
				collectObj.Name = fmt.Sprintf("%s(%s)", template.Name, template.Version)
			}
			var roleList []string
			err = x.SQL("select role from request_template_role where role_type='USE' and request_template= ?", collectObj.Id).Find(&roleList)
			if err != nil || len(roleList) == 0 {
				continue
			}
			collectObj.UseRole = strings.Join(roleList, ",")
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
