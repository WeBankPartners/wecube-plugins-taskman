package service

import (
	"strings"

	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/rpc"
)

type RoleService struct {
}

func (s RoleService) GetRoleMail(roleList []*models.RoleTable, userToken, language string) (mailList []string) {
	var roleMap map[string]*models.SimpleLocalRoleDto
	var err error
	if models.CoreToken.BaseUrl == "" || len(roleList) == 0 {
		return
	}
	roleMap, err = rpc.QueryAllRoles("N", userToken, language)
	if err != nil {
		log.Logger.Error("QueryAllRoles err:%+v", log.Error(err))
		return
	}
	for _, v := range roleList {
		tmpMail := ""
		for _, vv := range roleMap {
			if vv.Name == v.Id {
				tmpMail = vv.Email
				break
			}
		}
		if tmpMail != "" {
			mailList = append(mailList, tmpMail)
		}
	}
	return
}

func (s RoleService) SyncCoreRole(userToken, language string) {
	var roleMap map[string]*models.SimpleLocalRoleDto
	var err error
	if models.CoreToken.BaseUrl == "" {
		return
	}
	roleMap, err = rpc.QueryAllRoles("N", userToken, language)
	var roleTable, addRoleList, delRoleList []*models.RoleTable
	err = dao.X.SQL("select * from role").Find(&roleTable)
	if err != nil {
		log.Logger.Error("Try to sync core role fail", log.Error(err))
		return
	}
	for _, v := range roleMap {
		existFlag := false
		for _, vv := range roleTable {
			if v.Name == vv.Id {
				existFlag = true
				break
			}
		}
		if !existFlag {
			addRoleList = append(addRoleList, &models.RoleTable{Id: v.Name, DisplayName: v.DisplayName, CoreId: v.ID, Email: v.Email})
		}
	}
	for _, v := range roleTable {
		existFlag := false
		for _, vv := range roleMap {
			if v.Id == vv.Name {
				existFlag = true
				break
			}
		}
		if !existFlag {
			delRoleList = append(delRoleList, &models.RoleTable{Id: v.Id})
		}
	}
	var actions []*dao.ExecAction
	for _, role := range addRoleList {
		actions = append(actions, &dao.ExecAction{Sql: "insert into `role`(id,display_name,core_id,email,updated_time) value (?,?,?,?,NOW())", Param: []interface{}{role.Id, role.DisplayName, role.CoreId, role.Email}})
	}
	if len(delRoleList) > 0 {
		roleIdList := []string{}
		for _, role := range delRoleList {
			actions = append(actions, &dao.ExecAction{Sql: "delete from `role` where id=?", Param: []interface{}{role.Id}})
			roleIdList = append(roleIdList, role.Id)
		}
		actions = append(actions, &dao.ExecAction{Sql: "update request_template_group set manage_role=NULL where manage_role in ('" + strings.Join(roleIdList, "','") + "')"})
	}
	if len(actions) > 0 {
		err = dao.TransactionWithoutForeignCheck(actions)
		if err != nil {
			log.Logger.Error("Sync core role fail", log.Error(err))
		}
	}
}

func (s RoleService) GetRoleList(ids []string) (result []*models.RoleTable, err error) {
	result = []*models.RoleTable{}
	if len(ids) == 0 {
		err = dao.X.SQL("select * from role").Find(&result)
	} else {
		idFilterSql, idFilterParam := dao.CreateListParams(ids, "")
		err = dao.X.SQL("select * from role where id in ("+idFilterSql+")", idFilterParam...).Find(&result)
	}
	return
}

func (s RoleService) QueryUserByRoles(roles []string, userToken, language string) (result []string, err error) {
	result = []string{}
	var roleTable []*models.RoleTable
	rolesFilterSql, rolesFilterParam := dao.CreateListParams(roles, "")
	err = dao.X.SQL("select * from `role` where id in ("+rolesFilterSql+")", rolesFilterParam...).Find(&roleTable)
	if err != nil || len(roleTable) == 0 {
		return
	}
	existMap := make(map[string]int)
	for _, v := range roleTable {
		if v.CoreId == "" {
			continue
		}
		tmpResult, tmpErr := rpc.QueryRolesUsers(v.CoreId, userToken, language)
		if tmpErr != nil {
			err = tmpErr
			break
		}
		if len(tmpResult) == 0 {
			continue
		}
		for _, vv := range tmpResult {
			if _, b := existMap[vv.UserName]; !b {
				result = append(result, vv.UserName)
				existMap[vv.UserName] = 1
			}
		}
	}
	return
}
