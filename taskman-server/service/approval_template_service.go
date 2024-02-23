package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
)

type ApprovalTemplateService struct {
	approvalTemplateDao     dao.ApprovalTemplateDao
	approvalTemplateRoleDao dao.ApprovalTemplateRoleDao
}

func (s *ApprovalTemplateService) CreateApprovalTemplate(param *models.ApprovalTemplateCreateParam) (*models.ApprovalTemplateCreateResponse, error) {
	actions := []*dao.ExecAction{}
	// 查询现有审批模板列表
	var approvalTemplates []*models.ApprovalTemplateTable
	err := s.approvalTemplateDao.DB.SQL("SELECT * FROM approval_template WHERE request_template = ? ORDER BY sort", param.RequestTemplate).Find(&approvalTemplates)
	if err != nil {
		return nil, err
	}
	// 校验参数
	if param.Sort <= 0 || param.Sort > len(approvalTemplates)+1 {
		return nil, errors.New("param sort out of range")
	}
	// 插入新审批模板
	nowTime := time.Now().Format(models.DateTimeFormat)
	newApprovalTemplate := &models.ApprovalTemplateTable{
		Id:              guid.CreateGuid(),
		Sort:            param.Sort,
		RequestTemplate: param.RequestTemplate,
		Name:            param.Name,
		RoleType:        string(models.ApprovalTemplateRoleTypeCustom),
		CreatedTime:     nowTime,
		UpdatedTime:     nowTime,
	}
	action := &dao.ExecAction{Sql: "INSERT INTO approval_template (id,sort,request_template,name,expire_day,role_type,created_time,updated_time) VALUES (?,?,?,?,?,?,?,?)"}
	action.Param = []interface{}{newApprovalTemplate.Id, newApprovalTemplate.Sort, newApprovalTemplate.RequestTemplate, newApprovalTemplate.Name, newApprovalTemplate.ExpireDay, newApprovalTemplate.RoleType, newApprovalTemplate.CreatedTime, newApprovalTemplate.UpdatedTime}
	actions = append(actions, action)
	// 插入新审批处理模板
	newApprovalTemplateRole := &models.ApprovalTemplateRoleTable{
		Id:               guid.CreateGuid(),
		Sort:             1,
		ApprovalTemplate: newApprovalTemplate.Id,
		RoleType:         string(models.ApprovalTemplateRoleRoleTypeTemplate),
		HandlerType:      string(models.ApprovalTemplateRoleHandlerTypeTemplate),
	}
	action = &dao.ExecAction{Sql: "INSERT INTO approval_template_role (id,sort,approval_template,role_type,handler_type) VALUES (?,?,?,?,?)"}
	action.Param = []interface{}{newApprovalTemplateRole.Id, newApprovalTemplateRole.Sort, newApprovalTemplateRole.ApprovalTemplate, newApprovalTemplateRole.RoleType, newApprovalTemplateRole.HandlerType}
	actions = append(actions, action)
	// 如果不是尾插，则需更新现有数据的序号
	if param.Sort != len(approvalTemplates)+1 {
		for i := param.Sort; i < len(approvalTemplates)+1; i++ {
			t := approvalTemplates[i-1]
			t.Sort += 1
			t.UpdatedTime = nowTime
			action = &dao.ExecAction{Sql: "UPDATE approval_template SET sort = ?, updated_time = ? WHERE id = ?"}
			action.Param = []interface{}{t.Sort, t.UpdatedTime, t.Id}
			actions = append(actions, action)
		}
		tmp := make([]*models.ApprovalTemplateTable, len(approvalTemplates)+1)
		copy(tmp, approvalTemplates[:param.Sort-1])
		tmp[param.Sort-1] = newApprovalTemplate
		copy(tmp[param.Sort:], approvalTemplates[param.Sort-1:])
		approvalTemplates = tmp
	} else {
		approvalTemplates = append(approvalTemplates, newApprovalTemplate)
	}
	// 执行事务
	err = dao.Transaction(actions)
	if err != nil {
		return nil, err
	}
	// 构造返回结果
	result := &models.ApprovalTemplateCreateResponse{
		Id:              newApprovalTemplate.Id,
		Sort:            newApprovalTemplate.Sort,
		RequestTemplate: newApprovalTemplate.RequestTemplate,
		Name:            newApprovalTemplate.Name,
		Ids:             make([]*models.ApprovalTemplateIdObj, len(approvalTemplates)),
	}
	for i, approvalTemplate := range approvalTemplates {
		result.Ids[i] = &models.ApprovalTemplateIdObj{
			Id:   approvalTemplate.Id,
			Sort: approvalTemplate.Sort,
			Name: approvalTemplate.Name,
		}
	}
	return result, nil
}

