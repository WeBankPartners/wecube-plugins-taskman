package service

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/rpc"
)

type RoleService struct {
}

func (s *RoleService) GetRoleMap(userToken, language string) (roleMap map[string]*models.RoleTable, err error) {
	var roleDtoMap = make(map[string]*models.SimpleLocalRoleDto)
	if roleDtoMap, err = rpc.QueryAllRoles("N", userToken, language); err != nil {
		return
	}
	roleMap = make(map[string]*models.RoleTable)
	if len(roleDtoMap) > 0 {
		for _, roleDto := range roleDtoMap {
			roleMap[roleDto.Name] = &models.RoleTable{
				Id:          roleDto.Name,
				DisplayName: roleDto.DisplayName,
				CoreId:      roleDto.ID,
				Email:       roleDto.Email,
			}
		}
	}
	return
}

func (s *RoleService) GetRoleMail(roleList []*models.RoleTable, userToken, language string) (mailList []string) {
	var roleMap map[string]*models.SimpleLocalRoleDto
	var err error
	if models.CoreToken.BaseUrl == "" || len(roleList) == 0 {
		return
	}
	if roleMap, err = rpc.QueryAllRoles("N", userToken, language); err != nil {
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

func (s *RoleService) GetRoleList(ids []string, userToken, language string) (result []*models.RoleTable, err error) {
	var roleMap map[string]*models.SimpleLocalRoleDto
	result = []*models.RoleTable{}
	if roleMap, err = rpc.QueryAllRoles("N", userToken, language); err != nil {
		log.Logger.Error("QueryAllRoles err:%+v", log.Error(err))
		return
	}
	if len(ids) == 0 {
		for key, value := range roleMap {
			result = append(result, &models.RoleTable{
				Id:          key,
				DisplayName: value.DisplayName,
				CoreId:      value.ID,
				Email:       value.Email,
			})
		}
		return
	}
	for _, id := range ids {
		if v, ok := roleMap[id]; ok {
			result = append(result, &models.RoleTable{
				Id:          v.Name,
				DisplayName: v.DisplayName,
				CoreId:      v.ID,
				Email:       v.Email,
			})
		}
	}
	return
}

func (s *RoleService) QueryUserByRoles(roles []string, userToken, language string) (result []string, err error) {
	result = []string{}
	var roleMap map[string]*models.SimpleLocalRoleDto
	var list []*models.SimpleLocalRoleDto
	if roleMap, err = rpc.QueryAllRoles("N", userToken, language); err != nil {
		log.Logger.Error("QueryAllRoles err:%+v", log.Error(err))
		return
	}
	if len(roles) > 0 {
		for _, role := range roles {
			if v, ok := roleMap[role]; ok {
				list = append(list, v)
			}
		}
	}
	existMap := make(map[string]int)
	for _, v := range list {
		if v.ID == "" {
			continue
		}
		tmpResult, tmpErr := rpc.QueryRolesUsers(v.ID, userToken, language)
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

func (s *RoleService) GetRoleAdministrators(role string, userToken, language string) (result []string, err error) {
	result = []string{}
	roleMap, err := rpc.QueryAllRoles("Y", userToken, language)
	if err != nil {
		return
	}
	if v, ok := roleMap[role]; ok {
		result = append(result, v.Administrator)
	}
	return
}

func (s *RoleService) GetRoleDisplayName(userToken, language string) (displayNameMap map[string]string, err error) {
	displayNameMap = make(map[string]string)
	var roleMap map[string]*models.SimpleLocalRoleDto
	if roleMap, err = rpc.QueryAllRoles("Y", userToken, language); err != nil {
		log.Logger.Error("QueryAllRoles err:%+v", log.Error(err))
		return
	}
	for key, value := range roleMap {
		displayNameMap[key] = value.DisplayName
	}
	return
}

func (s *RoleService) GetUserInfo(userName, userToken, language string) (dto *models.UserDto, err error) {
	return rpc.GetUserInfo(userName, userToken, language)
}
