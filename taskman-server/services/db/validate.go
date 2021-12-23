package db

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"regexp"
	"strings"
)

func ValidateRequestForm(param []*models.RequestPreDataTableObj, userToken string) (err error) {
	for _, entityData := range param {
		if len(entityData.Value) == 0 {
			continue
		}
		err = validateFormDataRegular(entityData, userToken)
		if err != nil {
			err = fmt.Errorf("Entity:%s %s ", entityData.Entity, err.Error())
			break
		}
	}
	return err
}

func validateFormDataRegular(input *models.RequestPreDataTableObj, userToken string) (err error) {
	attrList, tmpErr := getCMDBAttributes(input.Entity, userToken)
	if tmpErr != nil {
		err = fmt.Errorf("Try to get CMDB attributes fail,%s ", tmpErr.Error())
		return
	}
	attrRegularMap := make(map[string]string)
	for _, v := range attrList {
		if v.RegularExpressionRule != "" {
			attrRegularMap[v.PropertyName] = v.RegularExpressionRule
		}
	}
	var titleIdList []string
	for _, v := range input.Title {
		titleIdList = append(titleIdList, v.Id)

	}
	formItemNameMap := getFormItemTemplateNameMap(titleIdList)
	for _, valueData := range input.Value {
		for k, v := range valueData.EntityData {
			if formItemObj, b := formItemNameMap[k]; b {
				vString := fmt.Sprintf("%s", v)
				if formItemObj.Regular != "" {
					if !validateRegular(vString, formItemObj.Regular) {
						err = fmt.Errorf("key:%s value:%s item form regular validate fail regular:%s", k, vString, formItemObj.Regular)
					}
				}
				if err != nil {
					break
				}
				if attrRegular, bb := attrRegularMap[k]; bb {
					if !validateRegular(vString, attrRegular) {
						err = fmt.Errorf("key:%s value:%s cmdb attribute regular validate fail regular:%s", k, vString, attrRegular)
					}
				}
				if err != nil {
					break
				}
			}
		}
		if err != nil {
			break
		}
	}
	return err
}

func getFormItemTemplateNameMap(idList []string) map[string]*models.FormItemTemplateTable {
	resultMap := make(map[string]*models.FormItemTemplateTable)
	var itemTemplateTable []*models.FormItemTemplateTable
	x.SQL("select * from form_item_template where id in ('" + strings.Join(idList, "','") + "')").Find(&itemTemplateTable)
	for _, v := range itemTemplateTable {
		resultMap[v.Name] = v
	}
	return resultMap
}

func validateRegular(input, regular string) bool {
	regObj, err := regexp.Compile(regular)
	if err != nil {
		log.Logger.Error("validate regular regexp compile fail", log.String("regular", regular), log.Error(err))
		return false
	}
	return regObj.MatchString(input)
}