func (s *ApprovalTemplateService) UpdateApprovalTemplate(param *models.ApprovalTemplateDto) error {
	actions := []*dao.ExecAction{}
	// 查询现有审批模板
	var approvalTemplates []*models.ApprovalTemplateTable
	err := s.approvalTemplateDao.DB.SQL("SELECT * FROM approval_template WHERE id = ?", param.Id).Find(&approvalTemplates)
	if err != nil {
		return err
	}
	if len(approvalTemplates) == 0 {
		return errors.New("no approval_template record found")
	}
	approvalTemplate := approvalTemplates[0]
	// 校验参数
	if approvalTemplate.Sort != param.Sort || approvalTemplate.RequestTemplate != param.RequestTemplate {
		return errors.New("param sort or requestTemplate wrong")
	}
	for _, roleObj := range param.RoleObjs {
		if roleObj.RoleType != string(models.ApprovalTemplateRoleRoleTypeTemplate) &&
			(roleObj.HandlerType == string(models.ApprovalTemplateRoleHandlerTypeTemplate) || roleObj.HandlerType == string(models.ApprovalTemplateRoleHandlerTypeTemplateSuggest)) {
			return fmt.Errorf("roleType %s not match handlerType %s", roleObj.RoleType, roleObj.HandlerType)
		}
		if roleObj.RoleType != string(models.ApprovalTemplateRoleRoleTypeTemplate) && roleObj.Role != "" {
			return fmt.Errorf("roleType %s not match role %s", roleObj.RoleType, roleObj.Role)
		}
		if roleObj.HandlerType != string(models.ApprovalTemplateRoleHandlerTypeTemplate) && roleObj.HandlerType != string(models.ApprovalTemplateRoleHandlerTypeTemplateSuggest) && roleObj.Handler != "" {
			return fmt.Errorf("handlerType %s not match handler %s", roleObj.HandlerType, roleObj.Handler)
		}
	}
	// 更新现有审批模板
	nowTime := time.Now().Format(models.DateTimeFormat)
	action := &dao.ExecAction{Sql: "UPDATE approval_template SET name = ?, expire_day = ?, description = ?, role_type = ?, updated_time = ? WHERE id = ?"}
	action.Param = []interface{}{param.Name, param.ExpireDay, param.Description, param.RoleType, nowTime, param.Id}
	actions = append(actions, action)
	// 增删改现有审批处理模板
	if param.RoleType == string(models.ApprovalTemplateRoleTypeAdmin) || param.RoleType == string(models.ApprovalTemplateRoleTypeAuto) {
		action = &dao.ExecAction{Sql: "DELETE FROM approval_template_role WHERE approval_template = ?"}
		action.Param = []interface{}{param.Id}
		actions = append(actions, action)
	} else {
		// 查询现有审批处理模板
		var approvalTemplateRoles []*models.ApprovalTemplateRoleTable
		err = s.approvalTemplateRoleDao.DB.SQL("SELECT * FROM approval_template_role WHERE approval_template = ? ORDER BY sort", param.Id).Find(&approvalTemplateRoles)
		if err != nil {
			return err
		}
		// 对比增删改
		for i, approvalTemplateRole := range approvalTemplateRoles {
			if i < len(param.RoleObjs) {
				roleObj := param.RoleObjs[i]
				action = &dao.ExecAction{Sql: "UPDATE approval_template_role SET role_type = ?, handler_type = ?, role = ?, handler = ? WHERE id = ?"}
				action.Param = []interface{}{roleObj.RoleType, roleObj.HandlerType, roleObj.Role, roleObj.Handler, approvalTemplateRole.Id}
				actions = append(actions, action)
			} else {
				action = &dao.ExecAction{Sql: "DELETE FROM approval_template_role WHERE id = ?"}
				action.Param = []interface{}{approvalTemplateRole.Id}
				actions = append(actions, action)
			}
		}
		for sort := len(approvalTemplateRoles) + 1; sort <= len(param.RoleObjs); sort++ {
			roleObj := param.RoleObjs[sort-1]
			action = &dao.ExecAction{Sql: "INSERT INTO approval_template_role (id,sort,approval_template,role_type,handler_type,role,handler) VALUES (?,?,?,?,?,?,?)"}
			action.Param = []interface{}{guid.CreateGuid(), sort, param.Id, roleObj.RoleType, roleObj.HandlerType, roleObj.Role, roleObj.Handler}
			actions = append(actions, action)
		}
	}
	// 执行事务
	err = dao.Transaction(actions)
	if err != nil {
		return err
	}
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
	// 查询现有审批模板列表
	var approvalTemplates []*models.ApprovalTemplateTable
	err := s.approvalTemplateDao.DB.SQL("SELECT * FROM approval_template WHERE request_template = ? ORDER BY sort", requestTemplateId).Find(&approvalTemplates)
	if err != nil {
		return nil, err
	}
	if len(approvalTemplates) == 0 {
		return result, nil
	}
	// 汇总审批模板列表
	approvalTemplateIds := make([]string, len(approvalTemplates))
	for i, approvalTemplate := range approvalTemplates {
		approvalTemplateIds[i] = approvalTemplate.Id
	}
	// 查询现有审批处理模板列表
	var approvalTemplateRoles []*models.ApprovalTemplateRoleTable
	err = s.approvalTemplateRoleDao.DB.SQL("SELECT * FROM approval_template_role WHERE approval_template IN ('" + strings.Join(approvalTemplateIds, "','") + "') ORDER BY approval_template, sort").Find(&approvalTemplateRoles)
	if err != nil {
		return nil, err
	}
	// 汇总审批处理模板列表
	approvalTemplateRoleMap := make(map[string][]*models.ApprovalTemplateRoleTable)
	approvalTemplateId := ""
	for _, approvalTemplateRole := range approvalTemplateRoles {
		if approvalTemplateId != approvalTemplateRole.ApprovalTemplate {
			approvalTemplateId = approvalTemplateRole.ApprovalTemplate
			approvalTemplateRoleMap[approvalTemplateId] = make([]*models.ApprovalTemplateRoleTable, 0)
		}
		approvalTemplateRoleMap[approvalTemplateId] = append(approvalTemplateRoleMap[approvalTemplateId], approvalTemplateRole)
	}
	// 构造返回结果
	result = make([]*models.ApprovalTemplateDto, len(approvalTemplates))
	for i, approvalTemplate := range approvalTemplates {
		result[i] = &models.ApprovalTemplateDto{
			Id:              approvalTemplate.Id,
			Sort:            approvalTemplate.Sort,
			RequestTemplate: approvalTemplate.RequestTemplate,
			Name:            approvalTemplate.Name,
			ExpireDay:       approvalTemplate.ExpireDay,
			Description:     approvalTemplate.Description,
			RoleType:        approvalTemplate.RoleType,
		}
		if roleObjs, ok := approvalTemplateRoleMap[approvalTemplate.Id]; ok {
			result[i].RoleObjs = make([]*models.ApprovalTemplateRoleDto, len(roleObjs))
			for j, roleObj := range roleObjs {
				result[i].RoleObjs[j] = &models.ApprovalTemplateRoleDto{
					RoleType:    roleObj.RoleType,
					HandlerType: roleObj.HandlerType,
					Role:        roleObj.Role,
					Handler:     roleObj.Handler,
				}
			}
		}
	}
	return result, nil
}
