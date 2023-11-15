package db

import (
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"time"
)

func AddTemplateCollect(param *models.CollectTemplateTable) error {
	param.Id = guid.CreateGuid()
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := x.Exec("insert into collect_template(id,request_template,account,created_time,updated_time) value (?,?,?,?,?)",
		param.Id, param.RequestTemplate, param.User, nowTime)
	if err != nil {
		err = fmt.Errorf("Insert database error:%s ", err.Error())
	}
	return err
}

func DeleteTemplateCollect(templateId, user string) error {
	_, err := x.Exec("delete from collect_template where request_template = ? and account = ?", templateId, user)
	if err != nil {
		err = fmt.Errorf("Delete database error:%s ", err.Error())
	}
	return err
}

func QueryTemplateCollect(param *models.QueryCollectTemplateObj) error {
	return nil
}
