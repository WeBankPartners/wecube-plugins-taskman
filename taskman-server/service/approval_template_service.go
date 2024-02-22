package service

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
)

type ApprovalTemplateService struct {
	approvalTemplateDao     dao.ApprovalTemplateDao
	approvalTemplateRoleDao dao.ApprovalTemplateRoleDao
}

func (s *ApprovalTemplateService) CreateApprovalTemplate(param *models.ApprovalTemplateCreateParam) (*models.ApprovalTemplateCreateResponse, error) {
	result := &models.ApprovalTemplateCreateResponse{}
	return result, nil
}

func (s *ApprovalTemplateService) UpdateApprovalTemplate(param *models.ApprovalTemplateDto) error {
	return nil
}

func (s *ApprovalTemplateService) DeleteApprovalTemplate(id string) (*models.ApprovalTemplateDeleteResponse, error) {
	result := &models.ApprovalTemplateDeleteResponse{}
	return result, nil
}

func (s *ApprovalTemplateService) GetApprovalTemplate(id string) (*models.ApprovalTemplateDto, error) {
	result := &models.ApprovalTemplateDto{}
	return result, nil
}

func (s *ApprovalTemplateService) ListApprovalTemplateIds(requestTemplateId string) (*models.ApprovalTemplateListIdsResponse, error) {
	result := &models.ApprovalTemplateListIdsResponse{}
	return result, nil
}

func (s *ApprovalTemplateService) ListApprovalTemplate(requestTemplateId string) ([]*models.ApprovalTemplateDto, error) {
	result := []*models.ApprovalTemplateDto{}
	return result, nil
}
