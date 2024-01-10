package exterror

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"

	"io/ioutil"
	"reflect"
	"strings"
)

type CustomError struct {
	Key           string        `json:"key"`           // 错误编码
	PassEnable    bool          `json:"passEnable"`    // 透传其它服务报错，不用映射
	Code          int           `json:"code"`          // 错误码
	Message       string        `json:"message"`       // 错误信息模版
	DetailErr     error         `json:"detail"`        // 错误信息
	MessageParams []interface{} `json:"messageParams"` // 消息参数列表
}

func (c CustomError) Error() string {
	return c.Message
}

type ErrorTemplate struct {
	CodeMessageMap map[int]string `json:"-"`
	CodeKeyMap     map[int]string `json:"-"`

	Language string `json:"language"`
	Success  string `json:"success"`
	// request param validate error
	RequestParamValidateError CustomError `json:"request_param_validate_error"`
	RequestReadBodyError      CustomError `json:"request_read_body_error"`
	RequestJsonUnmarshalError CustomError `json:"request_json_unmarshal_error"`
	RequestTokenValidateError CustomError `json:"request_token_validate_error"`
	RequestTokenExpireError   CustomError `json:"request_token_expire_error"`
	// database error
	DatabaseQueryError      CustomError `json:"database_query_error"`
	DatabaseQueryEmptyError CustomError `json:"database_query_empty_error"`
	DatabaseExecuteError    CustomError `json:"database_execute_error"`
	// sever handle error
	ServerHandleError CustomError `json:"server_handle_error"`
	// 业务报错
	ApiPermissionDeny  CustomError `json:"api_permission_deny"`
	DataPermissionDeny CustomError `json:"data_permission_deny"`
	// TemplateAlreadyCollect 模板已收藏
	TemplateAlreadyCollect CustomError `json:"template_already_collect"`
	// UploadFileTooLarge 上传文件太大
	UploadFileTooLarge CustomError `json:"upload_file_too_large"`
	// ChangeTaskStatusError 任务状态修改失败
	ChangeTaskStatusError CustomError `json:"change_task_status_error"`
	// SubmitRequestNotPermission 请求提交没有权限
	SubmitRequestNotPermission CustomError `json:"submit_request_not_permission"`
	// SaveRequestNotPermission 请求保存没有权限
	SaveRequestNotPermission CustomError `json:"save_request_not_permission"`
	// TaskApproveNotPermission 任务审批没有权限
	TaskApproveNotPermission CustomError `json:"task_approve_not_permission"`
	// TaskSaveNotPermission 任务保存没有权限
	TaskSaveNotPermission CustomError `json:"task_save_not_permission"`
	// UpdateRequestHandlerStatusError 更新请求处理人失败
	UpdateRequestHandlerStatusError CustomError `json:"update_request_handler_status_error"`
	// RevokeRequestError 撤回请求失败
	RevokeRequestError CustomError `json:"revoke_request_error"`
	// GetRequestPreviewDataError 获取请求模板配置表单项
	GetRequestPreviewDataError CustomError `json:"get_request_preview_data_error"`
	// ImportTemplateVersionConflictError 模板导入版本冲突
	ImportTemplateVersionConflictError CustomError `json:"import_template_version_conflict_error"`
}

var (
	TemplateList      []*ErrorTemplate
	ErrorDetailReturn bool
)

func InitErrorTemplateList(dirPath string, detailReturn bool) (err error) {
	ErrorDetailReturn = detailReturn
	if !strings.HasSuffix(dirPath, "/") {
		dirPath = dirPath + "/"
	}
	fs, readDirErr := ioutil.ReadDir(dirPath)
	if readDirErr != nil {
		return readDirErr
	}
	if len(fs) == 0 {
		return fmt.Errorf("dirPath:%s is empty dir", dirPath)
	}
	for _, v := range fs {
		if !strings.HasSuffix(v.Name(), ".json") {
			continue
		}
		tmpFileBytes, _ := ioutil.ReadFile(dirPath + v.Name())
		tmpErrorTemplate := ErrorTemplate{}
		tmpErr := json.Unmarshal(tmpFileBytes, &tmpErrorTemplate)
		if tmpErr != nil {
			err = fmt.Errorf("unmarshal json file :%s fail,%s ", v.Name(), tmpErr.Error())
			continue
		}
		tmpErrorTemplate.Language = strings.Replace(v.Name(), ".json", "", -1)
		tmpErrorTemplate.CodeMessageMap = make(map[int]string)
		tmpErrorTemplate.CodeKeyMap = make(map[int]string)
		tmpRt := reflect.TypeOf(tmpErrorTemplate)
		tmpVt := reflect.ValueOf(tmpErrorTemplate)
		for i := 0; i < tmpRt.NumField(); i++ {
			if tmpRt.Field(i).Type.Name() == "CustomError" {
				tmpC := tmpVt.Field(i).Interface().(CustomError)
				tmpErrorTemplate.CodeMessageMap[tmpC.Code] = tmpC.Message
				tmpErrorTemplate.CodeKeyMap[tmpC.Code] = tmpRt.Field(i).Tag.Get("json")
			}
		}
		TemplateList = append(TemplateList, &tmpErrorTemplate)
	}
	if err == nil && len(TemplateList) == 0 {
		err = fmt.Errorf("i18n error template list empty")
	}
	return err
}

func New() (et ErrorTemplate) {
	et = ErrorTemplate{}
	if len(TemplateList) > 0 {
		et = *TemplateList[0]
	}
	return
}

func Catch(customErr CustomError, err error) CustomError {
	customErr.DetailErr = err
	return customErr
}

func GetErrorResult(headerLanguage string, err error) (errorCode int, errorKey, errorMessage string) {
	customErr, b := err.(CustomError)
	if !b {
		return -1, models.DefaultHttpErrorCode, err.Error()
	} else {
		errorCode = customErr.Code
		if headerLanguage == "" || customErr.PassEnable {
			errorMessage = buildErrMessage(customErr.Message, customErr.MessageParams)
			if customErr.DetailErr != nil && ErrorDetailReturn {
				errorMessage = fmt.Sprintf("%s (%s)", errorMessage, customErr.DetailErr.Error())
			}
			return
		}
		headerLanguage = strings.Replace(headerLanguage, ";", ",", -1)
		for _, lang := range strings.Split(headerLanguage, ",") {
			if strings.HasPrefix(lang, "q=") {
				continue
			}
			lang = strings.ToLower(lang)
			for _, template := range TemplateList {
				if template.Language == lang {
					if message, exist := template.CodeMessageMap[errorCode]; exist {
						errorMessage = buildErrMessage(message, customErr.MessageParams)
						errorKey = template.CodeKeyMap[errorCode]
					}
					break
				}
			}
			if errorMessage != "" {
				break
			}
		}
		if errorMessage == "" {
			errorMessage = buildErrMessage(customErr.Message, customErr.MessageParams)
		}
	}
	if customErr.DetailErr != nil && ErrorDetailReturn {
		errorMessage = fmt.Sprintf("%s (%s)", errorMessage, customErr.DetailErr.Error())
	}
	return
}

func buildErrMessage(templateMessage string, params []interface{}) (message string) {
	message = templateMessage
	if strings.Count(templateMessage, "%") == 0 {
		return
	}
	message = fmt.Sprintf(message, params...)
	return
}

func IsBusinessErrorCode(errorCode int) bool {
	if strings.HasPrefix(fmt.Sprintf("%d", errorCode), "2") {
		return true
	}
	return false
}
