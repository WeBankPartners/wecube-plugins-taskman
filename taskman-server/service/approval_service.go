package service

import (
	"errors"
	"strings"

	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
)

type ApprovalService struct {
	approvalDao     *dao.ApprovalDao
	approvalRoleDao *dao.ApprovalRoleDao
}

func (s *ApprovalService) CreateApprovals(param []*models.ApprovalDto) ([]*dao.ExecAction, error) {
	if len(param) == 0 {
		return nil, nil
	}
	actions := []*dao.ExecAction{}
	// 校验参数
	approvalTemplateId := param[0].ApprovalTemplate
	requestId := param[0].Request
	approvalRoleHandlerTypes := make([][]string, len(param))
	for i, approval := range param {
		if approval.ApprovalTemplate != approvalTemplateId {
			return nil, errors.New("param approvalTemplate not the same")
		}
		if approval.Request != requestId {
			return nil, errors.New("param request not the same")
		}
		if len(approval.RoleObjs) == 0 {
			continue
		}
		approvalRoleHandlerTypes[i] = make([]string, len(approval.RoleObjs))
		for j, roleObj := range approval.RoleObjs {
			approvalTemplateRoleId := roleObj.ApprovalTemplateRole
			// 查询现有审批处理模板
			var approvalTemplateRoles []*models.ApprovalTemplateRoleTable
			err := GetApprovalTemplateService().approvalTemplateRoleDao.DB.SQL("SELECT * FROM approval_template_role WHERE id = ? AND approval_template = ?", approvalTemplateRoleId, approvalTemplateId).Find(&approvalTemplateRoles)
			if err != nil {
				return nil, err
			}
			if len(approvalTemplateRoles) == 0 {
				return nil, errors.New("no approval_template_role record found")
			}
			approvalRoleHandlerTypes[i][j] = approvalTemplateRoles[0].HandlerType
		}
	}
	// 查询现有审批模板
	var approvalTemplates []*models.ApprovalTemplateTable
	err := GetApprovalTemplateService().approvalTemplateDao.DB.SQL("SELECT * FROM approval_template WHERE id = ?", approvalTemplateId).Find(&approvalTemplates)
	if err != nil {
		return nil, err
	}
	if len(approvalTemplates) == 0 {
		return nil, errors.New("no approval_template record found")
	}
	approvalTemplate := approvalTemplates[0]
	// 查询现有审批列表
	var approvals []*models.ApprovalTable
	err = s.approvalDao.DB.SQL("SELECT * FROM approval WHERE request = ? ORDER BY sort", requestId).Find(&approvals)
	if err != nil {
		return nil, err
	}
	if len(approvals) != 0 {
		return nil, errors.New("approval record already exist")
	}
	// 插入数据
	for i, approval := range param {
		approval.Id = guid.CreateGuid()
		approval.Sort = i + 1
		// 插入审批
		action := &dao.ExecAction{Sql: "INSERT INTO approval (id,sort,approval_template,request,name,approve) VALUES (?,?,?,?,?,?)"}
		action.Param = []interface{}{approval.Id, approval.Sort, approval.ApprovalTemplate, approval.Request, approvalTemplate.Name, string(models.ApprovalApproveInit)}
		actions = append(actions, action)
		// 插入审批处理
		for j, roleObj := range approval.RoleObjs {
			action := &dao.ExecAction{Sql: "INSERT INTO approval_role (id,sort,approval_template_role,approval,role,handler,approve,handler_type) VALUES (?,?,?,?,?,?,?,?)"}
			action.Param = []interface{}{guid.CreateGuid(), j + 1, roleObj.ApprovalTemplateRole, approval.Id, roleObj.Role, roleObj.Handler, string(models.ApprovalRoleApproveInit), approvalRoleHandlerTypes[i][j]}
			actions = append(actions, action)
		}
	}
	return actions, nil
}

