package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
)

const (
	// pathCmdbQueryEntities cmdb Query
	pathCmdbQueryEntities = "/wecmdb/entities/%s/query"
)

func CmdbEntitiesQuery(param models.EntityQueryParam, ciType, userToken, language string) (response models.EntityResponse, err error) {
	postBytes, _ := json.Marshal(param)
	byteArr, err := HttpPost(fmt.Sprintf(models.Config.Wecube.BaseUrl+pathCmdbQueryEntities, ciType), userToken, language, postBytes)
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
