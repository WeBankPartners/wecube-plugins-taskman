package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
)

const (
	// pathQueryProcessDefinitions 查询编排
	pathQueryProcessDefinitions = "/platform/v1/process/definitions/list"
	// pathQueryModel 查询model
	pathQueryModel = "/platform/v1/models"
)

// QueryProcessDefinitionList 查询编排列表
func QueryProcessDefinitionList(userToken, language string, param models.QueryProcessDefinitionParam) (processList []*models.ProcDefDto, err error) {
	var response models.QueryProcessDefinitionResponse
	var procDefMap = make(map[string]*models.ProcDefDto)
	processList = make([]*models.ProcDefDto, 0)
	postBytes, _ := json.Marshal(param)
	byteArr, err := HttpPost(models.Config.Wecube.BaseUrl+pathQueryProcessDefinitions, userToken, language, postBytes)
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
		return
	}
	if len(response.Data) > 0 {
		for _, queryDto := range response.Data {
			if len(queryDto.ProcDefList) > 0 {
				for _, dto := range queryDto.ProcDefList {
					procDefMap[dto.Id] = dto
				}
			}
		}
	}
	if len(procDefMap) > 0 {
		if err != nil {
			return
		}
		for _, dto := range procDefMap {
			processList = append(processList, dto)
		}
	}
	return
}

// QueryAllModels 查询所有模型
func QueryAllModels(userToken, language string) (nodesList []*models.DataModel, err error) {
	var response models.QueryAllModelsResponse
	nodesList = make([]*models.DataModel, 0)
	byteArr, err := HttpGet(models.Config.Wecube.BaseUrl+pathQueryModel, userToken, language)
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
		nodesList = response.Data
	}
	return
}
