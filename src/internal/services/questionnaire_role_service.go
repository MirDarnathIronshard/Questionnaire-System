package services

import (
	"context"
	"errors"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
)

type QuestionnaireRoleService struct {
	repo repositories.QuestionnaireRoleRepository
}

func NewQuestionnaireRoleService(repo repositories.QuestionnaireRoleRepository) *QuestionnaireRoleService {
	return &QuestionnaireRoleService{repo: repo}
}

func (s *QuestionnaireRoleService) CreateQuestionnaireRole(ctx context.Context, data *models.QuestionnaireRole) error {
	if err := data.Validate(); err != nil {
		return err
	}

	err := s.repo.Create(ctx, data)

	return err
}

func (s *QuestionnaireRoleService) UpdateQuestionnaireRole(ctx context.Context, data *models.QuestionnaireRole) error {
	if err := data.Validate(); err != nil {
		return err
	}

	err := s.repo.Update(ctx, data)

	return err
}

func (s *QuestionnaireRoleService) DeleteQuestionnaireRole(ctx context.Context, id uint) error {
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if data == nil {
		return errors.New("data not found")
	}

	err = s.repo.Delete(ctx, id)
	return err
}

func (s *QuestionnaireRoleService) GetQuestionnaireRoleByID(ctx context.Context, id uint) (*models.QuestionnaireRole, error) {
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, errors.New("data not found")
	}

	return data, err
}

func (s *QuestionnaireRoleService) GetUserQuestionnaireRoles(ctx context.Context, QuestionnaireID uint, userID uint) ([]models.QuestionnaireRole, error) {

	data, err := s.repo.GetUserQuestionnaireRoles(userID, QuestionnaireID)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (s *QuestionnaireRoleService) AssignRole(ctx context.Context, role *models.QuestionnaireRole) error {
	existingRoles, err := s.repo.GetByQuestionnaireID(ctx, role.QuestionnaireID)
	if err != nil {
		return err
	}
	for _, existingRole := range existingRoles {
		if existingRole.UserID == role.UserID {
			return errors.New("user already has a role in this questionnaire")
		}
	}
	return s.repo.Create(ctx, role)
}

func (s *QuestionnaireRoleService) GetRolesByQuestionnaireID(ctx context.Context, questionnaireID uint) ([]models.QuestionnaireRole, error) {
	return s.repo.GetByQuestionnaireID(ctx, questionnaireID)
}

func (s *QuestionnaireRoleService) GetRolesByUserID(ctx context.Context, userID uint) ([]models.QuestionnaireRole, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *QuestionnaireRoleService) GetRoleByID(ctx context.Context, id uint) (*models.QuestionnaireRole, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *QuestionnaireRoleService) UpdateRole(ctx context.Context, role *models.QuestionnaireRole) error {
	return s.repo.Update(ctx, role)
}

func (s *QuestionnaireRoleService) RemoveRole(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *QuestionnaireRoleService) IsOwnerOrAdmin(ctx context.Context, questionnaireID uint, userID uint) (bool, error) {
	roles, err := s.repo.GetByQuestionnaireID(ctx, questionnaireID)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		if role.UserID == userID && (role.Name == "admin") {
			return true, nil
		}
	}
	return false, nil
}
