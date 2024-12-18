package services

import (
	"context"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
)

type ResponseService interface {
	CreateResponse(ctx context.Context, response *models.Response) error
	GetResponsesByQuestionnaireID(ctx context.Context, questionnaireID uint, page int, size int) ([]models.Response, error)
	GetTotalResponsesByQuestionnaireID(ctx context.Context, questionnaireID uint) (int, error)
	GetResponseByID(ctx context.Context, id uint) (*models.Response, error)
	UpdateResponse(ctx context.Context, response *models.Response) error
	DeleteResponse(ctx context.Context, id uint) error
	IsOwnerResponse(ctx context.Context, responseID uint, userID uint) (bool, error)
}
type responseService struct {
	repo repositories.ResponseRepository
}

func NewResponseService(repo repositories.ResponseRepository) ResponseService {
	return &responseService{repo}
}

func (s *responseService) CreateResponse(ctx context.Context, response *models.Response) error {
	return s.repo.Create(ctx, response)
}

func (s *responseService) GetResponsesByQuestionnaireID(ctx context.Context, questionnaireID uint, page int, size int) ([]models.Response, error) {
	return s.repo.GetByQuestionnaireID(ctx, questionnaireID, page, size)
}

func (s *responseService) GetTotalResponsesByQuestionnaireID(ctx context.Context, questionnaireID uint) (int, error) {
	return s.repo.GetTotalByQuestionnaireID(ctx, questionnaireID)
}

func (s *responseService) GetResponseByID(ctx context.Context, id uint) (*models.Response, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *responseService) UpdateResponse(ctx context.Context, response *models.Response) error {
	return s.repo.Update(ctx, response)
}

func (s *responseService) DeleteResponse(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *responseService) IsOwnerResponse(ctx context.Context, responseID uint, userID uint) (bool, error) {
	response, err := s.repo.GetByID(ctx, responseID)
	if err != nil {
		return false, err
	}
	return response.UserID == userID, nil
}
