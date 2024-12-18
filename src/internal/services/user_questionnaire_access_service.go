package services

import (
	"context"
	"errors"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"time"
)

type QuestionnaireAccessService interface {
	AssignRole(ctx context.Context, userID, questionnaireID, roleID uint, expiresAt *time.Time) error
	UpdateRole(ctx context.Context, accessID, newRoleID uint) error
	RevokeAccess(ctx context.Context, accessID uint) error
	GetUserPermissions(ctx context.Context, userID, questionnaireID uint) ([]models.QuestionnairePermission, error)
	ValidateAccess(ctx context.Context, userID, questionnaireID uint, action, resource string) error
	GetQuestionnaireUsers(ctx context.Context, questionnaireID uint) ([]models.UserQuestionnaireAccess, error)
	IsOwner(ctx context.Context, userID, questionnaireID uint) (bool, error)
}

type questionnaireAccessService struct {
	accessRepo   repositories.QuestionnaireAccessRepository
	questionRepo repositories.QuestionnaireRepository
}

func NewQuestionnaireAccessService(accessRepo repositories.QuestionnaireAccessRepository,
	questionRepo repositories.QuestionnaireRepository) QuestionnaireAccessService {
	return &questionnaireAccessService{
		accessRepo:   accessRepo,
		questionRepo: questionRepo,
	}
}

func (s *questionnaireAccessService) AssignRole(ctx context.Context, userID, questionnaireID, roleID uint, expiresAt *time.Time) error {

	existingAccess, err := s.accessRepo.GetUserAccess(ctx, userID, questionnaireID)
	if err == nil && existingAccess != nil {
		return errors.New("user already has access to this questionnaire")
	}

	access := &models.UserQuestionnaireAccess{
		UserID:          userID,
		QuestionnaireID: questionnaireID,
		RoleID:          roleID,
		ExpiresAt:       expiresAt,
		IsActive:        true,
	}

	return s.accessRepo.CreateAccess(ctx, access)
}

func (s *questionnaireAccessService) UpdateRole(ctx context.Context, accessID, newRoleID uint) error {
	access, err := s.accessRepo.GetAccessByID(ctx, accessID)
	if err != nil {
		return err
	}

	access.RoleID = newRoleID
	return s.accessRepo.UpdateAccess(ctx, access)
}

func (s *questionnaireAccessService) RevokeAccess(ctx context.Context, accessID uint) error {
	access, err := s.accessRepo.GetAccessByID(ctx, accessID)
	if err != nil {
		return err
	}

	access.IsActive = false
	return s.accessRepo.UpdateAccess(ctx, access)
}

func (s *questionnaireAccessService) GetUserPermissions(ctx context.Context, userID, questionnaireID uint) ([]models.QuestionnairePermission, error) {
	access, err := s.accessRepo.GetUserAccess(ctx, userID, questionnaireID)
	if err != nil {
		return nil, err
	}

	if access.Role == nil {
		return nil, errors.New("role not found")
	}

	return access.Role.Permissions, nil
}

func (s *questionnaireAccessService) ValidateAccess(ctx context.Context, userID, questionnaireID uint, action, resource string) error {

	isOwner, err := s.IsOwner(ctx, userID, questionnaireID)
	if err != nil {
		return err
	}
	if isOwner {
		return nil
	}

	hasPermission, err := s.accessRepo.HasPermission(ctx, userID, questionnaireID, action, resource)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("access denied")
	}

	access, err := s.accessRepo.GetUserAccess(ctx, userID, questionnaireID)
	if err != nil {
		return err
	}
	if access.ExpiresAt != nil && access.ExpiresAt.Before(time.Now()) {
		return errors.New("access expired")
	}

	return nil
}

func (s *questionnaireAccessService) IsOwner(ctx context.Context, userID, questionnaireID uint) (bool, error) {
	questionnaire, err := s.questionRepo.GetByID(questionnaireID)
	if err != nil {
		return false, err
	}
	return questionnaire.OwnerID == userID, nil
}

func (s *questionnaireAccessService) GetQuestionnaireUsers(ctx context.Context, questionnaireID uint) ([]models.UserQuestionnaireAccess, error) {
	return s.accessRepo.GetQuestionnaireUsers(ctx, questionnaireID)
}
