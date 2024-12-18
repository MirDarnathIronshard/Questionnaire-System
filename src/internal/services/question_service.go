package services

import (
	"context"
	"github.com/QBG-P2/Voting-System/internal/repositories"

	"github.com/QBG-P2/Voting-System/internal/models"
)

type QuestionService interface {
	CreateQuestion(ctx context.Context, question *models.Question) error
	GetQuestionsByQuestionnaireID(ctx context.Context, questionnaireID uint, page int, size int) ([]models.Question, error)
	GetQuestionByID(ctx context.Context, id uint) (*models.Question, error)
	UpdateQuestion(ctx context.Context, question *models.Question) error
	DeleteQuestion(ctx context.Context, id uint) error
	GetTotalQuestionsByQuestionnaireID(ctx context.Context, questionnaireID uint) (int, error)
	IsOwnerQuestionnaire(ctx context.Context, questionnaireID uint, userID uint) (bool, error)
}

type questionService struct {
	repo repositories.QuestionRepository
}

func NewQuestionService(repo repositories.QuestionRepository) QuestionService {
	return &questionService{repo}
}

func (s *questionService) CreateQuestion(ctx context.Context, question *models.Question) error {
	return s.repo.Create(ctx, question)
}

func (s *questionService) GetQuestionsByQuestionnaireID(ctx context.Context, questionnaireID uint, page int, size int) ([]models.Question, error) {
	return s.repo.GetByQuestionnaireID(ctx, questionnaireID, page, size)
}

func (s *questionService) GetQuestionByID(ctx context.Context, id uint) (*models.Question, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *questionService) UpdateQuestion(ctx context.Context, question *models.Question) error {
	return s.repo.Update(ctx, question)
}

func (s *questionService) DeleteQuestion(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *questionService) GetTotalQuestionsByQuestionnaireID(ctx context.Context, questionnaireID uint) (int, error) {
	return s.repo.GetTotalByQuestionnaireID(ctx, questionnaireID)
}

func (s *questionService) IsOwnerQuestionnaire(ctx context.Context, questionnaireID uint, userID uint) (bool, error) {

	return true, nil
}
