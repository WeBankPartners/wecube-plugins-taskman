package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/network"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
)

const (
	// pathRetrieveAllRoles  查询所有角色
	pathRetrieveAllRoles = "/platform/v1/roles/retrieve?all=%s"
	// pathRetrieveAllUser 查询所有用户
	pathRetrieveAllUser = "/platform/v1/users/retrieve"
)

// QueryAllRoles 查询所有角色
func QueryAllRoles(requiredAll, userToken, language string) (roleMap map[string]*models.SimpleLocalRoleDto, err error) {
	var response models.QueryRolesResponse
	var userMap map[string]*models.UserDto
	roleMap = make(map[string]*models.SimpleLocalRoleDto)
	byteArr, err := network.HttpGet(fmt.Sprintf(models.Config.Wecube.BaseUrl+pathRetrieveAllRoles, requiredAll), userToken, language)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteArr, &response)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if response.Status != "OK" {
		err = fmt.Errorf(response.Message)
	}
	if len(response.Data) > 0 {
		userMap, err = QueryAllUser(userToken, language)
		if err != nil {
			return
		}
		for _, data := range response.Data {
			if len(userMap) > 0 && userMap[data.Administrator] != nil {
				data.Administrator = userMap[data.Administrator].UserName
			}
			roleMap[data.Name] = data
		}
	}
	return
}

// QueryAllUser 查询所有用户
func QueryAllUser(userToken, language string) (userMap map[string]*models.UserDto, err error) {
	var response models.QueryUserResponse
	userMap = make(map[string]*models.UserDto)
	byteArr, err := network.HttpGet(models.Config.Wecube.BaseUrl+pathRetrieveAllUser, userToken, language)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteArr, &response)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if response.Status != "OK" {
		err = fmt.Errorf(response.Message)
	}
	if len(response.Data) > 0 {
		for _, data := range response.Data {
			userMap[data.ID] = data
		}
	}
	return
}
