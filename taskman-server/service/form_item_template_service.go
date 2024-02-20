package service

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
)

type FormItemTemplateService struct {
	formItemTemplateDao dao.FormItemTemplateDao
	requestTemplateDao  dao.RequestTemplateDao
}

func (s FormItemTemplateService) UpdateFormTemplateItemGroup(param models.FormTemplateGroupConfigureDto) (err error) {
	return
}

func (s FormItemTemplateService) DeleteFormTemplateItemGroup(formTemplateId, itemGroupName string) (err error) {
	return
}

func (s FormItemTemplateService) UpdateFormTemplateItemGroupCustomData(param models.FormTemplateGroupCustomDataDto) (err error) {
	return
}

func (s FormItemTemplateService) SortFormTemplateItemGroup(param models.FormTemplateGroupSortDto) (err error) {

	return
}
