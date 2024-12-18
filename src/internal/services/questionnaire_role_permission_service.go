package services

import (
	"context"
	"errors"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
)

type QuestionnaireRolePermissionService interface {
	AssignPermissions(ctx context.Context, questionnaireID uint, roleID uint, permissionIDs []uint) error
	RemovePermissions(ctx context.Context, questionnaireID uint, roleID uint, permissionIDs []uint) error
	GetRolePermissions(ctx context.Context, questionnaireID uint, roleID uint) ([]models.QuestionnairePermission, error)
}

type questionnaireRolePermissionService struct {
	rolePermRepo repositories.QuestionnaireRolePermissionRepository
}

func NewQuestionnaireRolePermissionService(rolePermRepo repositories.QuestionnaireRolePermissionRepository) QuestionnaireRolePermissionService {
	return &questionnaireRolePermissionService{
		rolePermRepo: rolePermRepo,
	}
}

func (s *questionnaireRolePermissionService) AssignPermissions(ctx context.Context, questionnaireID uint, roleID uint, permissionIDs []uint) error {
	if len(permissionIDs) == 0 {
		return errors.New("no permissions provided")
	}

	for _, permID := range permissionIDs {
		rolePermission := &models.QuestionnaireRolePermission{
			QuestionnaireRoleID:       roleID,
			QuestionnairePermissionID: permID,
		}

		if err := rolePermission.Validate(); err != nil {
			return err
		}

		if err := s.rolePermRepo.Create(rolePermission); err != nil {
			return err
		}
	}

	return nil
}

func (s *questionnaireRolePermissionService) RemovePermissions(ctx context.Context, questionnaireID uint, roleID uint, permissionIDs []uint) error {
	if len(permissionIDs) == 0 {
		return errors.New("no permissions provided")
	}

	for _, permID := range permissionIDs {
		rolePermission := &models.QuestionnaireRolePermission{
			QuestionnaireRoleID:       roleID,
			QuestionnairePermissionID: permID,
		}

		if err := s.rolePermRepo.Delete(rolePermission); err != nil {
			return err
		}
	}

	return nil
}

func (s *questionnaireRolePermissionService) GetRolePermissions(ctx context.Context, questionnaireID uint, roleID uint) ([]models.QuestionnairePermission, error) {
	rolePerms, err := s.rolePermRepo.GetByRoleID(roleID)
	if err != nil {
		return nil, err
	}

	var permissions []models.QuestionnairePermission
	for _, rp := range rolePerms {
		permissions = append(permissions, rp.QuestionnairePermission)
	}

	return permissions, nil
}
