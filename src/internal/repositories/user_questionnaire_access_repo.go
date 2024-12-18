package repositories

import (
	"context"
	"errors"
	"github.com/QBG-P2/Voting-System/internal/models"
	"gorm.io/gorm"
)

type QuestionnaireAccessRepository interface {
	CreateAccess(ctx context.Context, access *models.UserQuestionnaireAccess) error
	GetUserAccess(ctx context.Context, userID, questionnaireID uint) (*models.UserQuestionnaireAccess, error)
	UpdateAccess(ctx context.Context, access *models.UserQuestionnaireAccess) error
	DeleteAccess(ctx context.Context, accessID uint) error
	GetQuestionnaireUsers(ctx context.Context, questionnaireID uint) ([]models.UserQuestionnaireAccess, error)
	HasPermission(ctx context.Context, userID, questionnaireID uint, action, resource string) (bool, error)
	GetAccessByID(ctx context.Context, accessID uint) (*models.UserQuestionnaireAccess, error)
	TransferAccess(ctx context.Context, sellerID, buyerID, questionnaireID uint) error
}

type questionnaireAccessRepository struct {
	db *gorm.DB
}

func NewQuestionnaireAccessRepository(db *gorm.DB) QuestionnaireAccessRepository {
	return &questionnaireAccessRepository{db: db}
}

func (r *questionnaireAccessRepository) CreateAccess(ctx context.Context, access *models.UserQuestionnaireAccess) error {
	return r.db.WithContext(ctx).Create(access).Error
}

func (r *questionnaireAccessRepository) GetUserAccess(ctx context.Context, userID, questionnaireID uint) (*models.UserQuestionnaireAccess, error) {
	var access models.UserQuestionnaireAccess
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND questionnaire_id = ? AND is_active = true", userID, questionnaireID).
		Preload("Role.Permissions").
		First(&access).Error
	return &access, err
}

func (r *questionnaireAccessRepository) UpdateAccess(ctx context.Context, access *models.UserQuestionnaireAccess) error {
	return r.db.WithContext(ctx).Save(access).Error
}

func (r *questionnaireAccessRepository) DeleteAccess(ctx context.Context, accessID uint) error {
	return r.db.WithContext(ctx).Delete(&models.UserQuestionnaireAccess{}, accessID).Error
}

func (r *questionnaireAccessRepository) GetQuestionnaireUsers(ctx context.Context, questionnaireID uint) ([]models.UserQuestionnaireAccess, error) {
	var accesses []models.UserQuestionnaireAccess
	err := r.db.WithContext(ctx).
		Where("questionnaire_id = ? AND is_active = true", questionnaireID).
		Preload("User").
		Preload("Role.Permissions").
		Find(&accesses).Error
	return accesses, err
}

func (r *questionnaireAccessRepository) HasPermission(ctx context.Context, userID, questionnaireID uint, action, resource string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("questionnaire_role_permissions").Table("user_questionnaire_accesses").
		Joins("JOIN questionnaire_roles ON user_questionnaire_accesses.role_id = questionnaire_roles.id").
		Joins("JOIN questionnaire_role_permissions ON questionnaire_roles.id = questionnaire_role_permissions.questionnaire_role_id").
		Joins("JOIN questionnaire_permissions ON questionnaire_role_permissions.questionnaire_permission_id = questionnaire_permissions.id").
		Where("user_questionnaire_accesses.user_id = ? AND user_questionnaire_accesses.questionnaire_id = ? AND "+
			"questionnaire_permissions.action = ? AND questionnaire_permissions.resource = ? AND "+
			"user_questionnaire_accesses.is_active = true", userID, questionnaireID, action, resource).
		Count(&count).Error

	return count > 0, err
}

func (r *questionnaireAccessRepository) GetAccessByID(ctx context.Context, accessID uint) (*models.UserQuestionnaireAccess, error) {
	var access models.UserQuestionnaireAccess
	err := r.db.WithContext(ctx).
		Where("id = ? AND is_active = true", accessID).
		Preload("Role.Permissions").
		First(&access).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("access not found")
		}
		return nil, err
	}
	return &access, nil
}

func (r *questionnaireAccessRepository) TransferAccess(ctx context.Context, sellerID, buyerID, questionnaireID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.UserQuestionnaireAccess{}).
		Where("user_id = ? AND questionnaire_id = ?", sellerID, questionnaireID).
		Updates(map[string]interface{}{
			"user_id":        buyerID,
			"transaction_id": gorm.Expr("id"),
		}).Error
}