func (s *ApprovalService) UpdateApprovals(param []*models.ApprovalDto) ([]*dao.ExecAction, error) {
	if len(param) == 0 {
		return nil, nil
	}
	actions := []*dao.ExecAction{}
	// 校验参数
	approvalTemplateId := param[0].ApprovalTemplate
	requestId := param[0].Request
	approvalRoleHandlerTypes := make([][]string, len(param))
	approvalRoleHandlerTypeMap := make([]map[string]string, len(param))
	for i, approval := range param {
		if approval.ApprovalTemplate != approvalTemplateId {
			return nil, errors.New("param approvalTemplate not the same")
		}
		if approval.Request != requestId {
			return nil, errors.New("param request not the same")
		}
		if approval.Sort != i+1 {
			return nil, errors.New("param sort wrong")
		}
		if len(approval.RoleObjs) == 0 {
			continue
		}
		approvalRoleHandlerTypes[i] = make([]string, len(approval.RoleObjs))
		approvalRoleHandlerTypeMap[i] = make(map[string]string)
		for j, roleObj := range approval.RoleObjs {
			approvalTemplateRoleId := roleObj.ApprovalTemplateRole
			// 查询现有审批处理模板
			var approvalTemplateRoles []*models.ApprovalTemplateRoleTable
			err := GetApprovalTemplateService().approvalTemplateRoleDao.DB.SQL("SELECT * FROM approval_template_role WHERE id = ? AND approval_template = ?", approvalTemplateRoleId, approvalTemplateId).Find(&approvalTemplateRoles)
			if err != nil {
				return nil, err
			}
			if len(approvalTemplateRoles) == 0 {
				return nil, errors.New("no approval_template_role record found")
			}
			approvalRoleHandlerTypes[i][j] = approvalTemplateRoles[0].HandlerType
			approvalRoleHandlerTypeMap[i][approvalTemplateRoleId] = approvalTemplateRoles[0].HandlerType
		}
	}
	// 查询现有审批模板
	var approvalTemplates []*models.ApprovalTemplateTable
	err := GetApprovalTemplateService().approvalTemplateDao.DB.SQL("SELECT * FROM approval_template WHERE id = ?", approvalTemplateId).Find(&approvalTemplates)
	if err != nil {
		return nil, err
	}
	if len(approvalTemplates) == 0 {
		return nil, errors.New("no approval_template record found")
	}
	approvalTemplate := approvalTemplates[0]
	// 查询现有审批列表
	var approvals []*models.ApprovalTable
	err = s.approvalDao.DB.SQL("SELECT * FROM approval WHERE request = ? AND approval_template = ? ORDER BY sort", requestId, approvalTemplateId).Find(&approvals)
	if err != nil {
		return nil, err
	}
	// 汇总审批列表
	approvalIds := make([]string, len(approvals))
	for i, approval := range approvals {
		approvalIds[i] = approval.Id
	}
	// 查询现有审批处理列表
	var approvalRoles []*models.ApprovalRoleTable
	err = s.approvalRoleDao.DB.SQL("SELECT * FROM approval_role WHERE approval IN ('" + strings.Join(approvalIds, "','") + "') ORDER BY approval, sort").Find(&approvalRoles)
	if err != nil {
		return nil, err
	}
	// 汇总审批处理列表
	approvalRoleMap := make(map[string][]*models.ApprovalRoleTable)
	approvalId := ""
	for _, approvalRole := range approvalRoles {
		if approvalId != approvalRole.Approval {
			approvalId = approvalRole.Approval
			approvalRoleMap[approvalId] = make([]*models.ApprovalRoleTable, 0)
		}
		approvalRoleMap[approvalId] = append(approvalRoleMap[approvalId], approvalRole)
	}
	// 对比请求和现有数据
	for i, approval := range param {
		if approval.Id == "" {
			approval.Id = guid.CreateGuid()
		}
		approvalRoles, ok := approvalRoleMap[approval.Id]
		if !ok {
			// 插入审批
			action := &dao.ExecAction{Sql: "INSERT INTO approval (id,sort,approval_template,request,name,approve) VALUES (?,?,?,?,?,?)"}
			action.Param = []interface{}{approval.Id, approval.Sort, approval.ApprovalTemplate, approval.Request, approvalTemplate.Name, string(models.ApprovalApproveInit)}
			actions = append(actions, action)
			// 插入审批处理
			for j, roleObj := range approval.RoleObjs {
				action := &dao.ExecAction{Sql: "INSERT INTO approval_role (id,sort,approval_template_role,approval,role,handler,approve,handler_type) VALUES (?,?,?,?,?,?,?,?)"}
				action.Param = []interface{}{guid.CreateGuid(), j + 1, roleObj.ApprovalTemplateRole, approval.Id, roleObj.Role, roleObj.Handler, string(models.ApprovalRoleApproveInit), approvalRoleHandlerTypes[i][j]}
				actions = append(actions, action)
			}
		} else {
			// 更新审批
			action := &dao.ExecAction{Sql: "UPDATE approval SET sort = ?, name = ? WHERE id = ?"}
			action.Param = []interface{}{approval.Sort, approvalTemplate.Name, approval.Id}
			actions = append(actions, action)
			// 增删改审批处理
			for j, approvalRole := range approvalRoles {
				if j < len(approval.RoleObjs) {
					roleObj := approval.RoleObjs[j]
					action = &dao.ExecAction{Sql: "UPDATE approval_role SET role = ?, handler = ?, approve = ?, handler_type = ? WHERE id = ?"}
					action.Param = []interface{}{roleObj.Role, roleObj.Handler, roleObj.Approve, approvalRoleHandlerTypeMap[i][roleObj.ApprovalTemplateRole], approvalRole.Id}
					actions = append(actions, action)
				} else {
					action = &dao.ExecAction{Sql: "DELETE FROM approval_role WHERE id = ?"}
					action.Param = []interface{}{approvalRole.Id}
					actions = append(actions, action)
				}
			}
			for sort := len(approvalRoles) + 1; sort <= len(approval.RoleObjs); sort++ {
				roleObj := approval.RoleObjs[sort-1]
				action = &dao.ExecAction{Sql: "INSERT INTO approval_role (id,sort,approval_template_role,approval,role,handler,approve,handler_type) VALUES (?,?,?,?,?,?,?,?)"}
				action.Param = []interface{}{guid.CreateGuid(), sort, roleObj.ApprovalTemplateRole, approval.Id, roleObj.Role, roleObj.Handler, string(models.ApprovalRoleApproveInit), approvalRoleHandlerTypeMap[i][roleObj.ApprovalTemplateRole]}
				actions = append(actions, action)
			}
			// 处理完，从审批处理列表删除
			delete(approvalRoleMap, approval.Id)
		}
	}
	// 删除剩余的审批和审批处理
	for approvalId := range approvalRoleMap {
		action := &dao.ExecAction{Sql: "DELETE FROM approval_role WHERE approval = ?"}
		action.Param = []interface{}{approvalId}
		actions = append(actions, action)

		action = &dao.ExecAction{Sql: "DELETE FROM approval WHERE id = ?"}
		action.Param = []interface{}{approvalId}
		actions = append(actions, action)
	}
	return actions, nil
}

