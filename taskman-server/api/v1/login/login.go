package login

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/WeBankPartners/go-common-lib/cipher"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/exterror"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/rpc"
	"github.com/gin-gonic/gin"
)

// GetSeed 获取 seed
func GetSeed(c *gin.Context) {
	md5sum := cipher.Md5Encode(models.Config.EncryptSeed)
	middleware.ReturnData(c, md5sum[0:16])
	return
}

func Login(c *gin.Context) {
	reqParam := models.LoginReq{}
	var err error
	if err = c.ShouldBindJSON(&reqParam); err != nil {
		middleware.ReturnError(c, exterror.Catch(exterror.New().RequestParamValidateError, err))
		return
	}

	if pwdBytes, pwdErr := base64.StdEncoding.DecodeString(reqParam.Password); pwdErr == nil {
		reqParam.Password = hex.EncodeToString(pwdBytes)
	} else {
		err = pwdErr
		middleware.ReturnError(c, exterror.Catch(exterror.New().RequestParamValidateError, err))
		return
	}

	if decodePwd, decodeErr := cipher.AesDePassword(models.Config.EncryptSeed, reqParam.Password); decodeErr == nil {
		reqParam.Password = decodePwd
	} else {
		err = decodeErr
		middleware.ReturnError(c, exterror.Catch(exterror.New().RequestParamValidateError, err))
		return
	}

	retData, err := rpc.RemoteLogin(&reqParam, "", c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnError(c, exterror.Catch(exterror.New().ServerHandleError, err))
		return
	} else {
		middleware.ReturnData(c, retData)
	}
	return
}
