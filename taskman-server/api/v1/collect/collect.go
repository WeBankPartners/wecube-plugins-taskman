package collect

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/services/db"
	"github.com/gin-gonic/gin"
	"strings"
)

// AddTemplateCollect 添加模板收藏,需要区分新老数据.之前旧模板,需要先更新数据
func AddTemplateCollect(c *gin.Context) {
	var param models.AddCollectTemplateParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if strings.TrimSpace(param.TemplateId) == "" {
		err := fmt.Errorf("templateId is empty")
		middleware.ReturnParamValidateError(c, err)
		return
	}
	var parentId string
	requestTemplate, err := db.GetSimpleRequestTemplate(param.TemplateId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	parentId = strings.TrimSpace(requestTemplate.ParentId)
	if parentId == "" {
		// parentId为空说明 模板为老数据,需要更新该名称的模板
		parentId = db.UpdateRequestTemplateParentId(requestTemplate)
		// 可能由于到导入到模板,模板的 recordId值是错误的,导致parentId 还是空,此处直接更新当前模板parentId值
		if parentId == "" {
			parentId = param.TemplateId
			err = db.UpdateRequestTemplateParentIdById(param.TemplateId, parentId)
			if err != nil {
				middleware.ReturnServerHandleError(c, err)
				return
			}
		}
	}
	// 判断模板是否已经收藏
	if db.CheckUserCollectExist(parentId, middleware.GetRequestUser(c)) {
		middleware.ReturnTemplateAlreadyCollectError(c)
		return
	}
	collectTemplate := &models.CollectTemplateTable{
		RequestTemplate: parentId,
		User:            middleware.GetRequestUser(c),
		Role:            param.Role,
		Type:            requestTemplate.Type,
	}
	err = db.AddTemplateCollect(collectTemplate)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

// CancelTemplateCollect  取消收藏模板
func CancelTemplateCollect(c *gin.Context) {
	templateId := c.Param("templateId")
	parentId, err := getRequestTemplateParentId(templateId)
	if err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err = db.DeleteTemplateCollect(parentId, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

// QueryTemplateCollect  获取收藏模板列表
func QueryTemplateCollect(c *gin.Context) {
	var param models.QueryCollectTemplateParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	pageInfo, rowData, err := db.QueryTemplateCollect(&param, middleware.GetRequestUser(c), c.GetHeader("Authorization"), middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnPageData(c, pageInfo, rowData)
}

// FilterItem 过滤数据
func FilterItem(c *gin.Context) {
	var param models.FilterRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	data, err := db.GetCollectFilterItem(&param, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, data)
}

// getRequestTemplateParentId  根据模板id查找 最开始版本模板id
func getRequestTemplateParentId(templateId string) (string, error) {
	// 根据 templateId 查找parent_id,模板会变更产生多个版本,只需要关联最开始版本
	requestTemplate, err := db.GetSimpleRequestTemplate(templateId)
	if err != nil {
		return "", err
	}
	if requestTemplate.Id == "" {
		return "", fmt.Errorf("templateId:%s err", templateId)
	}
	if requestTemplate.ParentId != "" {
		return requestTemplate.ParentId, nil
	}
	return templateId, nil
}
