package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
)

const (
	// pathPlatformQueryEntities  Query
	pathPlatformQueryEntities = "/platform/v1/packages/%s/entities/%s/query"
)

func EntitiesQuery(param models.EntityQueryParam, packageName, entity, userToken, language string) (response models.EntityResponse, err error) {
	postBytes, _ := json.Marshal(param)
	byteArr, err := HttpPost(fmt.Sprintf(models.Config.Wecube.BaseUrl+pathPlatformQueryEntities, packageName, entity), userToken, language, postBytes)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteArr, &response)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	return
}
