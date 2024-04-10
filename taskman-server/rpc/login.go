package rpc

import (
	"encoding/json"
	"fmt"

	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
)

func RemoteLogin(reqParam *models.LoginReq, userToken, language string) (result interface{}, err error) {
	var byteArr []byte
	var response models.RemoteLoginResp
	postBytes, _ := json.Marshal(reqParam)
	log.Logger.Info("Start request", log.String("param", string(postBytes)))
	byteArr, err = HttpPost(models.Config.Wecube.BaseUrl+"/auth/v1/api/taskLogin", userToken, language, postBytes)
	if err != nil {
		return
	}

	err = json.Unmarshal(byteArr, &response)
	if err != nil {
		err = fmt.Errorf("try to json unmarshal response body fail,%s ", err.Error())
		return
	}

	if response.Status != "OK" {
		err = fmt.Errorf(response.Message)
		return
	}

	result = response.Data
	return
}
