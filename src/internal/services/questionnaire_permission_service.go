package services

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
)

type QuestionnairePermissionService interface {
	Create(permission *models.QuestionnairePermission) error
	GetByID(id uint) (*models.QuestionnairePermission, error)
	GetByQuestionnaireID(questionnaireID uint) ([]models.QuestionnairePermission, error)
	GetByUserID(userID uint) ([]models.QuestionnairePermission, error)
	Update(permission *models.QuestionnairePermission) error
	Delete(id uint) error
}

type questionnairePermissionService struct {
	repo repositories.QuestionnairePermissionRepository
}

func NewQuestionnairePermissionService(repo repositories.QuestionnairePermissionRepository) QuestionnairePermissionService {
	return &questionnairePermissionService{repo}
}

func (s *questionnairePermissionService) Create(permission *models.QuestionnairePermission) error {
	return s.repo.Create(permission)
}

func (s *questionnairePermissionService) GetByID(id uint) (*models.QuestionnairePermission, error) {
	return s.repo.GetByID(id)
}

func (s *questionnairePermissionService) GetByQuestionnaireID(questionnaireID uint) ([]models.QuestionnairePermission, error) {
	return s.repo.GetByQuestionnaireID(questionnaireID)
}

func (s *questionnairePermissionService) GetByUserID(userID uint) ([]models.QuestionnairePermission, error) {
	return s.repo.GetByUserID(userID)
}

func (s *questionnairePermissionService) Update(permission *models.QuestionnairePermission) error {
	return s.repo.Update(permission)
}

func (s *questionnairePermissionService) Delete(id uint) error {
	return s.repo.Delete(id)
}
