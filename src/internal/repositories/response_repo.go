package repositories

import (
	"context"
	"github.com/QBG-P2/Voting-System/internal/models"
	"gorm.io/gorm"
)

type ResponseRepository interface {
	Create(ctx context.Context, response *models.Response) error
	GetByQuestionnaireID(ctx context.Context, questionnaireID uint, page int, size int) ([]models.Response, error)
	GetByID(ctx context.Context, id uint) (*models.Response, error)
	Update(ctx context.Context, response *models.Response) error
	Delete(ctx context.Context, id uint) error
	GetTotalByQuestionnaireID(ctx context.Context, questionnaireID uint) (int, error)
	DeleteByUserAndQuestionnaire(ctx context.Context, userID, questionnaireID uint) error
	GetUserResponseCount(ctx context.Context, userID, questionnaireID uint) (int, error)
}

type responseRepository struct {
	db *gorm.DB
}

func NewResponseRepository(db *gorm.DB) ResponseRepository {
	return &responseRepository{db}
}

func (r *responseRepository) Create(ctx context.Context, response *models.Response) error {
	return r.db.WithContext(ctx).Create(response).Error
}

func (r *responseRepository) GetByQuestionnaireID(ctx context.Context, questionnaireID uint, page int, size int) ([]models.Response, error) {
	var responses []models.Response
	offset := (page - 1) * size
	err := r.db.WithContext(ctx).
		Where("questionnaire_id = ?", questionnaireID).
		Limit(size).
		Offset(offset).
		Find(&responses).Error
	return responses, err
}

func (r *responseRepository) GetByID(ctx context.Context, id uint) (*models.Response, error) {
	var response models.Response
	err := r.db.WithContext(ctx).First(&response, id).Error
	return &response, err
}

func (r *responseRepository) Update(ctx context.Context, response *models.Response) error {
	return r.db.WithContext(ctx).Save(response).Error
}

func (r *responseRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Response{}, id).Error
}

func (r *responseRepository) GetTotalByQuestionnaireID(ctx context.Context, questionnaireID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Response{}).
		Where("questionnaire_id = ?", questionnaireID).
		Count(&count).Error
	return int(count), err
}

func (r *responseRepository) DeleteByUserAndQuestionnaire(ctx context.Context, userID, questionnaireID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND questionnaire_id = ?", userID, questionnaireID).
		Delete(&models.Response{}).Error
}

func (r *responseRepository) GetUserResponseCount(ctx context.Context, userID, questionnaireID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Response{}).
		Where("user_id = ? AND questionnaire_id = ?", userID, questionnaireID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
