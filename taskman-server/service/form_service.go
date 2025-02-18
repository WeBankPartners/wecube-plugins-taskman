package service

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
)

type FormService struct {
	formDao     *dao.FormDao
	formItemDao *dao.FormItemDao
}

func (s *FormService) QueryFormItemListByRequest(requestId string) (result []*models.FormItemTable, err error) {
	if result, err = s.formItemDao.QueryByRequestId(requestId); err != nil {
		return
	}
	for _, formItem := range result {
		if form, err := s.formDao.Get(formItem.Form); err != nil {
			continue
		} else {
			formItem.RowDataId = form.DataId
		}
	}
	return
}