func (s *ApprovalService) DeleteApprovals(requestId string) ([]*dao.ExecAction, error) {
	actions := []*dao.ExecAction{}
	// 查询现有审批列表
	var approvals []*models.ApprovalTable
	err := s.approvalDao.DB.SQL("SELECT * FROM approval WHERE request = ? ORDER BY sort", requestId).Find(&approvals)
	if err != nil {
		return nil, err
	}
	if len(approvals) == 0 {
		return nil, nil
	}
	// 汇总审批列表
	approvalIds := make([]string, len(approvals))
	for i, approval := range approvals {
		approvalIds[i] = approval.Id
	}
	// 删除现有审批处理列表
	action := &dao.ExecAction{Sql: "DELETE FROM approval_role WHERE approval IN ('" + strings.Join(approvalIds, "','") + "')"}
	actions = append(actions, action)
	// 删除现有审批列表
	action = &dao.ExecAction{Sql: "DELETE FROM approval WHERE id IN ('" + strings.Join(approvalIds, "','") + "')"}
	actions = append(actions, action)
	return actions, nil
}

func (s *ApprovalService) ListApprovals(requestId string) ([]*models.ApprovalDto, error) {
	// 查询现有审批列表
	var approvals []*models.ApprovalTable
	err := s.approvalDao.DB.SQL("SELECT * FROM approval WHERE request = ? ORDER BY sort", requestId).Find(&approvals)
	if err != nil {
		return nil, err
	}
	if len(approvals) == 0 {
		return nil, nil
	}
	// 汇总审批列表
	approvalIds := make([]string, len(approvals))
	for i, approval := range approvals {
		approvalIds[i] = approval.Id
	}
	// 查询现有审批处理列表
	var approvalRoles []*models.ApprovalRoleTable
	err = s.approvalRoleDao.DB.SQL("SELECT * FROM approval_role WHERE approval IN ('" + strings.Join(approvalIds, "','") + "') ORDER BY approval, sort").Find(&approvalRoles)
	if err != nil {
		return nil, err
	}
	// 汇总审批处理列表
	approvalRoleMap := make(map[string][]*models.ApprovalRoleTable)
	approvalId := ""
	for _, approvalRole := range approvalRoles {
		if approvalId != approvalRole.Approval {
			approvalId = approvalRole.Approval
			approvalRoleMap[approvalId] = make([]*models.ApprovalRoleTable, 0)
		}
		approvalRoleMap[approvalId] = append(approvalRoleMap[approvalId], approvalRole)
	}
	// 构造返回结果
	result := make([]*models.ApprovalDto, len(approvals))
	for i, approval := range approvals {
		result[i] = &models.ApprovalDto{
			Id:               approval.Id,
			Sort:             approval.Sort,
			ApprovalTemplate: approval.ApprovalTemplate,
			Request:          approval.Request,
		}
		if roleObjs, ok := approvalRoleMap[approval.Id]; ok {
			result[i].RoleObjs = make([]*models.ApprovalRoleDto, len(roleObjs))
			for j, roleObj := range roleObjs {
				result[i].RoleObjs[j] = &models.ApprovalRoleDto{
					ApprovalTemplateRole: roleObj.ApprovalTemplateRole,
					Role:                 roleObj.Role,
					Handler:              roleObj.Handler,
					Approve:              roleObj.Approve,
				}
			}
		}
	}
	return result, nil
}
