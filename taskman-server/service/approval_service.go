package service

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
)

type ApprovalService struct {
	approvalDao     dao.ApprovalDao
	approvalRoleDao dao.ApprovalRoleDao
}

func (s *ApprovalService) CreateApprovals(param []*models.ApprovalDto) ([]*dao.ExecAction, error) {
	actions := []*dao.ExecAction{}
	return actions, nil
}

func (s *ApprovalService) UpdateApprovals(param []*models.ApprovalDto) ([]*dao.ExecAction, error) {
	actions := []*dao.ExecAction{}
	return actions, nil
}

func (s *ApprovalService) DeleteApprovals(param []*models.ApprovalDto) ([]*dao.ExecAction, error) {
	actions := []*dao.ExecAction{}
	return actions, nil
}

func (s *ApprovalService) ListApproval(requestId string) ([]*models.ApprovalDto, error) {
	result := []*models.ApprovalDto{}
	return result, nil
}
