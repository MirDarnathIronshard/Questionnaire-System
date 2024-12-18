// File: ./src/internal/repositories/question_repository.go
package repositories

import (
	"context"

	"github.com/QBG-P2/Voting-System/internal/models"
	"gorm.io/gorm"
)

type QuestionRepository interface {
	Create(ctx context.Context, question *models.Question) error
	GetByQuestionnaireID(ctx context.Context, questionnaireID uint, page int, size int) ([]models.Question, error)
	GetByID(ctx context.Context, id uint) (*models.Question, error)
	Update(ctx context.Context, question *models.Question) error
	Delete(ctx context.Context, id uint) error
	GetTotalByQuestionnaireID(ctx context.Context, questionnaireID uint) (int, error)
}

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{db}
}

func (r *questionRepository) Create(ctx context.Context, question *models.Question) error {
	return r.db.WithContext(ctx).Create(question).Error
}

func (r *questionRepository) GetByQuestionnaireID(ctx context.Context, questionnaireID uint, page int, size int) ([]models.Question, error) {
	var questions []models.Question
	offset := (page - 1) * size
	err := r.db.WithContext(ctx).
		Where("questionnaire_id = ?", questionnaireID).
		Preload("Option").
		Preload("Response").
		Limit(size).
		Offset(offset).
		Find(&questions).Error
	return questions, err
}

func (r *questionRepository) GetByID(ctx context.Context, id uint) (*models.Question, error) {
	var question models.Question
	err := r.db.WithContext(ctx).Preload("Option").Preload("Response").First(&question, id).Error
	return &question, err
}

func (r *questionRepository) Update(ctx context.Context, question *models.Question) error {
	return r.db.WithContext(ctx).Save(question).Error
}

func (r *questionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Question{}, id).Error
}

func (r *questionRepository) GetTotalByQuestionnaireID(ctx context.Context, questionnaireID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Question{}).
		Where("questionnaire_id = ?", questionnaireID).
		Count(&count).Error
	return int(count), err
}

func (r *questionRepository) GetMonitoringData() {

}
