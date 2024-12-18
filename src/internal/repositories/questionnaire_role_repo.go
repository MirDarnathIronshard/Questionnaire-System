package repositories

import (
	"context"
	"github.com/QBG-P2/Voting-System/internal/models"
	"gorm.io/gorm"
)

type QuestionnaireRoleRepository interface {
	Create(ctx context.Context, role *models.QuestionnaireRole) error
	GetByQuestionnaireID(ctx context.Context, questionnaireID uint) ([]models.QuestionnaireRole, error)
	GetByUserID(ctx context.Context, userID uint) ([]models.QuestionnaireRole, error)
	GetByID(ctx context.Context, id uint) (*models.QuestionnaireRole, error)
	Update(ctx context.Context, role *models.QuestionnaireRole) error
	Delete(ctx context.Context, id uint) error
	GetUserQuestionnaireRoles(userID uint, questionnaireID uint) ([]models.QuestionnaireRole, error)
}

type questionnaireRoleRepository struct {
	db *gorm.DB
}

func NewQuestionnaireRoleRepository(db *gorm.DB) QuestionnaireRoleRepository {
	return &questionnaireRoleRepository{db}
}

func (r *questionnaireRoleRepository) Create(ctx context.Context, role *models.QuestionnaireRole) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *questionnaireRoleRepository) GetByQuestionnaireID(ctx context.Context, questionnaireID uint) ([]models.QuestionnaireRole, error) {
	var roles []models.QuestionnaireRole
	err := r.db.WithContext(ctx).
		Where("questionnaire_id = ?", questionnaireID).
		Find(&roles).Error
	return roles, err
}

func (r *questionnaireRoleRepository) GetByUserID(ctx context.Context, userID uint) ([]models.QuestionnaireRole, error) {
	var roles []models.QuestionnaireRole
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

func (r *questionnaireRoleRepository) GetByID(ctx context.Context, id uint) (*models.QuestionnaireRole, error) {
	var role models.QuestionnaireRole
	err := r.db.WithContext(ctx).First(&role, id).Error
	return &role, err
}

func (r *questionnaireRoleRepository) Update(ctx context.Context, role *models.QuestionnaireRole) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *questionnaireRoleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.QuestionnaireRole{}, id).Error
}

func (r *questionnaireRoleRepository) GetUserQuestionnaireRoles(userID uint, questionnaireID uint) ([]models.QuestionnaireRole, error) {
	var list []models.QuestionnaireRole
	r.db.Where("user_id ? and questionnaire_id", userID, questionnaireID).Find(&list)
	return list, nil
}
